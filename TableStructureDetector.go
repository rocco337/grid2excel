package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/lucasb-eyer/go-colorful"
	"gocv.io/x/gocv"
)

//TableStructureDetector ...
type TableStructureDetector struct {
}

//Detect ...
func (t *TableStructureDetector) Detect(filename string) (TableStructure, error) {

	return t.detectEdgesInternal(filename), nil
}

func (t *TableStructureDetector) detectEdgesInternal(filename string) TableStructure {
	mat := gocv.IMRead(filename, gocv.IMReadColor)

	matCanny := gocv.NewMat()
	matGray := gocv.NewMat()
	matLines := gocv.NewMat()

	window := gocv.NewWindow("detected lines")

	gocv.CvtColor(mat, &matGray, gocv.ColorBGRToGray)

	gocv.Canny(matGray, &matCanny, 100, 150)
	gocv.HoughLinesPWithParams(matCanny, &matLines, 1, math.Pi/180, 200, 2, 4)

	rows := make([]gocv.Veci, 0)
	columns := make([]gocv.Veci, 0)
	for i := 0; i < matLines.Rows(); i++ {
		line := matLines.GetVeciAt(i, 0)

		if line[1] > line[3] {
			//vertical
			columns = append(columns, line)
		} else if line[0] < line[2] {
			//horizontal
			rows = append(rows, line)
		}
	}

	result := TableStructure{Rows: make([][]Column, len(rows))}

	for rowIndex, row := range rows {
		result.Rows[rowIndex] = make([]Column, len(columns))

		for columnIndex, column := range columns {
			x, y, err := t.intersection(columns[columnIndex], rows[rowIndex])
			if err != nil {
				log.Println("Cannot find intersection", column, row)
				continue
			}

			point := []int32{x, y}

			result.Rows[rowIndex][columnIndex].Points = append(result.Rows[rowIndex][columnIndex].Points, point)

			//set intersection as end of previous column
			if columnIndex-1 >= 0 {
				result.Rows[rowIndex][columnIndex-1].Points = append(result.Rows[rowIndex][columnIndex-1].Points, point)
			}

			//set column intersection to previous row
			if rowIndex-1 >= 0 {
				result.Rows[rowIndex-1][columnIndex].Points = append(result.Rows[rowIndex-1][columnIndex].Points, point)
			}
		}

		if rowIndex == 1 {
			break

		}
	}

	// intersections := make([]gocv.Veci, 0)
	// for _, column := range columns {
	// 	for _, row := range rows {

	// 		x, y, err := t.intersection(column, row)
	// 		if err != nil {
	// 			log.Println("Cannot find intersection", column, row)
	// 		} else {
	// 			intersections = append(intersections, gocv.Veci{x, y})
	// 			result.Rows = append(result.Rows, Column{})
	// 		}

	// 	}
	// }

	for _, row := range result.Rows {
		fmt.Println("--------------row--------------")
		for _, columns := range row {
			fmt.Println("--------------column--------------")
			randomColor := colorful.Hcl(180.0+rand.Float64()*50.0, 0.2+rand.Float64()*0.8, 0.3+rand.Float64()*0.7)

			for _, point := range columns.Points {

				fmt.Println(point)
				pt1 := image.Pt(int(point[0]), int(point[1]))
				gocv.Line(&mat, pt1, pt1, color.RGBA{uint8(randomColor.R), uint8(randomColor.G), uint8(randomColor.B), 50}, 4)
				// gocv.PutText(&mat, fmt.Sprint(pt1), pt1, gocv.FontHersheySimplex, 0.4, color.RGBA{0, 255, 0, 50}, 2)
			}
		}
	}

	for {
		window.ResizeWindow(15000, 20000)
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}

	return result
}

//copied from https://stackoverflow.com/questions/20677795/how-do-i-compute-the-intersection-point-of-two-lines
func (t *TableStructureDetector) intersection(line1, line2 gocv.Veci) (int32, int32, error) {
	xdiff := []int32{line1[0] - line1[2], line2[0] - line2[2]}
	ydiff := []int32{line1[1] - line1[3], line2[1] - line2[3]}

	det := func(a, b []int32) int32 {
		return a[0]*b[1] - a[1]*b[0]
	}

	div := det(xdiff, ydiff)
	if div == 0 {
		return 0, 0, errors.New("Cannnot find interseciton")
	}

	d := []int32{det([]int32{line1[0], line1[1]}, []int32{line1[2], line1[3]}), det([]int32{line2[0], line2[1]}, []int32{line2[2], line2[3]})}
	x := det(d, xdiff) / div
	y := det(d, ydiff) / div

	return x, y, nil
}

//TableStructure ...
type TableStructure struct {
	Rows [][]Column
}

//Column ...
type Column struct {
	Points [][]int32
}
