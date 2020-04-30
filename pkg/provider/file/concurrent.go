package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"os"
	"sync"
)

type ConcurrentImageProvider struct {
	Dir string
}

func listDir(c chan<- pkg.Image, wg *sync.WaitGroup, info os.FileInfo, path string) {
	if info.IsDir() {
		fs, _ := ioutil.ReadDir(path + "/" + info.Name())
		for _, f := range fs {
			if f.Name()[0] != '.' {
				wg.Add(1)
				go func(f os.FileInfo) {
					listDir(c, wg, f, path+"/"+info.Name())
				}(f)
			}
		}
		wg.Done()
		return
	}

	image := pkg.Image{
		Name: info.Name(),
		Path: path + "/" + info.Name(),
	}

	f, err := pkg.FeatureVector(image)
	if err == nil {
		image.Features = f
		c <- image
	}
	wg.Done()
}

func (c ConcurrentImageProvider) Images() ([]*pkg.Image, error) {
	fs, err := ioutil.ReadDir(c.Dir)
	if err != nil {
		return nil, err
	}

	imageChannel := make(chan pkg.Image, 8)
	allDone := make(chan bool)

	var wg sync.WaitGroup

	// launch a go routine for each file
	for _, f := range fs {
		wg.Add(1)
		go func(f os.FileInfo) {
			listDir(imageChannel, &wg, f, c.Dir)
		}(f)
	}

	var list []*pkg.Image

	go func() {
		// when all routines finished notify the main routine
		wg.Wait()
		allDone <- true
	}()

	for {
		select {
		case image := <-imageChannel:
			// a routine successfully calculated a feature vector. We should add it to the list
			list = append(list, &image)

		case <-allDone:
			// all routines are done and the list is complete
			return list, nil
		}
	}
}

func (c ConcurrentImageProvider) Get(path string) *pkg.Image {
	_, err := os.Stat(c.Dir + "/" + path)
	if err != nil {
		return nil
	}

	return &pkg.Image{
		Id:   path,
		Path: c.Dir + "/" + path,
		Name: path,
	}
}

func (c ConcurrentImageProvider) Persist(*pkg.Image) error {
	// file provider keep images in memory only
	// noop
	return nil
}

var _ pkg.ImageProvider = (*ConcurrentImageProvider)(nil)
