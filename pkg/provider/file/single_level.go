package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"os"
	"strings"
)

type SingleLevelImageProvider struct {
	dir string
}

func (f SingleLevelImageProvider) Images() ([]*pkg.Image, error) {
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

func (f SingleLevelImageProvider) Get(id string) *pkg.Image {
	i := newImage(f.dir + "/" + strings.Trim(id, "\n"))
	return &i
}

func (f SingleLevelImageProvider) Persist(i *pkg.Image) error {
	// noop
	return nil
}

func newImage(path string) pkg.Image {
	parts := strings.Split(path, "/")
	return pkg.Image{
		Id:   parts[len(parts)-1],
		Path: path,
		Name: parts[len(parts)-1],
	}
}

func NewSingleLevelProvider(dir string) SingleLevelImageProvider {
	return SingleLevelImageProvider{dir: dir}
}

var _ pkg.ImageProvider = (*SingleLevelImageProvider)(nil)
