package squashfs

import (
	"bytes"
	"io"

	"github.com/rclone/squashfs/internal/toreader"
)

type fragEntry struct {
	Start uint64
	Size  uint32
	_     uint32
}

func (r Reader) fragReader(index uint32, fragSize uint32) (io.Reader, error) {
	realSize := r.fragEntries[index].Size &^ (1 << 24)
	if realSize == 0 {
		return bytes.NewReader(make([]byte, fragSize)), nil
	}
	rdr := io.LimitReader(toreader.NewReader(r.r, int64(r.fragEntries[index].Start)), int64(realSize))
	if realSize != r.fragEntries[index].Size {
		return rdr, nil
	}
	return r.d.Reader(rdr)
}
