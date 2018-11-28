// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package daq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

var (
	evtBeg = make([]byte, 4)
	evtEnd = make([]byte, 4)
	detEnd = make([]byte, 4)
)

func init() {
	binary.LittleEndian.PutUint32(evtBeg, EvtBeg)
	binary.LittleEndian.PutUint32(evtEnd, EvtEnd)
	binary.LittleEndian.PutUint32(detEnd, DetEnd)
}

func NewEvtScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(splitEvts)
	return s
}

func NewDetScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(splitDets)
	return s
}

func splitEvts(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, evtEnd); i >= 0 {
		return i + 4, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func splitDets(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, detEnd); i >= 0 {
		return i + 4, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func dropMagic(data []byte) []byte {
	if len(data) > 4 {
		return data[0 : len(data)-4]
	}
	return data
}
