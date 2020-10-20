package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	gosseract "github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

func main() {
	//load document
	//save pages as images?
	//open cv- detect edges
	//ocr - get bonding boxes
	// map edges to row/columns and letters
	detectEdges("test1.jpg")
	//ocr("test2.jpg")
}

func ocr(filename string) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(filename)

	bb, err := client.GetBoundingBoxes(gosseract.RIL_SYMBOL)
	fmt.Println("-> err:", err)

	fmt.Println(bb)
}

func detectEdges(filename string) {
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
			fmt.Println("vertical")
			columns = append(columns, line)

		} else if line[0] < line[2] {
			fmt.Println("horizontal")
			rows = append(columns, line)
		}

		fmt.Println(line)

		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)
		// gocv.PutText(&mat, fmt.Sprint(pt1), pt1, gocv.FontHersheySimplex, 0.4, color.RGBA{0, 255, 0, 50}, 2)
		// gocv.PutText(&mat, fmt.Sprint(pt2), pt2, gocv.FontHersheySimplex, 0.4, color.RGBA{255, 0, 0, 50}, 2)
		if i > 10 {
			break
		}
	}

	fmt.Println(rows)
	fmt.Println(columns)

	for {
		window.ResizeWindow(15000, 20000)
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
