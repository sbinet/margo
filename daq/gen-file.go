// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"hash/adler32"
	"io"
	"math"
	"os"

	"github.com/sbinet/margo/daq"
)

func main() {
	enc := newEncoder(os.Stdout)
	enc.Write(daq.MagicHdr[:])
	enc.WriteU32(daq.MagicBeg)

	det1 := daq.Calorimeter{
		CellID: 4242,
		Ene:    10,
	}
	det2 := daq.Env{
		Name: "raspi-01",
		H:    0.98,
		P:    1,
		T:    25.2,
	}
	buf := new(bytes.Buffer)
	{
		enc := newEncoder(buf)
		enc.WriteU8(uint8(daq.CaloDet))
		enc.WriteU64(det1.CellID)
		enc.WriteF64(det1.Ene)

		enc.WriteU8(uint8(daq.EnvSysDet))
		enc.WriteStr(det2.Name)
		enc.WriteF64(det2.H)
		enc.WriteF64(det2.P)
		enc.WriteF64(det2.T)
	}
	hdr := daq.Header{
		Len:    int64(buf.Len()),
		Sum:    adler32.Checksum(buf.Bytes()),
		RunNbr: 42,
		EvtNbr: 66,
	}

	enc.WriteU64(uint64(hdr.Len))
	enc.WriteU32(hdr.Sum)
	enc.WriteU64(uint64(hdr.RunNbr))
	enc.WriteU64(uint64(hdr.EvtNbr))
	enc.Write(buf.Bytes())
	enc.WriteU32(daq.MagicEnd)
}

type encoder struct {
	w   io.Writer
	err error
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
}

func (enc *encoder) Write(p []byte) (int, error) {
	if enc.err != nil {
		return 0, enc.err
	}
	n, err := enc.w.Write(p)
	enc.err = err
	return n, err
}

func (enc *encoder) WriteU8(v uint8) {
	if enc.err != nil {
		return
	}
	var buf = [1]byte{v}
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *encoder) WriteU32(v uint32) {
	if enc.err != nil {
		return
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *encoder) WriteU64(v uint64) {
	if enc.err != nil {
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *encoder) WriteF64(v float64) {
	if enc.err != nil {
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
	_, enc.err = enc.w.Write(buf[:])
}

func (enc *encoder) WriteStr(v string) {
	if enc.err != nil {
		return
	}
	enc.WriteU64(uint64(len(v)))
	enc.Write([]byte(v))
}
