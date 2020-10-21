package main

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math"
	"sort"

	"gocv.io/x/gocv"
)

//TableStructureDetector ...
type TableStructureDetector struct {
}

//Detect ...
func (t *TableStructureDetector) Detect(filename string) (TableStructure, gocv.Mat, error) {
	return t.detectEdgesInternal(filename)
}

func (t *TableStructureDetector) detectEdgesInternal(filename string) (TableStructure, gocv.Mat, error) {
	mat := gocv.IMRead(filename, gocv.IMReadColor)

	matCanny := gocv.NewMat()
	matGray := gocv.NewMat()
	matLines := gocv.NewMat()

	gocv.CvtColor(mat, &matGray, gocv.ColorBGRToGray)

	gocv.Canny(matGray, &matCanny, 100, 150)
	gocv.HoughLinesPWithParams(matCanny, &matLines, 1, math.Pi/180, 200, 2, 4)

	// remove detected line - paint them white  should work?
	removeLine := func(line []int) {
		pt1 := image.Pt(int(line[0]), int(line[1]))
		pt2 := image.Pt(int(line[2]), int(line[3]))
		white := color.RGBA{255, 255, 255, 255}
		gocv.Line(&mat, pt1, pt2, white, 1)
	}

	rows := make([][]int, 0)
	columns := make([][]int, 0)
	for i := 0; i < matLines.Rows(); i++ {
		line := matLines.GetVeciAt(i, 0)
		lineAsIntArray := []int{int(line[0]), int(line[1]), int(line[2]), int(line[3])}

		if line[1] > line[3] {
			//vertical
			columns = appendSortedBy(columns, lineAsIntArray, 0)

		} else if line[0] < line[2] {
			//horizontal
			rows = appendSortedBy(rows, lineAsIntArray, 1)
		}
		removeLine(lineAsIntArray)

	}

	//remove duplicates
	for rowIndex, row := range rows {
		if rowIndex > 1 && row[1]-rows[rowIndex-1][1] < 5 {
			copy(rows[rowIndex:], rows[rowIndex+1:])
		}
	}
	for columnIndex, column := range columns {
		if columnIndex > 1 && column[1]-columns[columnIndex-1][1] < 5 {
			copy(columns[columnIndex:], columns[columnIndex+1:])
		}
	}

	result := TableStructure{Rows: make([][]Column, len(rows))}

	for rowIndex, row := range rows {
		result.Rows[rowIndex] = make([]Column, len(columns))
		// pt1 := image.Pt(int(row[0]), int(row[1]))
		// pt2 := image.Pt(int(row[2]), int(row[3]))
		// gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)

		//todo - detect row treshold remove line if cv detected to lines one on top of each other .for instance where distancebetween lines is 1 pixel
		for columnIndex, column := range columns {
			// pt1 := image.Pt(int(column[0]), int(column[1]))
			// pt2 := image.Pt(int(column[2]), int(column[3]))
			// gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)

			x, y, err := t.intersection(columns[columnIndex], rows[rowIndex])
			if err != nil {
				log.Println("Cannot find intersection", column, row)
				continue
			}

			point := []int{x, y}

			result.Rows[rowIndex][columnIndex].Points = append(result.Rows[rowIndex][columnIndex].Points, point)

			//set intersection as end of previous column
			if columnIndex-1 >= 0 {
				result.Rows[rowIndex][columnIndex-1].Points = append(result.Rows[rowIndex][columnIndex-1].Points, point)
			}

			//set column intersection to previous row
			if rowIndex-1 >= 0 {
				result.Rows[rowIndex-1][columnIndex].Points = append(result.Rows[rowIndex-1][columnIndex].Points, point)
			}

			if rowIndex-1 >= 0 && columnIndex-1 >= 0 {
				result.Rows[rowIndex-1][columnIndex-1].Points = append(result.Rows[rowIndex-1][columnIndex-1].Points, point)
			}
		}
	}

	// for _, row := range result.Rows {
	// 	// fmt.Println("--------------row--------------")
	// 	for _, columns := range row {
	// 		fmt.Println("--------------column--------------")
	// 		randomColor := colorful.Hcl(180.0+rand.Float64()*50.0, 0.2+rand.Float64()*0.8, 0.3+rand.Float64()*0.7)

	// 		for _, point := range columns.Points {

	// 			fmt.Println(point)
	// 			pt1 := image.Pt(int(point[0]), int(point[1]))
	// 			gocv.Line(&mat, pt1, pt1, color.RGBA{uint8(randomColor.R), uint8(randomColor.G), uint8(randomColor.B), 50}, 4)
	// 			gocv.PutText(&mat, fmt.Sprint(pt1), pt1, gocv.FontHersheySimplex, 0.4, color.RGBA{uint8(randomColor.R), uint8(randomColor.G), uint8(randomColor.B), 50}, 1)
	// 		}
	// 	}
	// }

	return result, mat, nil
}

//copied from https://stackoverflow.com/questions/20677795/how-do-i-compute-the-intersection-point-of-two-lines
func (t *TableStructureDetector) intersection(line1, line2 []int) (int, int, error) {
	xdiff := []int{line1[0] - line1[2], line2[0] - line2[2]}
	ydiff := []int{line1[1] - line1[3], line2[1] - line2[3]}

	det := func(a, b []int) int {
		return a[0]*b[1] - a[1]*b[0]
	}

	div := det(xdiff, ydiff)
	if div == 0 {
		return 0, 0, errors.New("Cannnot find interseciton")
	}

	d := []int{det([]int{line1[0], line1[1]}, []int{line1[2], line1[3]}), det([]int{line2[0], line2[1]}, []int{line2[2], line2[3]})}
	x := det(d, xdiff) / div
	y := det(d, ydiff) / div

	return x, y, nil
}

func appendSortedBy(lines [][]int, line []int, valueIndex int) [][]int {
	if len(lines) == 0 {
		lines = append(lines, line)
		return lines
	}

	sortValues := make([]int, 0)
	for _, point := range lines {
		sortValues = append(sortValues, point[valueIndex])
	}

	i := sort.SearchInts(sortValues, line[valueIndex])
	lines = append(lines, []int{})
	copy(lines[i+1:], lines[i:])
	lines[i] = line
	return lines
}

//TableStructure ...
type TableStructure struct {
	Rows [][]Column
}

//Column ...
type Column struct {
	Points [][]int
}
