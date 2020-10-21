package main

import (
	"fmt"
	"io/ioutil"

	gosseract "github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

func main() {
	//load document
	//save pages as images?
	//open cv- detect edges
	//ocr - get bonding boxes
	// map edges to row/columns and letters

	imageFilename := "8207_datagrid.jpg"
	structDetector := TableStructureDetector{}
	tableStructure, matrix, _ := structDetector.Detect(imageFilename)

	// fmt.Println(strucutre)
	window := gocv.NewWindow("detected lines")
	for {
		window.ResizeWindow(15000, 20000)
		window.IMShow(matrix)
		if window.WaitKey(10) >= 0 {
			break
		}
	}

	//convert matrix to image so ocr can read it
	buff, err := gocv.IMEncode(gocv.JPEGFileExt, matrix)
	fmt.Println("encode err", err)

	ioutil.WriteFile("converted.jpg", buff, 0644)

	boundingBoxes := ocrFromBytes(buff)
	connectTableColumnsAndOcrCharacters(boundingBoxes, tableStructure)
}

func ocrFromBytes(imageBytes []byte) []gosseract.BoundingBox {
	client := gosseract.NewClient()
	defer client.Close()
	err := client.SetImageFromBytes(imageBytes)
	fmt.Println("-> err:", err)

	bb, err := client.GetBoundingBoxes(gosseract.RIL_SYMBOL)
	fmt.Println("-> err:", err)

	// for _, bb := range bb {
	// 	fmt.Println(bb.Word)
	// }
	return bb
}

func doBoxesIntersect(topLeft, bottomRight []int, point []int) bool {
	if topLeft[0] <= point[0] && point[0] <= bottomRight[0] && topLeft[1] <= point[0] && point[1] <= bottomRight[1] {
		return true
	}
	return false
}

func connectTableColumnsAndOcrCharacters(ocrBoundingBoxes []gosseract.BoundingBox, tableStructure TableStructure) {

	content := make([][]string, 0)
	for _, row := range tableStructure.Rows {
		columns := make([]string, 0)
		fmt.Println("")
		for _, column := range row {
			if len(column.Points) < 4 {
				// fmt.Println(column.Points)
				continue
			}
			topLeft := column.Points[0]
			bottomRight := column.Points[3]

			columnContent := ""
			for _, box := range ocrBoundingBoxes {
				if doBoxesIntersect(topLeft, bottomRight, []int{box.Box.Min.X, box.Box.Min.Y}) || doBoxesIntersect(topLeft, bottomRight, []int{box.Box.Max.X, box.Box.Max.Y}) {
					columnContent += box.Word
				}
			}
			columns = append(columns, columnContent)
		}

		if len(columns) > 0 {
			content = append(content, columns)
		}
	}

	for _, row := range content {
		columns := ""
		for _, column := range row {
			columns += column + " | "
		}
		fmt.Println(columns)
	}
}
