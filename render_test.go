package tinysvg

import (
	"fmt"
	"testing"
)

func TestRender(t *testing.T) {
	document, svg := NewTinySVG(256, 256)
	svg.Describe("Diagram")

	roundedRectangle := svg.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")

	fmt.Println(document)

	_, err := Render(document.Bytes())
	if err != nil {
		t.Error(err)
	}

	//fmt.Println(img)
}
