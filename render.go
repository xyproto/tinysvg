package tinysvg

import (
	"bytes"
	"fmt"
	"github.com/rustyoz/svg"
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

	go func() {
		for err := range errorChan {
			if err != nil {
				fmt.Println("error", err)
			} else {
				fmt.Println("no error")
			}
		}
	}()

	go func() {
		for drawingInstruction := range drawingInstructionChan {
			fmt.Println("drawing instruction", drawingInstruction)
		}
	}()

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{256, 256}})
	return img, nil
}
