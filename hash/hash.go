package hash

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

type BinaryMarshaller interface {
	MarshalBinary() []byte
}

func H(key interface{}) uint64 {
	if key == nil {
		return 0
	}
	switch h := key.(type) {
	case int:
		return uint64(h)
	case int32:
		return uint64(h)
	case int64:
		return uint64(h)
	case uint32:
		return uint64(h)
	case uint64:
		return h
	case byte:
		return uint64(h)
	case []byte:
		enc := fnv.New64a()
		_, err := enc.Write(h)
		if err != nil {
			panic(err)
		}
		return enc.Sum64()
	case string:
		enc := fnv.New64a()
		_, err := enc.Write([]byte(h))
		if err != nil {
			panic("failed to encode fnv: " + err.Error())
		}
		return enc.Sum64()
	case BinaryMarshaller:
		byt := h.MarshalBinary()
		fnvEnc := fnv.New64a()
		_, err := fnvEnc.Write(byt)
		if err != nil {
			panic("failed to encode fnv: " + err.Error())
		}
		return fnvEnc.Sum64()
	default:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(key); err != nil {
			panic("failed to encode gob" + err.Error())
		}
		fnvEnc := fnv.New64a()
		_, err := fnvEnc.Write(buf.Bytes())
		if err != nil {
			panic("failed to encode fnv: " + err.Error())
		}
		return fnvEnc.Sum64()
	}
}
