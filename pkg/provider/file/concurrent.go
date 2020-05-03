package file

import (
	"imba28/images/pkg"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

type ConcurrentImageProvider struct {
	Dir string
}

func workerPool(workerCount int, filePool <-chan *pkg.Image, resultPool chan<- *pkg.Image, wg *sync.WaitGroup, exit chan bool) {
	signals := make([]chan bool, workerCount)

	for i := 0; i < workerCount; i++ {
		signals[i] = make(chan bool)
		go func(stopSignal <-chan bool) {
			for {
				select {
				case image := <-filePool:
					f, err := pkg.FeatureVector(*image)
					if err == nil {
						image.Features = f
						resultPool <- image
					} else {
						resultPool <- nil
					}
				case <-stopSignal:
					return
				}
			}
		}(signals[i])
	}

	<-exit
	for i := range signals {
		signals[i] <- true
	}
}

func readDir(path string, imagePool chan *pkg.Image, wg *sync.WaitGroup) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range fs {
		if f.IsDir() {
			readDir(path+"/"+f.Name(), imagePool, wg)
			continue
		}

		wg.Add(1)
		image := &pkg.Image{
			Name: f.Name(),
			Path: path + "/" + f.Name(),
		}
		imagePool <- image
	}
}

func (c ConcurrentImageProvider) Images() ([]*pkg.Image, error) {
	cores := runtime.NumCPU()
	imagePool := make(chan *pkg.Image, 10)
	resultPool := make(chan *pkg.Image, 10)
	done := make(chan bool)
	doneWorker := make(chan bool)
	var wg sync.WaitGroup
	var l []*pkg.Image

	go workerPool(cores, imagePool, resultPool, &wg, doneWorker)
	go func() {
		for {
			select {
			case image := <-resultPool:
				if image != nil {
					l = append(l, image)
				}
				wg.Done()
			case <-done:
				return
			}
		}
	}()

	readDir(c.Dir, imagePool, &wg)
	wg.Wait()

	doneWorker <- true
	done <- true

	return l, nil
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
