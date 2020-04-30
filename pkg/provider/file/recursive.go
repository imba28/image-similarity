package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"os"
)

type RecursiveImageProvider struct {
	Dir string
}

func rec(l *[]*pkg.Image, f os.FileInfo, dir string) {
	if f.IsDir() {
		fs, _ := ioutil.ReadDir(dir + "/" + f.Name())
		for _, ff := range fs {
			rec(l, ff, dir+"/"+f.Name())
		}
	} else {
		image := pkg.Image{
			Name: f.Name(),
			Path: dir + "/" + f.Name(),
		}

		f, err := pkg.FeatureVector(image)
		if err == nil {
			image.Features = f
			*l = append(*l, &image)
		}

	}
}

func (r RecursiveImageProvider) Images() ([]*pkg.Image, error) {
	var l []*pkg.Image

	fs, _ := ioutil.ReadDir(r.Dir)
	for _, f := range fs {
		rec(&l, f, r.Dir)
	}

	return l, nil
}

func (r RecursiveImageProvider) Get(path string) *pkg.Image {
	_, err := os.Stat(r.Dir + "/" + path)
	if err != nil {
		return nil
	}

	return &pkg.Image{
		Id:   path,
		Path: r.Dir + "/" + path,
		Name: path,
	}
}

func (r RecursiveImageProvider) Persist(*pkg.Image) error {
	// noop
	return nil
}

var _ pkg.ImageProvider = (*RecursiveImageProvider)(nil)
