package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
)

func ShowMask(mat gocv.Mat) {
	if true {
		return
	}

	window := gocv.NewWindow("Mask")
	defer window.Close()

	for {
		window.IMShow(mat)
		if window.WaitKey(0) >= 0 {
			break
		}
	}
}

func DisplayImage(path string) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("error reading image from %q\n", path)
		return
	}

	window := gocv.NewWindow("Result")
	defer window.Close()

	for {
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
