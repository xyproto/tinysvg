package tinysvg

import (
	"bytes"
	"testing"
)

func TestSVG(t *testing.T) {
	document, svg := NewTinySVG(256, 256)
	svg.Describe("Diagram")
	roundedRectangle := svg.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")
	document.SaveSVG("/tmp/output.svg")
}

func TestString(t *testing.T) {
	document, svg := NewTinySVG(256, 256)
	svg.Describe("Diagram")
	roundedRectangle := svg.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")
	s := document.String()
	if len(s) != 258 {
		t.Fatalf("1: length is not 258 but %d\n", len(s))
	}
	s = document.String()
	if len(s) != 258 {
		t.Fatalf("2: length is not 258 but %d\n", len(s))
	}
}

func TestWriteTo(t *testing.T) {
	document, svg := NewTinySVG(256, 256)
	svg.Describe("Diagram")
	roundedRectangle := svg.AddRoundedRect(30, 10, 5, 5, 20, 20)
	roundedRectangle.Fill("red")
	var buf bytes.Buffer
	document.WriteTo(&buf)
	s := buf.String()
	if len(s) != 258 {
		t.Fatalf("1: length is not 258 but %d\n", len(s))
	}
	var buf2 bytes.Buffer
	document.WriteTo(&buf2)
	s = buf2.String()
	if len(s) != 258 {
		t.Fatalf("2: length is not 258 but %d\n", len(s))
	}
}
