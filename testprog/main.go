package main

import (
	"bytes"
	"fmt"
	"github.com/rustyoz/svg"
	"github.com/xyproto/tinysvg"
	"image"
)

func Render(svgdata []byte) (*image.RGBA, error) {
	p, err := svg.ParseSvgFromReader(bytes.NewReader(svgdata), "Untitled", 1.0)
	if err != nil {
		return nil, err
	}
	fmt.Println("title:", p.Title)
	fmt.Println("groups:")
	for _, group := range p.Groups {
		fmt.Println("group:", group)
	}
	fmt.Println("width:", p.Width)
	fmt.Println("height:", p.Height)
	fmt.Println("viewbox", p.ViewBox)
	for _, element := range p.Elements {
		fmt.Println("element:", element)
	}
	fmt.Println("name:", p.Name)
	fmt.Println("transform:", p.Transform)

	drawingInstructionChan, errorChan := p.ParseDrawingInstructions()

	for err := range errorChan {
		if err != nil {
			fmt.Println("error", err)
		} else {
			fmt.Println("no error")
		}
	}

	for drawingInstruction := range drawingInstructionChan {
		fmt.Println("drawing instruction", drawingInstruction)
	}

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{256, 256}})
	return img, nil
}

func main() {
	document, svgTag := tinysvg.NewTinySVG(256, 256)
	svgTag.Describe("Diagram")

	roundedRectangle := svgTag.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")

	fmt.Println(document)

	_, err := Render(document.Bytes())
	if err != nil {
		panic(err)
	}
}
