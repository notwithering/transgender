package main

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	screen, err := capture(display)
	if err != nil {
		fmt.Println(err)
		return
	}

	boxes, err := getBoxes(screen)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, box := range boxes {
		if box.Confidence >= threshold {
			fmt.Println(box.Word)
		}
	}

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_POPUP)
	if err != nil {
		fmt.Println(err)
		return
	}

	win.Connect("destroy", gtk.MainQuit)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.ShowAll()
	gtk.Main()
}

func capture(display int) ([]byte, error) {
	bounds := screenshot.GetDisplayBounds(display)

	img, err := screenshot.Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func getBoxes(img []byte) (out []gosseract.BoundingBox, err error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(img)

	return client.GetBoundingBoxes(gosseract.RIL_WORD)
}
