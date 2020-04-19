package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"strings"
)

type ImageIndex struct {
	dir         string
	descriptors []ImageDescriptor
}

func (i ImageIndex) Descriptors() []ImageDescriptor {
	return i.descriptors
}

func (i ImageIndex) Search(referenceImagePath string, distanceThreshold float64, limit int) ([]ImageDistance, error) {
	distances, err := CalculateDistances(i.dir+"/"+strings.Trim(referenceImagePath, "\n"), i.descriptors)
	if err != nil {
		return nil, err
	}

	var results []ImageDistance

	for i, distance := range distances {
		if distance.Distance > distanceThreshold || i > limit {
			break
		}
		results = append(results, distance)
	}

	return results, nil
}

func NewIndex(dir string) (*ImageIndex, error) {
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
			Name:     file.Name(),
			Path:     imagePath,
			Features: features,
		})
	}

	return &ImageIndex{
		dir:         dir,
		descriptors: images,
	}, nil
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
