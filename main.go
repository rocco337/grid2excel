package main

import (
	"fmt"

	gosseract "github.com/otiai10/gosseract"
)

func main() {
	//load document
	//save pages as images?
	//open cv- detect edges
	//ocr - get bonding boxes
	// map edges to row/columns and letters

	structDetector := TableStructureDetector{}
	structDetector.Detect("test1.jpg")
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

// func detectEdges(filename string) {
// 	mat := gocv.IMRead(filename, gocv.IMReadColor)

// 	matCanny := gocv.NewMat()
// 	matGray := gocv.NewMat()

// 	matLines := gocv.NewMat()

// 	window := gocv.NewWindow("detected lines")

// 	gocv.CvtColor(mat, &matGray, gocv.ColorBGRToGray)

// 	gocv.Canny(matGray, &matCanny, 100, 150)
// 	gocv.HoughLinesPWithParams(matCanny, &matLines, 1, math.Pi/180, 200, 2, 4)

// 	rows := make([]gocv.Veci, 0)
// 	columns := make([]gocv.Veci, 0)
// 	for i := 0; i < matLines.Rows(); i++ {
// 		line := matLines.GetVeciAt(i, 0)

// 		if line[1] > line[3] {
// 			fmt.Println("vertical")
// 			columns = append(columns, line)
// 		} else if line[0] < line[2] {
// 			fmt.Println("horizontal")
// 			rows = append(rows, line)
// 		}
// 		// pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
// 		// pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
// 		// gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 2)
// 		// gocv.PutText(&mat, fmt.Sprint(pt1), pt1, gocv.FontHersheySimplex, 0.4, color.RGBA{0, 255, 0, 50}, 2)
// 		// gocv.PutText(&mat, fmt.Sprint(pt2), pt2, gocv.FontHersheySimplex, 0.4, color.RGBA{255, 0, 0, 50}, 2)
// 		// if i > 10 {
// 		// 	break
// 		// }
// 	}

// 	intersections := make([]gocv.Veci, 0)
// 	for _, column := range columns {
// 		for _, row := range rows {

// 			x, y, err := intersection(column, row)
// 			if err != nil {
// 				log.Println("Cannot find intersection", column, row)
// 			} else {
// 				intersections = append(intersections, gocv.Veci{x, y})
// 			}

// 		}
// 	}

// 	fmt.Println(intersections)
// 	for _, intersection := range intersections {
// 		pt1 := image.Pt(int(intersection[0]), int(intersection[1]))

// 		gocv.Line(&mat, pt1, pt1, color.RGBA{0, 255, 0, 50}, 4)
// 		// gocv.PutText(&mat, fmt.Sprint(pt1), pt1, gocv.FontHersheySimplex, 0.4, color.RGBA{0, 255, 0, 50}, 2)
// 	}

// 	// fmt.Println(rows)
// 	// fmt.Println(columns)

// 	for {
// 		window.ResizeWindow(15000, 20000)
// 		window.IMShow(mat)

// 		if window.WaitKey(10) >= 0 {
// 			time.Sleep(10 * time.Minute)
// 		}
// 	}
// }

// func intersection(line1, line2 gocv.Veci) (int32, int32, error) {
// 	xdiff := []int32{line1[0] - line1[2], line2[0] - line2[2]}
// 	ydiff := []int32{line1[1] - line1[3], line2[1] - line2[3]}

// 	// det := func(line gocv.Veci) int32 {
// 	// 	return line[0]*line[2] - line[1]*line[3]
// 	// }

// 	det := func(a, b []int32) int32 {
// 		return a[0]*b[1] - a[1]*b[0]
// 	}

// 	div := det(xdiff, ydiff)
// 	if div == 0 {
// 		return 0, 0, errors.New("Cannnot find interseciton")
// 	}

// 	d := []int32{det([]int32{line1[0], line1[1]}, []int32{line1[2], line1[3]}), det([]int32{line2[0], line2[1]}, []int32{line2[2], line2[3]})}
// 	x := det(d, xdiff) / div
// 	y := det(d, ydiff) / div
// 	return x, y, nil
// }

// def line_intersection(line1, line2):
//     xdiff = (line1[0][0] - line1[1][0], line2[0][0] - line2[1][0])
//     ydiff = (line1[0][1] - line1[1][1], line2[0][1] - line2[1][1])

//     def det(a, b):
//         return a[0] * b[1] - a[1] * b[0]

//     div = det(xdiff, ydiff)
//     if div == 0:
//        raise Exception('lines do not intersect')

//     d = (det(*line1), det(*line2))
//     x = det(d, xdiff) / div
//     y = det(d, ydiff) / div
// 	return x, y

// bool intersection(Point2f o1, Point2f p1, Point2f o2, Point2f p2,
// 	Point2f &r)
// {
// Point2f x = o2 - o1;
// Point2f d1 = p1 - o1;
// Point2f d2 = p2 - o2;

// float cross = d1.x*d2.y - d1.y*d2.x;
// if (abs(cross) < /*EPS*/1e-8)
// return false;

// double t1 = (x.x * d2.y - x.y * d2.x)/cross;
// r = o1 + d1 * t1;
// return true;
// }
