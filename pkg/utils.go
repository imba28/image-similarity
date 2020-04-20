package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
)

type Image struct {
	Id   string
	Path string
	Name string
}

type ImageProvider interface {
	Images() ([]Image, error)
	Get(string) *Image
}

type ImageIndex struct {
	provider    ImageProvider
	descriptors []ImageDescriptor
}

func (i ImageIndex) Get(id string) *Image {
	return i.provider.Get(id)
}

func (i ImageIndex) Descriptors() []ImageDescriptor {
	return i.descriptors
}

func (i ImageIndex) Search(referenceImage Image, distanceThreshold float64, limit int) ([]ImageDistance, error) {
	distances, err := CalculateDistances(referenceImage, i.descriptors)
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

func NewIndex(p ImageProvider) (*ImageIndex, error) {
	images, err := p.Images()
	if err != nil {
		return nil, err
	}

	var imageDescriptors []ImageDescriptor
	for _, image := range images {
		features, err := FeatureVector(image)
		if err != nil {
			fmt.Println(err)
			continue
		}
		imageDescriptors = append(imageDescriptors, ImageDescriptor{
			Image:    image,
			Features: features,
		})
	}

	return &ImageIndex{
		provider:    p,
		descriptors: imageDescriptors,
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
