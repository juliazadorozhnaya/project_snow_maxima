package main

import (
	"strconv"
)

// ==========точка============
type Dot struct {
	x, y, z int
}

func NewDot(x, y, z int) Dot {
	var d Dot
	d.x = x
	d.y = y
	d.z = z
	return d
}

func (a *Dot) DistanceXY(b Dot) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)
}

func (a *Dot) DistanceXYZ(b Dot) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y) + (a.z-b.z)*(a.z-b.z)
}

func (a Dot) String() string {
	return ".(" + strconv.Itoa(a.x) + ", " + strconv.Itoa(a.y) + ", " + strconv.Itoa(a.z) + ")"
}
