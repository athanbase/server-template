package pkg

import (
	"math"

	"github.com/pkg/errors"
)

// pre 32 bit		 suf 16 bit
// blockNo		event index
// 4294967295	65535
const (
	logIndexMask = ^(-1 << 16) // ( 1 << 16)-1
	blockNoMask  = ^(-1 << 32) // (1 << 32)-1
)

func EncodeLogID(blockNo uint64, logIndex uint) (int64, error) {
	if blockNo > math.MaxUint32 ||
		logIndex > math.MaxUint16 {
		err := errors.Errorf(
			"number overflow, blockNo: %d, event index: %d",
			blockNo,
			logIndex,
		)
		return 0, err
	}

	var logID uint64
	logID |= blockNo << 31
	logID |= uint64(logIndex)
	return int64(logID), nil
}

func DecodeLogID(logID int64) (blockNo uint64, logIndex uint16) {
	blockNo = uint64(logID >> 31 & blockNoMask)
	logIndex = uint16(logID & logIndexMask)
	return
}
