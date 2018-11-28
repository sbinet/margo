// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/sbinet/margo/daq"
)

func main() {
	addr := flag.String("dial", "localhost:8000", "address to dial")

	flag.Parse()

	var (
		r   io.Reader
		err error
	)

	switch flag.NArg() {
	case 1:
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	default:
		log.Printf("dialing %q...", *addr)

		c, err := net.Dial("tcp", *addr)
		if err != nil {
			log.Fatalf("could not dial %q: %v", *addr, err)
		}
		defer c.Close()
		r = c
	}

	dec := daq.NewDecoder(r)
	hdr := dec.ReadStreamHeader()
	if hdr != daq.MagicHdr {
		log.Fatalf("not a DAQ file: %q", hdr)
	}

	sc := daq.NewEvtScanner(r)
	for sc.Scan() {
		log.Printf(">>>...")
		dec := daq.NewDecoder(bytes.NewReader(sc.Bytes()))
		blk := dec.ReadU32()
		if blk != daq.EvtBeg {
			log.Fatalf("expected a magic-block beg, got: %#x", blk)
		}

		hdr := dec.ReadHeader()
		fmt.Printf("hdr: %#v\n", hdr)

		buf := make([]byte, hdr.Len)
		_, err = dec.Read(buf)
		if err != nil {
			log.Fatalf("could not read payload: %v", err)
		}

		det := daq.NewDetScanner(bytes.NewReader(buf))
		for det.Scan() {
			dec := daq.NewDecoder(bytes.NewReader(det.Bytes()))
			switch typ := daq.DetKind(dec.ReadU8()); typ {
			case daq.CaloDet:
				calo := dec.ReadCalorimeter()
				fmt.Printf("cal: %#v\n", calo)
			case daq.EnvSysDet:
				env := dec.ReadEnvSys()
				fmt.Printf("env: %#v\n", env)
			default:
				log.Fatalf("unknown detector type: %x", typ)
			}
		}
		if det.Err() != nil {
			log.Fatalf("sub error: %v", det.Err())
		}
	}
	if sc.Err() != nil {
		log.Fatalf("scanner error: %v", sc.Err())
	}
}
