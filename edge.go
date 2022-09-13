package main

// ==========отрезок============
type Edge struct {
	b Dot // left-bottom corner always!
	e Dot
	v Vector
}

func NewEdge(b, e Dot) Edge {
	var edg Edge
	edg.v = NewVector(b, e)
	if b.x < e.x {
		edg.b = b
		edg.e = e
		return edg
	}

	if b.x > e.x {
		edg.b = e
		edg.e = b
		return edg
	}

	if b.y < e.y {
		edg.b = b
		edg.e = e
	} else {
		edg.b = e
		edg.e = b
	}
	return edg
}

// сумма квадратов длин от краёв отрезка до точки на проекции в 0XY
func (e *Edge) L2normXY(d Dot) int {
	return e.b.DistanceXY(d) + e.e.DistanceXY(d)
}

// сумма квадратов длин от краёв отрезка в 3D
func (e *Edge) L2normXYZ(d Dot) int {
	return e.b.DistanceXYZ(d) + e.e.DistanceXYZ(d)
}

// квадрат длины отрезка, спроецированного в плоскость 0XY
func (e *Edge) L2lenXY() int {
	return e.b.DistanceXY(e.e)
}

// квадрат длины отрезка в 3D
func (e *Edge) L2lenXYZ() int {
	return e.b.DistanceXYZ(e.e)
}

// проверка на нахождение точки на одной прямой с отрезком
func (e *Edge) DotOMW(d Dot) bool {
	db := NewVector(d, e.b)
	de := NewVector(d, e.e)
	return db.SinProduct(de) == 0
}

// проверка на пересечение двух отрезков внутренними частями, без учёта точек на краю
func (e *Edge) InnerCrossing(j *Edge) bool {
	ac := NewVector(e.b, j.b)
	ad := NewVector(e.b, j.e)
	ca := NewVector(j.b, e.b) // подумать как от нее избавиться. Это же -1*ac
	cb := NewVector(j.b, e.e)
	S1 := e.v.SinProduct(ac) * e.v.SinProduct(ad)
	S2 := j.v.SinProduct(ca) * j.v.SinProduct(cb)
	if (S1 < 0) && (S2 < 0) {
		return true
	}
	return false
}

func (e Edge) String() string {
	return "_[" + e.b.String() + ", " + e.e.String() + "]"
}
