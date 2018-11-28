// Copyright 2018 The margo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/sbinet/margo/daq"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type monitor struct {
	Calo *plot.Plot
	Env  *plot.Plot
	envs map[string]plotter.XYs
}

func (mon *monitor) update(evt Event) {
	if mon.envs == nil {
		mon.envs = make(map[string]plotter.XYs)
	}
	mon.calo(evt.Calo)
	mon.env(evt.Envs)
}

func (mon *monitor) calo(calo []daq.Calorimeter) {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Calorimeter"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	pal := palette.Heat(12, 1)
	grd := newGrid(calo)
	h := plotter.NewHeatMap(grd, pal)
	p.Add(h)
	p.Add(plotter.NewGrid())

	mon.Calo = p
}

func (mon *monitor) env(envs []daq.Env) {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Temperature"
	p.X.Label.Text = "Time"
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "2006-01-02\n15:04:05",
	}
	p.Y.Min = 10

	now := float64(time.Now().UTC().Unix())
	for _, e := range envs {
		data := mon.envs[e.Name]
		data = append(data, struct{ X, Y float64 }{
			X: now,
			Y: e.T,
		})
		if len(data) > 1024 {
			copy(data, data[512:])
			data = data[:512]
		}
		mon.envs[e.Name] = data
	}

	for k, v := range mon.envs {
		lines, points, err := plotter.NewLinePoints(v)
		if err != nil {
			log.Fatal(err)
		}
		c := plotColors[k]
		points.Color = c
		lines.Color = c
		p.Add(points, lines)
		p.Legend.Add(k, lines)
	}

	p.Add(plotter.NewGrid())

	mon.Env = p
}

var (
	plotColors = make(map[string]color.Color)
)

func init() {
	for i, c := range plotutil.SoftColors {
		plotColors[fmt.Sprintf("raspi-%02d", i)] = c
	}
}

type grid struct {
	data map[key]float64
}

func (grid) Dims() (int, int) { return N, N }
func (g grid) Z(c, r int) float64 {
	k := key{c: int64(c), r: int64(r)}
	ene, ok := g.data[k]
	if !ok {
		return math.NaN()
	}
	return ene
}
func (g grid) X(c int) float64 { return float64(c) }
func (g grid) Y(r int) float64 { return float64(r) }

func newGrid(calo []daq.Calorimeter) grid {
	g := grid{data: make(map[key]float64)}
	for _, c := range calo {
		g.data[keyFrom(int64(c.CellID))] = c.Ene
	}
	return g
}

type key struct {
	c, r int64
}

const N = 16

func keyFrom(cell int64) key {
	r := cell / N
	c := cell - r*N
	return key{c: c, r: r}
}
