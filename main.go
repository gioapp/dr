package main

import (
	"flag"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gop9/gorminal/mod"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	editor   = new(widget.Editor)
	mainList = &layout.List{
		Axis: layout.Horizontal,
	}
	listFolder = &layout.List{
		Axis: layout.Vertical,
	}
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	textColor = color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}
)

func main() {
	com := mod.CommandsHistory{}

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	gofont.Register()
	th := material.NewTheme()
	col := color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}
	flag.Parse()
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(1400), unit.Dp(800)),
			app.Title("ParallelCoin"),
		)
		things := make(map[string]*Thing)
		listThings, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		for _, t := range listThings {
			things[t.Name()] = &Thing{
				Name: t.Name(),
			}
			if t.IsDir() {
				things[t.Name()].Type = "space"
			}
		}
		for _, t := range listThings {
			things[t.Name()] = &Thing{
				Name: t.Name(),
			}
		}
		gtx := layout.NewContext(w.Queue())
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				gtx.Reset(e.Config, e.Size)
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(0.765, func() {
						widgets := []func(){
							func() {
								layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
									fill(gtx)

									listFolder.Layout(gtx, len(listThings), func(i int) {
										if listThings[i].IsDir() {
											col := color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0x30}
											things[listThings[i].Name()].Layout(pwd+"/"+listThings[i].Name(), col, gtx)
										} else {
											things[listThings[i].Name()].Layout(pwd+"/"+listThings[i].Name(), col, gtx)
										}
									})
								})
							},

							func() {
								layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
									th.Editor("").Layout(gtx, editor)
								})
							},
						}
						mainList.Layout(gtx, len(widgets), func(i int) {
							layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
						})

					}),
					layout.Flexed(0.235, func() {

						fill(gtx)
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, func() {
								layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
									layout.Flex{
										Axis:    layout.Vertical,
										Spacing: layout.SpaceAround,
									}.Layout(gtx,
										layout.Flexed(1, func() {
											consoleOutputList.Layout(gtx, len(com.Commands), func(i int) {
												t := com.Commands[i]
												layout.Flex{
													Alignment: layout.End,
												}.Layout(gtx,
													layout.Rigid(func() {
														out := th.Body1(t.Out)
														//out.Font.Size = unit.Dp(fontSize)
														out.Color = textColor
														out.Layout(gtx)
													}),
												)
											})
										}),
										layout.Rigid(func() {
											pwd, _ := exec.Command("pwd").Output()
											layout.Flex{}.Layout(gtx,
												layout.Rigid(func() {
													p := th.Body1(out(pwd))
													p.Font.Style = text.Regular
													//p.Font.Size = unit.Dp(fontSize)
													p.Color = textColor
													p.Layout(gtx)
												}),
												layout.Rigid(func() {
													input := th.Editor("")
													input.Font.Style = text.Regular
													//input.Font.Size = unit.Dp(fontSize)
													input.Color = textColor
													input.Layout(gtx, consoleInputField)
												}),
											)
											for _, e := range consoleInputField.Events(gtx) {
												if e, ok := e.(widget.SubmitEvent); ok {
													splitted := strings.Split(e.Text, " ")
													cmd, _ := exec.Command(splitted[0], splitted[1:]...).CombinedOutput()
													com.Commands = append(com.Commands, mod.Command{
														ComID: e.Text,
														Time:  time.Time{},
														Out:   out(cmd),
													})
													consoleInputField.SetText("")
												}
											}
										}))
								})
							}),
						)

					}),
				)

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
