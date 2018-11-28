// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"hash/adler32"
	"os"

	"github.com/sbinet/margo/daq"
)

func main() {
	enc := daq.NewEncoder(os.Stdout)
	enc.Write(daq.MagicHdr[:])

	for _, evt := range []struct {
		run  int64
		evt  int64
		det1 daq.Calorimeter
		det2 daq.Env
	}{
		{
			run: 42,
			evt: 66,
			det1: daq.Calorimeter{
				CellID: 4242,
				Ene:    10,
			},
			det2: daq.Env{
				Name: "raspi-01",
				H:    0.98,
				P:    1,
				T:    25.2,
			},
		},
		{
			run: 42,
			evt: 1337,
			det1: daq.Calorimeter{
				CellID: 0x42,
				Ene:    11,
			},
			det2: daq.Env{
				Name: "raspi-02",
				H:    0.986,
				P:    1.02,
				T:    28.2,
			},
		},
	} {

		det1 := evt.det1
		det2 := evt.det2
		buf := new(bytes.Buffer)
		{
			enc := daq.NewEncoder(buf)
			enc.WriteU8(uint8(daq.CaloDet))
			enc.WriteU64(det1.CellID)
			enc.WriteF64(det1.Ene)
			enc.WriteU32(daq.DetEnd)

			enc.WriteU8(uint8(daq.EnvSysDet))
			enc.WriteStr(det2.Name)
			enc.WriteF64(det2.H)
			enc.WriteF64(det2.P)
			enc.WriteF64(det2.T)
			enc.WriteU32(daq.DetEnd)
		}
		hdr := daq.Header{
			Len:    int64(buf.Len()),
			Sum:    adler32.Checksum(buf.Bytes()),
			RunNbr: evt.run,
			EvtNbr: evt.evt,
		}

		enc.WriteU32(daq.EvtBeg)
		enc.WriteU64(uint64(hdr.Len))
		enc.WriteU32(hdr.Sum)
		enc.WriteU64(uint64(hdr.RunNbr))
		enc.WriteU64(uint64(hdr.EvtNbr))
		enc.Write(buf.Bytes())
		enc.WriteU32(daq.EvtEnd)
	}
}
