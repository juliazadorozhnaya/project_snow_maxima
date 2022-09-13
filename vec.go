package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// ==========Облако точек============
type PointCloud struct {
	d []Dot //массив точек
}

func NewPointCloud(x []int, y []int, z []int) *PointCloud {
	var p PointCloud
	p.d = make([]Dot, len(x))
	for i, xp := range x {
		p.d[i] = Dot{xp, y[i], z[i]}
	}
	sort.Slice(p.d, func(i, j int) bool {
		return p.d[i].x < p.d[j].x
	})
	return &p
}

// собственно триангуляция из точек
func (p *PointCloud) Triangulate(minLen int, updateScene func(edges []Edge)) *TriangleCloud {
	var trc TriangleCloud
	// надо расширить слайсы trc.e и trc.t сразу до максимально
	// возможного размера исходя из Эйлеровой характеристики
	trc.d = p.d
	trc.eMap = make(map[Edge]struct{})
	trc.tMap = make(map[Triangle]struct{})
	fmt.Println("Triangulation started")
	if !(trc.StartTriangle(minLen)) {
		fmt.Println("There is no any correct edge/triangle")
	}

	var e Edge
	var d Dot
	var f bool
	for {
		e, d, f = trc.GetNextTriangle(minLen)
		if f {
			trc.AppendTriangle(e, d)
			if len(trc.t)%100 == 0 {
				updateScene(trc.e)
				fmt.Println("Added triangle", len(trc.t), "with dot", d, ":", NewTriangle(e.b, e.e, d))
			}
		} else {
			break
		}
	}
	updateScene(trc.e)
	fmt.Println("Added triangle", len(trc.t), "with dot", d, ":", NewTriangle(e.b, e.e, d))
	return &trc
}

func (p PointCloud) String() string {
	var s string
	for i, d := range p.d {
		s += strconv.Itoa(i) + ") " + d.String() + "\n"
	}
	return s
}

// ==========Результат триангуляции по точкам============
type TriangleCloud struct {
	t []Triangle //массив треугольников
	e []Edge     //массив рёбер. Возможно не нужен будет
	d []Dot      //собственно исходные точки

	eMap map[Edge]struct{}
	tMap map[Triangle]struct{}
}

// инициализация рёбер. первое ребро - 2 ближайшие точки
func (trc *TriangleCloud) StartEdge(minLen int) bool {
	trc.e = trc.e[:0]
	var resultDot, closestDot Dot
	var currentDistance int
	var resultBool bool
	for leftBound, leftDot := range trc.d[:len(trc.d)-1] {
		for _, tmpDot := range trc.d[leftBound+1:] {
			currentDistance = trc.d[0].DistanceXY(tmpDot)
			if currentDistance < minLen {
				minLen = currentDistance
				resultDot = leftDot
				closestDot = tmpDot
				resultBool = true
			} else if (tmpDot.x-leftDot.x)*(tmpDot.x-leftDot.x) > minLen {
				break // надо по х проверять
			}
		}
	}
	if resultBool {
		fmt.Println("Closest edge found between", resultDot, "and", closestDot)
		trc.AppendEdge(NewEdge(resultDot, closestDot))
	}
	return resultBool
}

// добавляет новое ребро в список рёбер (без дубликатов)
func (trc *TriangleCloud) AppendEdge(e Edge) bool {
	_, edgeInside := trc.eMap[e]
	if !edgeInside {
		trc.eMap[e] = struct{}{}
		trc.e = append(trc.e, e)
		sort.Slice(trc.e, func(i, j int) bool {
			return trc.e[i].b.x < trc.e[j].b.x
		})
	}
	return edgeInside
}

// удаляет ребро из множества рёбер внешней оболочки
func (trc *TriangleCloud) RemoveEdge(e Edge) bool {
	_, edgeInside := trc.eMap[e]
	if edgeInside {
		delete(trc.eMap, e)
		lenOfEdges := len(trc.e)
		pos, _ := sort.Find(lenOfEdges, func(in int) int { return e.b.x - trc.e[in].b.x })
		for ; pos < lenOfEdges; pos++ {
			if trc.e[pos] == e {
				break
			}
		}
		if pos < lenOfEdges {
			trc.e = append(trc.e[:pos], trc.e[pos+1:]...)
			sort.Slice(trc.e, func(i, j int) bool {
				return trc.e[i].b.x < trc.e[j].b.x
			})
		} else {
			edgeInside = false
		}
	}
	return edgeInside
}

// инициализация треугольников и рёбер.
func (trc *TriangleCloud) StartTriangle(minLen int) bool {
	trc.t = trc.t[:0]
	var resultBool bool
	var closestDot Dot
	//поиск ближайшей точки к ребру, если такое есть
	//предусмотреть перевыбор стартовой точки если не найдётся треугольник
	if trc.StartEdge(minLen) {
		//var currentDistance int
		currentEdge := trc.e[0]
		// можно оптимизировать, уменьшив количество проверяемых точек длинной по 0Х
		// а лучше вообще заменить на GetClosestFreePointToEdge
		closestDot, minLen, resultBool = trc.GetClosestFreePointToEdge(currentEdge, minLen)
		// for _, tmpDot := range trc.d {
		// 	currentDistance = currentEdge.L2normXY(tmpDot)
		// 	if currentDistance >= minLen {
		// 		continue
		// 	}
		// 	if currentEdge.DotOMW(tmpDot) {
		// 		continue
		// 	}
		// 	minLen = currentDistance
		// 	closestDot = tmpDot
		// 	resultBool = true
		// }
	}
	if resultBool {
		trc.AppendTriangle(trc.e[0], closestDot)
	}
	return resultBool
}

