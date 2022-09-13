package main

import "fmt"

func main1() {
	//code for tests
	dota := NewDot(11, 0, 0)
	dotb := NewDot(0, 11, 0)
	dotc := NewDot(0, 0, 0)
	dotd := NewDot(2, 2, 0)
	t := NewTriangle(dota, dotb, dotc)
	fmt.Println(t.a, t.b, t.c)
	fmt.Println(t.s)
	fmt.Println(t.v)
	fmt.Println(t.DotInside(dotd))
	mm := map[string]bool{
		"aaa": true,
		"aab": true,
		"bbb": false,
	}
	fmt.Println(mm)
	delete(mm, "aaa")
	fmt.Println(mm)
	delete(mm, "aaa")
	fmt.Println(mm)

	//code for triangulate
	pcd := NewPointCloud(GetScanData(-1))
	fmt.Println(pcd)
	triangulatedData := pcd.Triangulate(8, func(edges []Edge) {})
	fmt.Println(len(triangulatedData.d), ":", triangulatedData.d)
	fmt.Println(len(triangulatedData.e), ":", triangulatedData.e)
	fmt.Println(len(triangulatedData.t), ":", triangulatedData.t)
}
