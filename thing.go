package main

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	fmt "github.com/gosve/react/forks/fmtless"
	"image"
	"image/color"
	"io/ioutil"
)

type Thing struct {
	Name     string
	Type     string
	out      interface{}
	pressed  bool
	selected bool
}

func (t *Thing) Info(col color.RGBA, gtx *layout.Context) {
	th := material.NewTheme()
	th.H6(t.Name).Layout(gtx)
}

func (t *Thing) Layout(f string, col color.RGBA, gtx *layout.Context) {
	cs := gtx.Constraints
	for _, e := range gtx.Events(t) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				t.pressed = true
				t.selected = true
				readFile(f)
				fmt.Println(f)
			case pointer.Release:
				t.pressed = false
			}
		}
	}
	th := material.NewTheme()

	if t.pressed {
		col = color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}
	}
	pointer.Rect(
		image.Rectangle{Max: image.Point{X: cs.Width.Max, Y: cs.Height.Max}},
	).Add(gtx.Ops)
	pointer.InputOp{Key: t}.Add(gtx.Ops)
	square := f32.Rectangle{
		Max: f32.Point{X: float32(cs.Width.Max), Y: float32(cs.Height.Max)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: square}.Add(gtx.Ops)
	th.H6(t.Name).Layout(gtx)
}

func readFile(f string) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
	}
	editor.SetText(string(dat))
}
func out(i []byte) string {
	return "[" + string(i) + "]$"
}

func fill(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
