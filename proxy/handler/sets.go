package handler

import (
	"context"
	"github.com/juju/errors"
	"tedis/proxy/log"
	"tedis/proxy/structure"
)

func (h *TxTikvHandler) SADD(key []byte, members [][]byte) (int, error) {

	if len(key) == 0 {
		return 0, errArguments("len(key) = %d, expect != 0", len(key))
	}

	requestContext := newRequestContext("SADD")

	if kerr := CheckKeySize(key); kerr != nil {
		log.Errorf("%s  key:%s", requestContext.id, key, kerr)
		return 0, kerr
	}

	res, err := CallWithRetry(requestContext, func() (interface{}, error) {
		txn, err := h.Store.Begin()
		if err != nil {
			return 0, errors.Trace(ErrBegionTXN)
		}

		tx := structure.NewStructure(txn, txn, h.NameSpace, h.IgnoreTTL)
		res, err := tx.SADD(key, members)
		if err == nil {
			err = txn.Commit(context.Background())
		}

		if err != nil {
			txn.Rollback()
		}
		return res, err
	})



	log.Infof("%s SADD %s", requestContext.id, key)
	return res.(int), errors.Trace(err)
}

//func (h *TxTikvHandler) SCARD( key []byte) (int, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return 0, err
//	}
//
//	context := newRequestContext("SCARD")
//	log.Infof("%s SCARD %s", context.id, key)
//
//
//	return int(0), nil
//}
//
//func (h *TxTikvHandler) SISMEMBER( key []byte, member []byte) (int, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return 0, err
//	}
//
//	if h.IsCompress {
//		member = util.Compress(member)
//	}
//
//	context := newRequestContext("SISMEMBER")
//	log.Infof("%s SISMEMBER %s", context.id, key)
//	err = h.Client.SIsmember(key, member)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return 0, nil
//	} else if err != nil {
//		return 0, err
//	}
//
//	return 1, nil
//
//}
//
//func (h *TxTikvHandler) SMEMBERS( key []byte) ([][]byte, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return nil, err
//	}
//
//	context := newRequestContext("SMEMBERS")
//	log.Infof("%s SMEMBERS %s", context.id, key)
//	members, err := h.Client.SMembers(key)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	var com_values [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			value, err = util.Uncompress(value)
//			if err != nil {
//				return nil, err
//			}
//			com_values = append(com_values, value)
//		}
//		return com_values, nil
//	}
//
//	return members, nil
//}
//
//func (h *TxTikvHandler) SPOP( key []byte) ([]byte, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return nil, err
//	}
//
//	context := newRequestContext("SPOP")
//	log.Infof("%s SPOP %s", context.id, key)
//	member, err := h.Client.SPop(key)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	if h.IsCompress {
//		member, err = util.Uncompress(member)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return member, nil
//
//}
//func (h *TxTikvHandler) SRANDMEMBER( key []byte, count []byte) ([][]byte, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return nil, err
//	}
//
//	cnt, err := strconv.Atoi(string(count))
//	if err != nil {
//		return nil, err
//	}
//	context := newRequestContext("SRANDMEMBER")
//	log.Infof("%s SRANDMEMBER %s", context.id, key)
//	members, err := h.Client.SRandmember(key, int32(cnt))
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	var com_values [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			value, err = util.Uncompress(value)
//			if err != nil {
//				return nil, err
//			}
//			com_values = append(com_values, value)
//		}
//		return com_values, nil
//	}
//
//	return members, nil
//
//}
//
//func (h *TxTikvHandler) SREM( key []byte, members [][]byte) (int, error) {
//	err := CheckKeySize(key)
//	if err != nil {
//		return 0, err
//	}
//
//	var com_members [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			com_members = append(com_members, util.Compress(value))
//		}
//	}
//
//	context := newRequestContext("SREM")
//	log.Infof("%s SREM %s", context.id, key)
//	var count int32
//	if h.IsCompress {
//		count, err = h.Client.SRem(key, com_members)
//	} else {
//		count, err = h.Client.SRem(key, members)
//	}
//	if err != nil && err == phxkv.ErrKeyExist {
//		return 0, nil
//	} else if err != nil {
//		return 0, err
//	}
//
//	return int(count), nil
//}
//
//func (h *TxTikvHandler) SDIFF( keys [][]byte) ([][]byte, error) {
//	for _, key := range keys {
//		err := CheckKeySize(key)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	context := newRequestContext("SDiff")
//	log.Infof("%s SDIFF %s", context.id)
//	members, err := h.Client.SDiff(keys)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	var com_values [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			value, err = util.Uncompress(value)
//			if err != nil {
//				return nil, err
//			}
//			com_values = append(com_values, value)
//		}
//		return com_values, nil
//	}
//
//	return members, nil
//}
//
//func (h *TxTikvHandler) SINTER( keys [][]byte) ([][]byte, error) {
//	for _, key := range keys {
//		err := CheckKeySize(key)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	context := newRequestContext("SINTER")
//	log.Infof("%s SINTER %s", context.id)
//	members, err := h.Client.SInter(keys)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	var com_values [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			value, err = util.Uncompress(value)
//			if err != nil {
//				return nil, err
//			}
//			com_values = append(com_values, value)
//		}
//		return com_values, nil
//	}
//
//	return members, nil
//}
//
//func (h *TxTikvHandler) SUNION(keys [][]byte) ([][]byte, error) {
//	for _, key := range keys {
//		err := CheckKeySize(key)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	context := newRequestContext("SUNION")
//	log.Infof("%s SUNION %s", context.id)
//	members, err := h.Client.SUnion(keys)
//	if err != nil && err == phxkv.ErrKeyExist {
//		return nil, nil
//	} else if err != nil {
//		return nil, err
//	}
//
//	var com_values [][]byte
//	if h.IsCompress {
//		for _, value := range members {
//			value, err = util.Uncompress(value)
//			if err != nil {
//				return nil, err
//			}
//			com_values = append(com_values, value)
//		}
//		return com_values, nil
//	}
//
//	return members, nil
//
//}
