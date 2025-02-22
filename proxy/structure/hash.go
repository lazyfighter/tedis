// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.
//
// The following only applies to changes made to this file as part of ELEME development.
//
// Portions Copyright (c) 2019 ELEME, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied.  See the License for the specific language governing permissions and limitations
// under the License.

package structure

import (
	"bytes"
	"github.com/pingcap/parser/terror"
	"strconv"

	"github.com/juju/errors"
	"github.com/pingcap/tidb/kv"
	"tedis/proxy/prometheus"
	"tedis/proxy/util"
	"time"
)

// HashPair is the pair for (field, value) in a hash.
type HashPair struct {
	Field []byte
	Value []byte
}

type hashMeta struct {
	ExpireAt   int64
	FieldCount int64
}

func (meta hashMeta) Value() []byte {
	//buf := make([]byte, 8)
	//binary.BigEndian.PutUint64(buf[0:8], uint64(meta.FieldCount))
	//return buf
	return EncodeHashMetaValue(meta.ExpireAt, meta.FieldCount)
}

func (meta hashMeta) IsEmpty() bool {
	return meta.FieldCount <= 0
}

// HSet sets the string value of a hash field.
func (t *TxStructure) HSet(key []byte, field []byte, value []byte) (int, error) {
	if t.readWriter == nil {
		return 0, errWriteOnSnapshot
	}
	return t.updateHash(key, field, func([]byte) ([]byte, error) {
		return value, nil
	})
}

func (t *TxStructure) HMSet(key []byte, elements []*HashPair) ([]byte, error) {

	ms := &util.MarkSet{}
	if len(elements) > SeekThreshold {

		omap := make(map[string][]byte)
		err := t.iterateHash(key, func(field []byte, value []byte) error {
			omap[string(append([]byte{}, field...))] = append([]byte{}, value...)
			return nil
		})
		if err != nil {
			return nil, errors.Trace(err)
		}
		for _, e := range elements {
			field := string(e.Field)
			if omap[field] == nil {
				ms.Set(e.Field)
			}

			if bytes.Equal(omap[field], e.Value) {
				continue
			}

			dataKey := t.encodeHashDataKey(key, e.Field)
			if err = t.readWriter.Set(dataKey, e.Value); err != nil {
				return nil, errors.Trace(err)
			}
		}
	} else {
		for _, e := range elements {
			dataKey := t.encodeHashDataKey(key, e.Field)
			oldValue, err := t.loadHashValue(dataKey)
			if err != nil {
				return nil, errors.Trace(err)
			}

			if oldValue == nil {
				ms.Set(e.Field)
			}
			if bytes.Equal(oldValue, e.Value) {
				continue
			}
			if err = t.readWriter.Set(dataKey, e.Value); err != nil {
				return nil, errors.Trace(err)
			}
		}
	}

	metaKey := t.EncodeMetaKey(key)
	meta, err := t.loadHashMeta(metaKey)
	if err != nil {
		return nil, errors.Trace(err)
	}

	meta.FieldCount += int64(ms.Len())
	if err = t.readWriter.Set(metaKey, EncodeHashMetaValue(meta.ExpireAt, meta.FieldCount)); err != nil {
		return nil, errors.Trace(err)
	}
	return []byte("OK"), nil
}

// HGet gets the value of a hash field.
func (t *TxStructure) HGet(key []byte, field []byte) ([]byte, error) {
	dataKey := t.encodeHashDataKey(key, field)
	value, err := t.reader.Get(dataKey)
	if terror.ErrorEqual(err, kv.ErrNotExist) {
		err = nil
	}
	return value, errors.Trace(err)
}

// HInc increments the integer value of a hash field, by step, returns
// the value after the increment.
func (t *TxStructure) HInc(key []byte, field []byte, step int64) (int64, error) {
	if t.readWriter == nil {
		return 0, errWriteOnSnapshot
	}
	base := int64(0)
	_, err := t.updateHash(key, field, func(oldValue []byte) ([]byte, error) {
		if oldValue != nil {
			var err error
			base, err = strconv.ParseInt(string(oldValue), 10, 64)
			if err != nil {
				return nil, errors.Trace(err)
			}
		}
		base += step
		return []byte(strconv.FormatInt(base, 10)), nil
	})

	return base, errors.Trace(err)
}

// HGetInt64 gets int64 value of a hash field.
func (t *TxStructure) HGetInt64(key []byte, field []byte) (int64, error) {
	value, err := t.HGet(key, field)
	if err != nil || value == nil {
		return 0, errors.Trace(err)
	}

	var n int64
	n, err = strconv.ParseInt(string(value), 10, 64)
	return n, errors.Trace(err)
}

func (t *TxStructure) updateHash(key []byte, field []byte, fn func(oldValue []byte) ([]byte, error)) (int, error) {
	dataKey := t.encodeHashDataKey(key, field)
	oldValue, err := t.loadHashValue(dataKey)
	res := 0

	if err != nil {
		return 0, errors.Trace(err)
	}

	newValue, err := fn(oldValue)
	if err != nil {
		return 0, errors.Trace(err)
	}

	// Check if new value is equal to old value.
	if bytes.Equal(oldValue, newValue) {
		return 0, nil
	}

	if err = t.readWriter.Set(dataKey, newValue); err != nil {
		return 0, errors.Trace(err)
	}

	metaKey := t.EncodeMetaKey(key)
	meta, err := t.loadHashMeta(metaKey)
	if err != nil {
		return 0, errors.Trace(err)
	}

	if oldValue == nil {
		meta.FieldCount++
		if err = t.readWriter.Set(metaKey, EncodeHashMetaValue(meta.ExpireAt, meta.FieldCount)); err != nil {
			return 0, errors.Trace(err)
		}
		res = 1
	}

	return res, nil
}

