package main

import (
	"bufio"
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

type ImageFeature struct {
	Path string
	Features []float64
}

type ImageDistance struct {
	Path string
	Distance float64
}

func displayImage(path string) {
	window := gocv.NewWindow("Result")
	defer window.Close()
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("error reading image from %q\n", path)
		return
	}
	for {
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func featureVector(path string) ([]float64, error) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("error reading image from %q\n", path)
	}

	hsvImg := gocv.NewMat()
	defer hsvImg.Close()
	img.ConvertTo(&hsvImg, gocv.ColorBGRToHSV)

	/*width, height := hsvImg.Size()[1], hsvImg.Size()[0]
	cx, cy := int(width / 2), int(height / 2)

	segments := [][]int {
		{0, 0, cx, cy}, // top left
		{cx, 0, width, cy}, // top right
		{0, cy, cx, height}, // bottom left
		{cx, cy, width, height}, // bottom right
	}*/

	mask := gocv.NewMat()
	defer mask.Close()

	hist := gocv.NewMat()
	defer hist.Close()

	bins := []int{8, 12, 3}
	// h [0,180], s[0,256], v[0,256]
	ranges := []float64{0, 180, 0, 256, 0, 256}
	gocv.CalcHist([]gocv.Mat{hsvImg}, []int{0,1,2}, mask, &hist, bins, ranges, false)
	gocv.Normalize(hist, &hist, 1, 0, gocv.NormL2)

	hist2 := gocv.NewMat()
	hist.ConvertTo(&hist2, gocv.MatTypeCV64F)

	features,err := hist2.DataPtrFloat64()
	if err != nil {
		return nil, err
	}

	return features, nil

	// axesX, axesY := int((width * 0.75) / 2), int((height * 0.75) / 2)
}

func directoryIndex(dir string) []ImageFeature {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("could not open source directory")
		return nil
	}

	var images []ImageFeature
	for _, file := range files {
		imagePath := dir + "/" + file.Name()
		features, err := featureVector(imagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		images = append(images, ImageFeature{
			Path:     imagePath,
			Features: features,
		})
	}

	return images
}


func calculateDistances(referencePath string, vectors []ImageFeature) ([]ImageDistance, error) {
	referenceVector, err := featureVector(referencePath)
	if err != nil {
		return nil, err
	}

	var d []ImageDistance

	for _, vector := range vectors {
		d = append(d, ImageDistance{
			Path:     vector.Path,
			Distance: chi2Distance(referenceVector, vector.Features),
		})
	}

	return d, nil
}

func chi2Distance(v1, v2 []float64) float64 {
	d := 0.
	min := int(math.Min(float64(len(v1)), float64(len(v2))))

	for i := 0; i < min; i++ {
		distance := math.Pow(v1[i] - v2[i], 2) / (v1[i] + v2[i] + 1e-10)
		d += distance
	}

	return d * 0.5
}

type ByDistance []ImageDistance
func (a ByDistance) Len() int {
	return len(a)
}
func (a ByDistance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByDistance) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}

func main() {
	fmt.Println("Building index...")
	directoryPath := "images"
	index := directoryIndex(directoryPath)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Input image name: ")
		text, _ := reader.ReadString('\n')
		distances, err := calculateDistances("images/" + strings.Trim(text, "\n"), index)
		if err != nil {
			fmt.Println(err)
			return
		}

		sort.Sort(ByDistance(distances))

		for i, distance := range distances {
			if distance.Distance > 1 || i > 10 {
				break
			}
			fmt.Println(distance.Path, distance.Distance)
			displayImage(distance.Path)
		}
	}
}
