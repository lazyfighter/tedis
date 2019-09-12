package structure

import (
	"tedis/proxy/log"
	"tedis/proxy/util"
)

type setMeta struct {
	ExpireAt   int64
	FieldCount int64
}

func (t *TxStructure) SADD(key []byte, members [][]byte) (int, error) {
	if t.readWriter == nil {
		return 0, errWriteOnSnapshot
	}

	metaKey := t.encodeSetMetaKey(key)
	meta := t.loadSetMeta(metaKey)

	if util.IsExpired(meta.ExpireAt) {
		meta = &setMeta{}
	}
	for _, member := range members {

		log.Debugf("setKey is %s , member is %s", metaKey, member)
	}

	return 0, nil
}

func (t *TxStructure) updateMember(key []byte, value []byte) (int, error) {
	if t.readWriter == nil {
		return 0, errWriteOnSnapshot
	}

	return 0, nil
}

func (t *TxStructure) loadSetMeta(metaKey []byte) *setMeta {
	return &setMeta{
		ExpireAt:   1232313,
		FieldCount: 23,
	}
}
