package tinysvg

import (
	"testing"
)

func TestSVG(t *testing.T) {
	document, svg := NewTinySVG(256, 256)
	svg.Describe("Diagram")

	roundedRectangle := svg.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")

	document.SaveSVG("/tmp/output.svg")
}
