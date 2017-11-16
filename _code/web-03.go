package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	fmt.Println("please connect to localhost:7777")
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/img", imageHandle)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, rootPage)
}

const rootPage = `<html>
<head>
	<title>Displaying images with Go</title>
</head>

<body>
	<h1>Image display</h1>
	<div id="content"><img src="/img"></img></div>
</body>
`

func imageHandle(w http.ResponseWriter, r *http.Request) {
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
	err := png.Encode(w, img)
	if err != nil {
		log.Printf("error encoding image: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
