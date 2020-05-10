package file

import (
	"imba28/images/pkg"
	"strconv"
)

// sets guid&ids of images in ascending order
// when using a simple file provider we have to set an id to make images individually identifiable
type ImageGuidProvider struct {
	pkg.ImageProvider
	images *[]*pkg.Image
}

func (p *ImageGuidProvider) Images() ([]*pkg.Image, error) {
	if p.images != nil {
		return *p.images, nil
	}

	images, err := p.ImageProvider.Images()
	if err != nil {
		return nil, err
	}

	for i := range images {
		images[i].Guid = i
		images[i].Id = strconv.Itoa(i)
	}

	p.images = &images

	return images, nil
}

func (p *ImageGuidProvider) Get(guid string) *pkg.Image {
	if p.images == nil {
		_, err := p.Images()
		if err != nil {
			return nil
		}
	}

	guidInt, err := strconv.Atoi(guid)
	if err != nil {
		return nil
	}

	for i := range *p.images {
		if guidInt == (*p.images)[i].Guid {
			return (*p.images)[i]
		}
	}

	return nil
}

func NewImageGuidProvider(main pkg.ImageProvider) pkg.ImageProvider {
	return &ImageGuidProvider{ImageProvider: main}
}
