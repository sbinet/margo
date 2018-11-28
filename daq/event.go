// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package daq contains types and functions to build a simple DAQ system.
package daq

type Header struct {
	Len    int64
	Sum    uint32
	RunNbr int64
	EvtNbr int64
}

var (
	MagicHdr = [4]byte{0, 'd', 'a', 'q'}
)

const (
	MagicBeg = 0xdeadcafe
	MagicEnd = 0xdeadfeed
)

type DetKind byte

const (
	CaloDet   DetKind = 0
	EnvSysDet DetKind = 1
)

type Calorimeter struct {
	CellID uint64
	Ene    float64
}

type Env struct {
	Name string
	H    float64
	P    float64
	T    float64
}
