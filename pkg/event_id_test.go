package pkg

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeLogID(t *testing.T) {
	type args struct {
		blockNo  uint64
		logIndex uint
	}
	tests := []struct {
		name      string
		args      args
		wantLogID int64
	}{
		{
			name: "encode success",
			args: args{
				blockNo:  1,
				logIndex: 1,
			},
			wantLogID: 2147483649,
		},
		{
			name: "encode max uint32",
			args: args{
				blockNo:  math.MaxUint32,
				logIndex: 1,
			},
			wantLogID: 9223372034707292161,
		},
		{
			name: "encode uint32 uint16 max",
			args: args{
				blockNo:  math.MaxUint32,
				logIndex: math.MaxUint16,
			},
			wantLogID: 9223372034707357695,
		},
		{
			name: "encode zero",
			args: args{
				blockNo:  0,
				logIndex: 0,
			},
			wantLogID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgID, err := EncodeLogID(tt.args.blockNo, tt.args.logIndex)
			require.NoError(t, err)
			require.Equal(t, tt.wantLogID, msgID)
		})
	}
}

func TestDecodeLogID(t *testing.T) {
	type args struct {
		logID int64
	}
	tests := []struct {
		name         string
		args         args
		wantBlockNo  uint64
		wantLogIndex uint16
	}{
		{
			name: "decode success",
			args: args{
				logID: 2147483649,
			},
			wantBlockNo:  1,
			wantLogIndex: 1,
		},
		{
			name: "decode max uint32",
			args: args{
				logID: 9223372034707292161,
			},
			wantBlockNo:  math.MaxUint32,
			wantLogIndex: 1,
		},
		{
			name: "decode uint32 int16 uint16 max",
			args: args{
				logID: 9223372034707357695,
			},
			wantBlockNo:  math.MaxUint32,
			wantLogIndex: math.MaxUint16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNo, logIndex := DecodeLogID(tt.args.logID)

			require.Equal(t, tt.wantBlockNo, blockNo)
			require.Equal(t, tt.wantLogIndex, logIndex)
		})
	}
}
