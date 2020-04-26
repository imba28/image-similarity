package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
)

type Image struct {
	Id       string
	Guid     int
	Path     string
	Name     string
	Features []float64
}

type ImageProvider interface {
	Images() ([]*Image, error)
	Get(string) *Image
	Persist(*Image) error
}

type ImageIndex struct {
	provider ImageProvider
	images   []*Image
	imageMap map[string]*Image
}

func (i *ImageIndex) Add(image *Image) error {
	if image == nil || i.Has(image.Id) {
		return nil
	}

	if image.Features == nil {
		feature, err := FeatureVector(*image)
		if err != nil {
			return err
		}

		image.Features = feature

		err = i.provider.Persist(image)
		if err != nil {
			return err
		}
	}

	i.imageMap[image.Id] = image
	i.images = append(i.images, image)

	return nil
}

func (i ImageIndex) Has(id string) bool {
	_, ok := i.imageMap[id]
	return ok
}

func (i ImageIndex) Load(id string) *Image {
	return i.provider.Get(id)
}

func (i ImageIndex) Images() []*Image {
	return i.images
}

func (i ImageIndex) Search(referenceImage Image, distanceThreshold float64, limit int) ([]ImageDistance, error) {
	distances, err := CalculateDistances(referenceImage, i.images)
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

	var imageDescriptors []*Image
	imageMap := make(map[string]*Image, len(images))

	imageIndex := &ImageIndex{
		provider: p,
		images:   imageDescriptors,
		imageMap: imageMap,
	}

	for i := range images {
		err := imageIndex.Add(images[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return imageIndex, nil
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
