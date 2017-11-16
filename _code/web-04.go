package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var (
	datac = make(chan string)
)

func main() {
	fmt.Println("please connect to localhost:7777")
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/img", imageHandle)
	http.Handle("/chan", websocket.Handler(chanHandler))
	go generate(datac)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, rootPage)
}

const rootPage = `<html>
<head>
	<title>Displaying images with Go</title>
	<script type="text/javascript">
	var sock = null;

	function update(data) {
		var img = document.getElementById("img-node");
		img.src = "data:image/png;base64,"+data;
	};

	window.onload = function() {
		sock = new WebSocket("ws://localhost:7777/chan");
		sock.onmessage = function(event) {
			update(event.data);
		};
	};
	</script>
</head>

<body>
	<h1>Image display</h1>
	<div id="content"><img id="img-node" src="" alt="N/A"/></div>
</body>
`

func imageHandle(w http.ResponseWriter, r *http.Request) {
	img := newImage()
	err := png.Encode(w, img)
	if err != nil {
		log.Printf("error encoding image: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newImage() image.Image {
	const sz = 50
	// create the whole image canvas
	canvas := image.Rect(0, 0, 2*sz, 2*sz)
	img := image.NewRGBA(canvas)
	// draw a gray background
	draw.Draw(img, canvas, image.NewUniform(color.RGBA{0x66, 0x66, 0x66, 0xff}), image.ZP, draw.Src)
	// create a randomly sized, randomly centered, red square
	x1 := rand.Intn(sz)
	y1 := rand.Intn(sz)
	x2 := rand.Intn(sz) + sz
	y2 := rand.Intn(sz) + sz
	draw.Draw(img, image.Rect(x1, y1, x2, y2), image.NewUniform(color.RGBA{0xff, 0x00, 0x00, 0xff}), image.ZP, draw.Src)
	return img
}

func chanHandler(ws *websocket.Conn) {
	for data := range datac {
		err := websocket.Message.Send(ws, data)
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}

func generate(datac chan string) {
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for range tick.C {
		buf := new(bytes.Buffer)
		err := png.Encode(buf, newImage())
		if err != nil {
			log.Fatal(err)
		}
		datac <- base64.StdEncoding.EncodeToString(buf.Bytes())
	}
}
