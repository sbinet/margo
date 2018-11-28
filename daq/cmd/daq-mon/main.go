// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sbinet/margo/daq"
	"golang.org/x/net/websocket"
	"gonum.org/v1/plot/vg"
)

var (
	evts = make(chan Event)

	mu  sync.RWMutex
	mon monitor
)

func main() {
	dialAddr := flag.String("dial", "localhost:8000", "address to dial and read data from")
	servAddr := flag.String("addr", ":8080", "address to serve")

	flag.Parse()

	go dial(*dialAddr)
	go run()

	http.HandleFunc("/", rootHandler)
	http.Handle("/ws", websocket.Handler(cmdHandler))
	http.HandleFunc("/plot-calo", plotCalo)
	http.HandleFunc("/plot-env", plotEnv)

	log.Fatal(http.ListenAndServe(*servAddr, nil))
}

type Event struct {
	Header daq.Header
	Envs   []daq.Env
	Calo   []daq.Calorimeter
}

func run() {
	for evt := range evts {
		mu.Lock()
		mon.update(evt)
		mu.Unlock()
	}
}

func dial(addr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not dial DAQ server: %v", err)
	}
	defer c.Close()

	dec := daq.NewDecoder(c)
	hdr := dec.ReadStreamHeader()
	if hdr != daq.MagicHdr {
		log.Fatalf("not a DAQ stream: %q", hdr)
	}

	sc := daq.NewEvtScanner(c)
	for sc.Scan() {
		var evt Event

		dec := daq.NewDecoder(bytes.NewReader(sc.Bytes()))
		blk := dec.ReadU32()
		if blk != daq.EvtBeg {
			log.Fatalf("expected a magic-block beg, got: %#x", blk)
		}

		evt.Header = dec.ReadHeader()

		buf := make([]byte, evt.Header.Len)
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
				// fmt.Printf("cal: %#v\n", calo)
				evt.Calo = append(evt.Calo, calo)
			case daq.EnvSysDet:
				env := dec.ReadEnvSys()
				// fmt.Printf("env: %#v\n", env)
				evt.Envs = append(evt.Envs, env)
			default:
				log.Fatalf("unknown detector type: %x", typ)
			}
		}
		if det.Err() != nil {
			log.Fatalf("sub error: %v", det.Err())
		}

		select {
		case evts <- evt:
		default:
			// nobody's ready, drop event on the floor.
		}
	}
	if sc.Err() != nil {
		log.Fatalf("scanner error: %v", sc.Err())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, rootPage)
}

func cmdHandler(c *websocket.Conn) {
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	ping := []byte("ping")
	for range tick.C {
		err := websocket.Message.Send(c, ping)
		if err != nil {
			log.Printf("could not send ping command: %v", err)
			return
		}
	}
}

func plotCalo(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	o, err := mon.Calo.WriterTo(10*vg.Centimeter, 10*vg.Centimeter, "png")
	if err != nil {
		log.Printf("could not create calo writer plot: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/png")
	o.WriteTo(w)
}

func plotEnv(w http.ResponseWriter, r *http.Request) {
	o, err := mon.Env.WriterTo(15*vg.Centimeter, 10*vg.Centimeter, "png")
	if err != nil {
		log.Printf("could not create env writer plot: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/png")
	o.WriteTo(w)
}

const rootPage = `<html>
<head>
        <title>DAQ monitor</title>
        <script type="text/javascript">
        var sock = null;

        function update(data) {
                var calo = document.getElementById("calo");
                calo.src = "/plot-calo?random="+new Date().getTime();

				var temp = document.getElementById("temp");
                temp.src = "/plot-env?random="+new Date().getTime();
        };

        window.onload = function() {
                sock = new WebSocket("ws://"+location.host+"/ws");
                sock.onmessage = function(event) {
                        update(event.data);
                };
        };
        </script>
</head>

<body>
	<h1>DAQ Monitor</h1>
	<div><img id="calo" src="" alt="N/A"/></div>
	<div><img id="temp" src="" alt="N/A"/></div>
</body>

</html>
`
