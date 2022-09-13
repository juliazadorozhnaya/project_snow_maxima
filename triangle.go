package main

import (
	"fmt"
)

// ==========треугольник============
type Triangle struct {
	a  Dot // left-bottom corner always
	b  Dot // next point by clockwise!
	c  Dot
	ab Edge
	bc Edge
	ca Edge
	s  float32 // площадь под треугольником
	v  float32
}

func NewTriangle(a, b, c Dot) Triangle {
	var t Triangle
	if a.x > b.x {
		a, b = b, a
	}
	if a.x > c.x {
		a, c = c, a
	}
	if b.x > c.x {
		b, c = c, b
	}

	if a.x == b.x {
		if a.y > b.y {
			a, b = b, a
		}
	}

	// обход по часовой стрелке исходя из знака площади
	// ab.x*ac.y - ab.y*ac.x
	// ab.x = b.x-a.x
	// ab.y = b.y-a.y
	// ac.x = c.x-a.x
	// ac.y = c.y-a.y
	t.s = float32((b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x))
	if t.s > 0 {
		b, c = c, b
	} else {
		t.s = -t.s
	}
	t.s = t.s / 2
	t.a = a
	t.b = b
	t.c = c

	t.ab = NewEdge(a, b)
	t.bc = NewEdge(b, c)
	t.ca = NewEdge(c, a)

	t.V()
	return t
}

// площадь под треугольником
func (t *Triangle) S() float32 {
	t.s = float32(t.ca.v.SinProduct(t.ab.v.Mul(-1)))
	t.s = t.s / 2
	// не должно происходить. удалить если правда не будет случаться
	if t.s < 0 {
		fmt.Println("ERROR: triangle's square < 0")
	}
	return t.s
}

// объём пирамиды вокруг треугольника
func (t *Triangle) V() float32 {
	det := func(a, b, c Vector) int {
		r1 := a.x * b.y * c.z
		r2 := a.y * b.z * c.x
		r3 := a.z * b.x * c.y
		p1 := a.z * b.y * c.x
		p2 := a.x * b.z * c.y
		p3 := a.y * b.x * c.z
		return (r1 + r2 + r3 - p1 - p2 - p3)
	}

	//упорядочивание по высоте
	a := t.a
	b := t.b
	c := t.c
	if a.z > b.z {
		a, b = b, a
	}
	if a.z > c.z {
		a, c = c, a
	}
	if b.z > c.z {
		b, c = c, b
	}

	d := Dot{x: b.x, y: b.y, z: a.z} // точка под средней по высоте
	e := Dot{x: c.x, y: c.y, z: a.z} // точка под верхней

	ba := NewVector(b, a)
	bc := NewVector(b, c)
	bd := NewVector(b, d)
	be := NewVector(b, e)

	// считаем объёмы. пока без math.Abs
	vh := det(ba, bc, be)
	vl := det(ba, bd, be)
	if vh < 0 {
		vh = -vh
	}
	if vl < 0 {
		vl = -vl
	}
	t.v = float32(vh+vl) / 6
	return t.v
}

// проверка на наличие точки внутри треугольника
func (t *Triangle) DotInside(d Dot) bool {
	if d == t.a {
		return false
	}
	if d == t.b {
		return false
	}
	if d == t.c {
		return false
	}

	//тут можно вместо создания вектора - математику расписать. Сэкономлю
	ad := NewVector(t.a, d)
	bd := NewVector(t.b, d)
	cd := NewVector(t.c, d)
	return (t.ab.v.SinProduct(ad) <= 0) && (t.bc.v.SinProduct(bd) <= 0) && (t.ca.v.SinProduct(cd) <= 0)
}

func (t Triangle) String() string {
	return "△ [A" + t.a.String() + ", B" + t.b.String() + ", C" + t.c.String() + "]"
}
