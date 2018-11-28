// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package daq

import (
	"encoding/binary"
	"io"
	"math"
)

type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Err() error {
	return enc.err
}

func (enc *Encoder) Write(p []byte) (int, error) {
	if enc.err != nil {
		return 0, enc.err
	}
	n, err := enc.w.Write(p)
	enc.err = err
	return n, err
}

func (enc *Encoder) WriteU8(v uint8) {
	if enc.err != nil {
		return
	}
	var buf = [1]byte{v}
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *Encoder) WriteU32(v uint32) {
	if enc.err != nil {
		return
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *Encoder) WriteU64(v uint64) {
	if enc.err != nil {
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *Encoder) WriteF64(v float64) {
	if enc.err != nil {
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *Encoder) WriteStr(v string) {
	if enc.err != nil {
		return
	}
	enc.WriteU64(uint64(len(v)))
	enc.Write([]byte(v))
}
