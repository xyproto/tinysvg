package mtransform

import "math"

type Transform [3][3]float64

func (t *Transform) Apply(x float64, y float64) (float64, float64) {
	var X, Y float64
	X = t[0][0]*x + t[0][1]*y + t[0][2]
	Y = t[1][0]*x + t[1][1]*y + t[1][2]
	return X, Y
}

func Identity() Transform {
	var t Transform
	t[0][0] = 1
	t[1][1] = 1
	t[2][2] = 1
	return t
}
func NewTransform() *Transform {
	var t Transform
	t = Identity()
	return &t
}

func MultiplyTransforms(a Transform, b Transform) Transform {
	return Transform{
		{
			a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0],
			a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1],
			a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2],
		},
		{
			a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0],
			a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1],
			a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2],
		},
		{
			a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0],
			a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1],
			a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2],
		},
	}
}

func (a *Transform) MultiplyWith(b Transform) {
	*a = MultiplyTransforms(*a, b)
}

func (t *Transform) Scale(x float64, y float64) {
	a := Identity()
	a[0][0] = x
	a[1][1] = y
	t.MultiplyWith(a)
}
func (t *Transform) Translate(x float64, y float64) {
	a := Identity()

	a[0][2] = x
	a[1][2] = y
	t.MultiplyWith(a)
}

func (t *Transform) RotateOrigin(angle float64) {
	a := Identity()
	a[0][0] = math.Cos(angle)
	a[0][1] = -math.Sin(angle)
	a[1][0] = math.Sin(angle)
	a[1][1] = a[0][0]
	t.MultiplyWith(a)
}

func (t *Transform) RotatePoint(angle float64, x float64, y float64) {
	t.Translate(x, y)
	t.RotateOrigin(angle)
	t.Translate(-x, -x)
}

func (t *Transform) SkewX(angle float64) {
	a := Identity()
	a[0][1] = math.Tan(angle)
	t.MultiplyWith(a)
}

func (t *Transform) SkewY(angle float64) {
	a := Identity()
	a[1][0] = math.Tan(angle)
	t.MultiplyWith(a)
}
