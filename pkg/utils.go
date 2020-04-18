package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
)

func DirectoryIndex(dir string) ([]ImageDescriptor, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var images []ImageDescriptor
	for _, file := range files {
		imagePath := dir + "/" + file.Name()
		features, err := FeatureVector(imagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		images = append(images, ImageDescriptor{
			Path:     imagePath,
			Features: features,
		})
	}

	return images, nil
}

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
