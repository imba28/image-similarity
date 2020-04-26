package pb

import (
	"context"
	"errors"
	"imba28/images/pkg"
	"log"
	"strconv"
)

type ImageSimilarityService struct {
	index *pkg.ImageIndex
	dir   string
}

func (s ImageSimilarityService) GetSimilar(c context.Context, r *ImageRequest) (*ImageSimilarityResponse, error) {
	image := s.index.Load(strconv.Itoa(int(r.Image.Guid)))

	log.Printf("[GRPC] \"%s %d\"", "GetSimilar", r.Image.Guid)

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
			log.Printf("[GRPC] \"%s %d\" Could not convert id %q to int!", "GetSimilar", r.Image.Guid, imageDistances[i].Image.Id)
			return nil, err
		}
		images = append(images, &ImageSimilarity{
			Image:    &Image{Guid: int32(id), Path: imageDistances[i].Image.Path},
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
	id := strconv.Itoa(int(i.Guid))
	log.Printf("[GRPC] \"%s %s\"", "AddImage", id)

	if s.index.Has(id) {
		return i, nil
	}

	image := s.index.Load(id)
	if image != nil {
		return i, nil
	}

	image = &pkg.Image{
		Guid: int(i.Guid),
		Path: i.Path,
		Name: i.Name,
	}

	if len(s.dir) > 0 {
		sep := "/"
		if image.Path[0] == '/' || s.dir[len(s.dir)-1] == '/' {
			sep = ""
		}
		image.Path = s.dir + sep + image.Path
	}

	err := s.index.Add(image)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func NewImageSimilarityService(index *pkg.ImageIndex, dir string) ImageSimilarityService {
	return ImageSimilarityService{
		index: index,
		dir:   dir,
	}
}
