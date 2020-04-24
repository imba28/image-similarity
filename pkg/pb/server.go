package pb

import (
	"context"
	"errors"
	"fmt"
	"imba28/images/pkg"
	"strconv"
)

type ImageSimilarityService struct {
	index *pkg.ImageIndex
}

func (s ImageSimilarityService) GetSimilar(c context.Context, r *ImageRequest) (*ImageSimilarityResponse, error) {
	image := s.index.Load(strconv.Itoa(int(r.Image.Id)))

	fmt.Printf("finding similar images to %d\n", r.Image.Id)

	if image == nil {
		return nil, errors.New("image not found")
	}

	imageDistances, err := s.index.Search(*image, 10, int(r.Limit))
	if err != nil {
		return nil, err
	}

	var images []*ImageSimilarity
	for i := range imageDistances {
		id, err := strconv.Atoi(imageDistances[i].Image.Id)
		if err != nil {
			fmt.Printf("Could not convert id %q to int!", imageDistances[i].Image.Id)
		}
		images = append(images, &ImageSimilarity{
			Image:    &Image{Id: int32(id)},
			Distance: imageDistances[i].Distance,
		})
	}

	return &ImageSimilarityResponse{
		Similarities: images,
		Count:        int32(len(images)),
	}, nil
}

func (s ImageSimilarityService) AddImage(c context.Context, i *Image) (*Image, error) {
	if i == nil {
		return nil, errors.New("image not found")
	}
	id := strconv.Itoa(int(i.Id))
	fmt.Printf("adding image %d to the index.\n", i.Id)

	if s.index.Has(id) {
		return i, nil
	}

	image := s.index.Load(id)
	if image == nil {
		return nil, errors.New("image not found")
	}

	err := s.index.Add(image)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func NewImageSimilarityService(index *pkg.ImageIndex) ImageSimilarityService {
	return ImageSimilarityService{
		index: index,
	}
}
