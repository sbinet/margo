// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package daq

import (
	"encoding/binary"
	"io"
	"math"
)

type Decoder struct {
	r   io.Reader
	buf []byte
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:   r,
		buf: make([]byte, 8),
	}
}

func (dec *Decoder) Read(p []byte) (int, error) {
	if dec.err != nil {
		return 0, dec.err
	}
	var n int
	n, dec.err = io.ReadFull(dec.r, p)
	return n, dec.err
}

func (dec *Decoder) Err() error {
	return dec.err
}

func (dec *Decoder) ReadU8() uint8 {
	if dec.load(1) != nil {
		return 0
	}
	return uint8(dec.buf[0])
}

func (dec *Decoder) ReadU32() uint32 {
	if dec.load(4) != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(dec.buf[:4])
}

func (dec *Decoder) ReadU64() uint64 {
	if dec.load(8) != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(dec.buf)
}

func (dec *Decoder) ReadI64() int64 {
	return int64(dec.ReadU64())
}

func (dec *Decoder) ReadF64() float64 {
	if dec.load(8) != nil {
		return 0
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(dec.buf))
}

func (dec *Decoder) ReadStr() string {
	n := dec.ReadU64()
	str := make([]byte, n)
	_, dec.err = io.ReadFull(dec.r, str)
	if dec.err != nil {
		return ""
	}
	return string(str)
}

func (dec *Decoder) ReadStreamHeader() [4]byte {
	var hdr [4]byte
	_, dec.err = io.ReadFull(dec.r, hdr[:])
	return hdr
}

func (dec *Decoder) ReadHeader() Header {
	hdr := Header{
		Len:    dec.ReadI64(),
		Sum:    dec.ReadU32(),
		RunNbr: dec.ReadI64(),
		EvtNbr: dec.ReadI64(),
	}
	return hdr
}

func (dec *Decoder) ReadCalorimeter() Calorimeter {
	calo := Calorimeter{
		CellID: dec.ReadU64(),
		Ene:    dec.ReadF64(),
	}
	return calo
}

func (dec *Decoder) ReadEnvSys() Env {
	env := Env{
		Name: dec.ReadStr(),
		H:    dec.ReadF64(),
		P:    dec.ReadF64(),
		T:    dec.ReadF64(),
	}
	return env
}

func (dec *Decoder) load(n int) error {
	if dec.err != nil {
		return dec.err
	}
	_, dec.err = io.ReadFull(dec.r, dec.buf[:n])
	return dec.err
}
