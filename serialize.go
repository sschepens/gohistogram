package gohistogram

import (
	"encoding/binary"
	"math"
)

func (h *NumericHistogram) Bytes() []byte {
	b := make([]byte, 2*8*len(h.bins)+8+8+8)
	binary.LittleEndian.PutUint64(b, h.total)
	binary.LittleEndian.PutUint64(b[8:], uint64(h.maxbins))
	binary.LittleEndian.PutUint64(b[16:], uint64(len(h.bins)))
	p := 24
	for _, v := range h.bins {
		binary.LittleEndian.PutUint64(b[p:], math.Float64bits(v.Value))
		binary.LittleEndian.PutUint64(b[p+8:], v.Count)
		p += 16
	}
	return b
}

func NewHistogramBytes(buff []byte) *NumericHistogram {
	var total uint64
	var maxbins int
	var lbins int
	total = binary.LittleEndian.Uint64(buff)
	maxbins = int(binary.LittleEndian.Uint64(buff[8:]))
	lbins = int(binary.LittleEndian.Uint64(buff[16:]))
	bins := make([]Bin, lbins, maxbins+1)
	for i, p := 0, 24; i < lbins; i, p = i+1, p+16 {
		bins[i] = Bin{
			Value: math.Float64frombits(binary.LittleEndian.Uint64(buff[p:])),
			Count: binary.LittleEndian.Uint64(buff[p+8:]),
		}
	}
	return &NumericHistogram{
		bins:    bins,
		maxbins: maxbins,
		total:   total,
	}
}
