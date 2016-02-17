package gohistogram

// Copyright (c) 2013 VividCortex, Inc. All rights reserved.
// Please see the LICENSE file for applicable license terms.

import (
	"fmt"
)

type NumericHistogram struct {
	bins    []Bin
	maxbins int
	total   uint64
}

// NewHistogram returns a new NumericHistogram with a maximum of n bins.
//
// There is no "optimal" bin count, but somewhere between 20 and 80 bins
// should be sufficient.
func NewHistogram(n int) *NumericHistogram {
	return &NumericHistogram{
		bins:    make([]Bin, 0, n+1),
		maxbins: n,
		total:   0,
	}
}

func (h *NumericHistogram) Reset() {
	h.total = 0
	h.bins = h.bins[:0]
}

func (h *NumericHistogram) Add(n float64) {
	h.total++
	for i, bini := range h.bins {
		if bini.Value == n {
			h.bins[i].Count++
			return
		}

		if bini.Value > n {
			// Insert a bin at i
			tail := h.bins[i:]
			h.bins = h.bins[0 : len(h.bins)+1]
			copy(h.bins[i+1:], tail)
			h.bins[i] = Bin{Value: n, Count: 1}

			h.trim()
			return
		}
	}

	h.bins = append(h.bins, Bin{Count: 1, Value: n})
	h.trim()
}

func (h *NumericHistogram) Quantile(q float64) float64 {
	count := q * float64(h.total)
	for i := range h.bins {
		count -= float64(h.bins[i].Count)

		if count <= 0 {
			return h.bins[i].Value
		}
	}

	return -1
}

// CDF returns the value of the cumulative distribution function
// at x
func (h *NumericHistogram) CDF(x float64) float64 {
	count := 0.0
	for i := range h.bins {
		if h.bins[i].Value <= x {
			count += float64(h.bins[i].Count)
		}
	}

	return count / float64(h.total)
}

// Mean returns the sample mean of the distribution
func (h *NumericHistogram) Mean() float64 {
	if h.total == 0 {
		return 0
	}

	sum := 0.0

	for i := range h.bins {
		sum += h.bins[i].Value * float64(h.bins[i].Count)
	}

	return sum / float64(h.total)
}

// Variance returns the variance of the distribution
func (h *NumericHistogram) Variance() float64 {
	if h.total == 0 {
		return 0
	}

	sum := 0.0
	mean := h.Mean()

	for i := range h.bins {
		sum += (float64(h.bins[i].Count) * (h.bins[i].Value - mean) * (h.bins[i].Value - mean))
	}

	return sum / float64(h.total)
}

func (h *NumericHistogram) Count() uint64 {
	return h.total
}

// trim merges adjacent bins to decrease the bin count to the maximum value
func (h *NumericHistogram) trim() {
	if len(h.bins) <= 1 {
		return
	}
	for len(h.bins) > h.maxbins {
		// Find closest bins in terms of value
		minDelta := h.bins[1].Value - h.bins[0].Value
		minDeltaIndex := 1
		for i := 2; i < len(h.bins); i++ {
			if i <= 1 {
				continue
			}

			if delta := h.bins[i].Value - h.bins[i-1].Value; delta < minDelta {
				minDelta = delta
				minDeltaIndex = i
			}
		}

		binMinDeltaIdxSub1 := h.bins[minDeltaIndex-1]
		binMinDeltaIdx := h.bins[minDeltaIndex]

		// We need to merge bins minDeltaIndex-1 and minDeltaIndex
		totalCount := binMinDeltaIdxSub1.Count + binMinDeltaIdx.Count
		mergedbin := Bin{
			Value: (binMinDeltaIdxSub1.Value*
				float64(binMinDeltaIdxSub1.Count) +
				binMinDeltaIdx.Value*
					float64(binMinDeltaIdx.Count)) /
				float64(totalCount), // weighted average
			Count: totalCount, // summed heights
		}

		// Now remove bin minDeltaIndex by replacing minDeltaIndex-1 with mergedbin
		head := h.bins[0:minDeltaIndex]
		tail := h.bins[minDeltaIndex+1:]
		head[minDeltaIndex-1] = mergedbin
		copy(h.bins[minDeltaIndex:], tail)
		h.bins = h.bins[:len(h.bins)-1]
	}
}

// String returns a string reprentation of the histogram,
// which is useful for printing to a terminal.
func (h *NumericHistogram) String() (str string) {
	str += fmt.Sprintln("Total:", h.total)

	for i := range h.bins {
		var bar string
		for j := 0; j < int(float64(h.bins[i].Count)/float64(h.total)*200); j++ {
			bar += "."
		}
		str += fmt.Sprintln(h.bins[i].Value, "\t", bar)
	}

	return
}

func (h *NumericHistogram) Bins() []Bin {
	return h.bins
}