// добавляет новый треугольник в триангуляцию и занимается рёбрами
// Добработаю чтобы ещё учитывало внешнюю оболочку (hull) из рёбер
func (trc *TriangleCloud) AppendTriangle(e Edge, d Dot) bool {
	t := NewTriangle(e.b, e.e, d)

	_, triangleInside := trc.tMap[t]
	if !triangleInside {
		trc.tMap[t] = struct{}{}
		trc.t = append(trc.t, t)
		trc.AppendEdge(t.ab)
		trc.AppendEdge(t.bc)
		trc.AppendEdge(t.ca)
	}
	return triangleInside

}

// проверка на пересечение отрезков внутренними частями, без учёта точек на краю
func (trc *TriangleCloud) CheckInnerCrossing(e Edge, left_bound int, right_bound int) bool {
	for _, j := range trc.e[left_bound:right_bound] {
		if e.InnerCrossing(&j) {
			return true
		}
	}
	return false
}

// проверка на наличие точек внутри имеющихся треугольников. Возможно получится без неё
func (trc *TriangleCloud) CheckDotsInside(t Triangle, left_bound int, right_bound int) bool {
	for _, d := range trc.d[left_bound:right_bound] {
		if t.DotInside(d) {
			return true
		}
	}
	return false
}

// поиск ближайшей точки к ребру, но не образующую уже имеющийся треугольник
func (trc *TriangleCloud) GetClosestFreePointToEdge(j Edge, minLen int) (Dot, int, bool) {
	bound := int(math.Ceil(math.Sqrt(float64(minLen))))
	lenOfDots := len(trc.d)
	lenOfEdges := len(trc.e)
	l_coord, _ := sort.Find(lenOfDots, func(in int) int { return j.b.x - bound - trc.d[in].x })
	r_coord, _ := sort.Find(lenOfDots, func(in int) int { return j.e.x + bound + 1 - trc.d[in].x })

	var tmpMinLen, tmpLeftBound, tmpRightBound int
	var tmpClosestPoint Dot
	var dotFound bool
	for _, d := range trc.d[l_coord:r_coord] {
		tmpMinLen = j.L2normXY(d)
		if tmpMinLen > minLen {
			continue
		}
		t := NewTriangle(j.b, j.e, d) //вынести из цикла чтоб не вызывать аллокацию памяти в конструкторе
		_, triangleInside := trc.tMap[t]
		if triangleInside {
			continue
		}
		//проверка что точка на не одной прямой с ребром
		if j.DotOMW(d) {
			continue
		}
		bd := NewEdge(j.b, d) //вынести из цикла чтоб не вызывать аллокацию памяти в конструкторе
		ed := NewEdge(j.e, d) //вынести из цикла чтоб не вызывать аллокацию памяти в конструкторе
		//проверка на пересечение будущих рёбер с существующими рёбрами
		//тут бы переделать на оболочку вместо диапазонов упорядоченного массива
		tmpLeftBound, _ = sort.Find(lenOfEdges, func(in int) int { return bd.b.x - bound - trc.e[in].b.x })
		tmpRightBound, _ = sort.Find(lenOfEdges, func(in int) int { return bd.e.x + 1 - trc.e[in].b.x })
		if trc.CheckInnerCrossing(bd, tmpLeftBound, tmpRightBound) {
			continue
		}
		tmpLeftBound, _ = sort.Find(lenOfEdges, func(in int) int { return ed.b.x - bound - trc.e[in].b.x })
		tmpRightBound, _ = sort.Find(lenOfEdges, func(in int) int { return ed.e.x + 1 - trc.e[in].b.x })
		if trc.CheckInnerCrossing(ed, tmpLeftBound, tmpRightBound) {
			continue
		}
		//тут проверка на включение точки внутрь треугольника
		//с учётом рёбер, без учёта вершин
		tmpLeftBound, _ := sort.Find(lenOfDots, func(in int) int { return t.a.x - trc.d[in].x })
		tmpRightBound, _ := sort.Find(lenOfDots, func(in int) int { return t.c.x + 1 - trc.d[in].x })
		if trc.CheckDotsInside(t, tmpLeftBound, tmpRightBound) {
			continue
		}

		minLen = tmpMinLen
		tmpClosestPoint = d
		dotFound = true
	}
	return tmpClosestPoint, minLen, dotFound
}

// minLen - сумма квадратов длинн от краев ребра до точки
func (trc *TriangleCloud) GetNextTriangle(minLen int) (Edge, Dot, bool) {
	var tmpClosestPoint Dot
	var tmpMinLen int
	var tmpDotFound bool
	var resultDot Dot
	var resultEdge Edge
	var triangleFound bool
	for _, e := range trc.e {
		//fmt.Println("Analyzed edge", e)
		tmpClosestPoint, tmpMinLen, tmpDotFound = trc.GetClosestFreePointToEdge(e, minLen)
		if tmpDotFound {
			if tmpMinLen < minLen {
				minLen = tmpMinLen
				resultDot = tmpClosestPoint
				resultEdge = e
				triangleFound = true
			}
		}
	}
	return resultEdge, resultDot, triangleFound
}

// подсчёт объёма под триангулированной поверхностью
func (trc *TriangleCloud) GetVolume() float32 {
	var result float32
	for _, t := range trc.t {
		result += t.v + float32(t.a.z)*t.s
	}
	return result
}
