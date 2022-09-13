package main

import (
	"strconv"
)

// ==========вектор============
type Vector Dot

func NewVector(b Dot, e Dot) Vector {
	var v Vector
	v.x = e.x - b.x
	v.y = e.y - b.y
	v.z = e.z - b.z
	return v
}

func (v Vector) Len2() int {
	return v.x*v.x + v.y*v.y
}

// scalar product
func (a Vector) CosProduct(b Vector) int {
	return a.x*b.x + a.y*b.y
}

// vector product
func (a Vector) SinProduct(b Vector) int {
	return a.x*b.y - a.y*b.x
}

func (a Vector) Add(b Vector) Vector {
	var v Vector
	v.x = a.x + b.x
	v.y = a.y + b.y
	v.z = a.z + b.z
	return v
}

func (a Vector) Sub(b Vector) Vector {
	var v Vector
	v.x = a.x - b.x
	v.y = a.y - b.y
	v.z = a.z - b.z
	return v
}

func (a Vector) Mul(k int) Vector {
	var v Vector
	v.x = a.x * k
	v.y = a.y * k
	v.z = a.z * k
	return v
}

func (a Vector) String() string {
	return "→(" + strconv.Itoa(a.x) + ", " + strconv.Itoa(a.y) + ", " + strconv.Itoa(a.z) + ")"
}