// HLen gets the number of fields in a hash.
func (t *TxStructure) HLen(key []byte) (int, error) {
	metaKey := t.EncodeMetaKey(key)
	meta, err := t.loadHashMeta(metaKey)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return int(meta.FieldCount), nil
}

// HDel deletes one or more hash fields.
func (t *TxStructure) HDel(key []byte, fields [][]byte) (int, error) {
	if t.readWriter == nil {
		return 0, errWriteOnSnapshot
	}
	metaKey := t.EncodeMetaKey(key)
	meta, err := t.loadHashMeta(metaKey)
	if err != nil || meta.IsEmpty() {
		return 0, errors.Trace(err)
	}

	res := 0

	var value []byte
	for _, field := range fields {
		dataKey := t.encodeHashDataKey(key, field)

		value, err = t.loadHashValue(dataKey)
		if err != nil {
			return 0, errors.Trace(err)
		}

		if value != nil {
			if err = t.readWriter.Delete(dataKey); err != nil {
				return 0, errors.Trace(err)
			}
			res++
			meta.FieldCount--
		}
	}

	if meta.IsEmpty() {
		err = t.readWriter.Delete(metaKey)
	} else {
		err = t.readWriter.Set(metaKey, meta.Value())
	}

	return res, errors.Trace(err)
}

// HKeys gets all the fields in a hash.
func (t *TxStructure) HKeys(key []byte) ([][]byte, error) {
	var keys [][]byte
	err := t.iterateHash(key, func(field []byte, value []byte) error {
		keys = append(keys, append([]byte{}, field...))
		return nil
	})

	return keys, errors.Trace(err)
}

// HGetAll gets all the fields and values in a hash.
func (t *TxStructure) HGetAll(key []byte) ([][]byte, error) {
	var res []HashPair
	err := t.iterateHash(key, func(field []byte, value []byte) error {
		pair := HashPair{
			Field: append([]byte{}, field...),
			Value: append([]byte{}, value...),
		}
		res = append(res, pair)
		return nil
	})

	rets := make([][]byte, len(res)*2)
	for i, e := range res {
		rets[i*2], rets[i*2+1] = e.Field, e.Value
	}

	return rets, errors.Trace(err)
}

// HClear removes the hash value of the key.
func (t *TxStructure) HClear(key []byte) error {
	metaKey := t.EncodeMetaKey(key)
	meta, err := t.loadHashMeta(metaKey)
	if err != nil || meta.IsEmpty() {
		return errors.Trace(err)
	}

	err = t.iterateHash(key, func(field []byte, value []byte) error {
		k := t.encodeHashDataKey(key, field)
		return errors.Trace(t.readWriter.Delete(k))
	})

	if err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(t.readWriter.Delete(metaKey))
}

func (t *TxStructure) iterateHash(key []byte, fn func(k []byte, v []byte) error) error {
	//dataPrefix := t.hashDataKeyPrefix(key)
	//it, err := t.reader.Seek(dataPrefix)
	//if err != nil {
	//	return errors.Trace(err)
	//}
	//
	//var field []byte
	//
	//for it.Valid() {
	//	if !it.Key().HasPrefix(dataPrefix) {
	//		break
	//	}
	//
	//	_, field, err = t.decodeHashDataKey(it.Key())
	//	if err != nil {
	//		return errors.Trace(err)
	//	}
	//
	//	if err = fn(field, it.Value()); err != nil {
	//		return errors.Trace(err)
	//	}
	//
	//	err = it.Next()
	//	if err != nil {
	//		return errors.Trace(err)
	//	}
	//}

	return nil
}

func (t *TxStructure) loadHashMeta(metaKey []byte) (hashMeta, error) {
	v, err := t.reader.Get(metaKey)
	if terror.ErrorEqual(err, kv.ErrNotExist) {
		err = nil
	} else if err != nil {
		return hashMeta{}, errors.Trace(err)
	}

	meta := hashMeta{FieldCount: 0, ExpireAt: 0}
	if v == nil {
		return meta, nil
	}

	if len(v) != 17 {
		return meta, errInvalidHashMeta
	}

	flag, expireAt, count := DecodeMetaValue(v)
	if flag != HashData {
		return meta, errInvalidHashKeyFlag
	}
	meta.ExpireAt = expireAt
	meta.FieldCount = count

	return meta, nil
}

func (t *TxStructure) loadHashValue(dataKey []byte) ([]byte, error) {
	start := time.Now().UnixNano() / 1e6
	defer func() {
		end := time.Now().UnixNano() / 1e6
		prometheus.PrintInfoCounter.WithLabelValues("LoadHash").Inc()
		prometheus.CmdHistogram.WithLabelValues("LocadHash").Observe(float64(end - start))
	}()
	v, err := t.reader.Get(dataKey)
	if terror.ErrorEqual(err, kv.ErrNotExist) {
		err = nil
		v = nil
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	return v, nil
}
