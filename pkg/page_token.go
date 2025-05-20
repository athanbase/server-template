package pkg

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/pkg/errors"
)

func DecodePageToken(pageToken string) (id int64, err error) {
	data, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		err = errors.Wrap(err, "decode page token failed")
		return
	}

	id = BytesToInt64(data)
	return
}

func EncodePageToken(id int64) (pageToken string) {
	data := Int64ToBytes(id)
	return base64.StdEncoding.EncodeToString(data)
}

func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
