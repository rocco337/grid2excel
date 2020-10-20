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
	//detectEdges("test1.jpg")
	ocr("test2.jpg")
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

	fmt.Println(matLines.Cols())
	fmt.Println(matLines.Rows())
	for i := 0; i < matLines.Rows(); i++ {
		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)
	}

	for {
		window.ResizeWindow(1500, 2000)
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
