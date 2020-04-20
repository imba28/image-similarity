package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"strings"
)

type ImageProvider struct {
	dir string
}

func (f ImageProvider) Images() ([]pkg.Image, error) {
	files, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return nil, err
	}

	var images []pkg.Image

	for _, file := range files {
		images = append(images, pkg.Image{
			Path: f.dir + "/" + file.Name(),
			Name: file.Name(),
		})
	}

	return images, nil
}

func NewImage(path string) pkg.Image {
	parts := strings.Split(path, "/")
	return pkg.Image{
		Path: path,
		Name: parts[len(parts)-1],
	}
}

func New(dir string) ImageProvider {
	return ImageProvider{dir: dir}
}
