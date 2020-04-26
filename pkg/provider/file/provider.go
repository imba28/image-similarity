package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"os"
	"strings"
)

type ImageProvider struct {
	dir string
}

func (f ImageProvider) Images() ([]*pkg.Image, error) {
	files, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return nil, err
	}

	var images []*pkg.Image

	for _, file := range files {
		if file.IsDir() || file.Mode()&os.ModeSymlink != 0 || (len(file.Name()) > 0 && file.Name()[0] == '.') {
			continue
		}

		images = append(images, &pkg.Image{
			Id:   file.Name(),
			Path: f.dir + "/" + file.Name(),
			Name: file.Name(),
		})
	}

	return images, nil
}

func (f ImageProvider) Get(id string) *pkg.Image {
	i := NewImage(f.dir + "/" + strings.Trim(id, "\n"))
	return &i
}

func (f ImageProvider) Persist(i *pkg.Image) error {
	// noop
	return nil
}

func NewImage(path string) pkg.Image {
	parts := strings.Split(path, "/")
	return pkg.Image{
		Id:   parts[len(parts)-1],
		Path: path,
		Name: parts[len(parts)-1],
	}
}

func New(dir string) ImageProvider {
	return ImageProvider{dir: dir}
}

var _ pkg.ImageProvider = (*ImageProvider)(nil)
