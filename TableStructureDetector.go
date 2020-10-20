package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

//TableStructureDetector ...
type TableStructureDetector struct {
}

//Detect ...
func (t *TableStructureDetector) Detect(filename string) (TableStructure, error) {

	mat := gocv.IMRead(filename, gocv.IMReadColor)

	matCanny := gocv.NewMat()
	matGray := gocv.NewMat()
	matLines := gocv.NewMat()

	gocv.CvtColor(mat, &matGray, gocv.ColorBGRToGray)

	gocv.Canny(matGray, &matCanny, 100, 150)
	gocv.HoughLinesPWithParams(matCanny, &matLines, 1, math.Pi/180, 200, 2, 4)

	result := TableStructure{Rows: make([]Column, 0)}

	fmt.Println(matLines.Cols())
	fmt.Println(matLines.Rows())
	for i := 0; i < matLines.Rows(); i++ {
		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)

		result.Rows = append(result.Rows, Column{})
	}

	return result, nil
}

//TableStructure ...
type TableStructure struct {
	Rows []Column
}

//Column ...
type Column struct {
}
