// Copyright 2014, Vimeo, LLC. All rights reserved.
// Copyright 2021 A Bunch Tell LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found here: https://github.com/vimeo/go-iccjpeg/blob/master/LICENSE

// Package iccjpeg implements ICC profile extraction from JPEG files.
package iccjpeg

import (
	"errors"
	"io"
)

const (
	iccHeaderLen = 14
)

// GetICCRaw reads a JPEG from input and returns a buffer containing the raw ICC profile data.
// If no ICC profile is present, then the buffer may be of length 0.
func GetICCRaw(input io.Reader) ([]byte, error) {
	var err error

	/*
		var iccData [][]byte
		iccLength := 0
	*/
	readProfs := 0
	numMarkers := -1
	p := NewParser(input)
	p.ReadSOI()
	seg, err := p.GetSegment(app2Marker)
	if err != nil {
		return nil, err
	}
	if seg.Size < iccHeaderLen {
		return nil, errors.New("ICC segment invalid")
	}

	i := 11
	if string(seg.Data[:i]) != "ICC_PROFILE" || seg.Data[i] != 0 {
		return nil, errors.New("ICC segment invalid")
	}
	i++

	seqN := seg.Data[i]
	i++
	if seqN == 0 {
		return nil, errors.New("invalid sequence number")
	}

	num := seg.Data[i]
	i++
	if numMarkers == -1 {
		numMarkers = int(num)
		//iccData = make([][]byte, numMarkers)
	} else if int(num) != numMarkers {
		return nil, errors.New("invalid ICC segment (numMarkers != cur_num_markers)")
	}

	if int(seqN) > numMarkers {
		return nil, errors.New("invalid ICC segment (seqN > numMarkers)")
	}

	/*
		// Non-raw data
		iccData[seqN-1] = make([]byte, seg.Size-iccHeaderLen)
		copy(iccData[seqN-1], seg.Data[i:])
		iccLength += seg.Size - iccHeaderLen
	*/

	readProfs++

	if readProfs == numMarkers {
		return seg.Data, nil
	}

	return nil, nil
}
