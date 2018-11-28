// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"hash/adler32"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/sbinet/margo/daq"
)

func main() {
	addr := flag.String("listen", ":8000", "address to listen on")

	flag.Parse()

	log.Printf("listening on %q...", *addr)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %v", err)
		}

		go serve(c)
	}
}

func serve(c net.Conn) {
	log.Printf("accepted connection from %q...", c.RemoteAddr())
	defer c.Close()

	_, err := c.Write(daq.MagicHdr[:])
	if err != nil {
		log.Printf("could not write magic header: %v", err)
		return
	}

	buf := new(bytes.Buffer)

	for {
		buf.Reset()
		generate(buf)
		_, err = buf.WriteTo(c)
		if err != nil {
			log.Printf("could not send event: %v", err)
			return
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(15)*100))
	}
}

var (
	genCellID = func() uint64 {
		return uint64(rand.Int63n(16 * 16))
	}
	genEne = func() float64 {
		return 20 + 5*rand.NormFloat64()
	}
	genName = func() string {
		return fmt.Sprintf("raspi-%02d", rand.Intn(3))
	}
	genH = func() float64 {
		return 0.98 + 5*rand.NormFloat64()
	}
	genP = func() float64 {
		return 1 + 5*rand.NormFloat64()
	}
	genT = func() float64 {
		return 37.2 + 2*rand.NormFloat64()
	}
)

func generate(w io.Writer) {

	var (
		run int64 = 42
		evt int64 = 66
	)

	buf := new(bytes.Buffer)
	enc := daq.NewEncoder(buf)

	for i := 0; i < rand.Intn(50)+1; i++ {
		det := daq.Calorimeter{
			CellID: genCellID(),
			Ene:    genEne(),
		}

		enc.WriteU8(uint8(daq.CaloDet))
		enc.WriteU64(det.CellID)
		enc.WriteF64(det.Ene)
		enc.WriteU32(daq.DetEnd)
	}

	for i := 0; i < 3; i++ {
		det := daq.Env{
			Name: fmt.Sprintf("raspi-%02d", i),
			H:    genH(),
			P:    genP(),
			T:    genT(),
		}

		enc.WriteU8(uint8(daq.EnvSysDet))
		enc.WriteStr(det.Name)
		enc.WriteF64(det.H)
		enc.WriteF64(det.P)
		enc.WriteF64(det.T)
		enc.WriteU32(daq.DetEnd)
	}

	hdr := daq.Header{
		Len:    int64(buf.Len()),
		Sum:    adler32.Checksum(buf.Bytes()),
		RunNbr: run,
		EvtNbr: evt,
	}

	enc = daq.NewEncoder(w)
	enc.WriteU32(daq.EvtBeg)
	enc.WriteU64(uint64(hdr.Len))
	enc.WriteU32(hdr.Sum)
	enc.WriteU64(uint64(hdr.RunNbr))
	enc.WriteU64(uint64(hdr.EvtNbr))
	enc.Write(buf.Bytes())
	enc.WriteU32(daq.EvtEnd)
}
