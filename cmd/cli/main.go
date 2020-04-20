package main

import (
	"bufio"
	"flag"
	"fmt"
	"imba28/images/pkg"
	"imba28/images/pkg/provider/file"
	"os"
)

const (
	distanceThreshold  = 10
	maxResultSetLength = 10
)

func main() {
	dir := flag.String("directory", "images", "Directory that contains the images set")
	flag.Parse()

	fmt.Println("Building index...")
	index, err := pkg.NewIndex(file.New(*dir))
	if err != nil {
		fmt.Printf("could not open image directory %q\n", *dir)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Input image name: ")
		text, _ := reader.ReadString('\n')

		image := index.Get(text)
		if image == nil {
			fmt.Println("image not found")
			return
		}
		distances, err := index.Search(*image, distanceThreshold, maxResultSetLength)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, distance := range distances {
			fmt.Println(distance.Image.Path, distance.Distance)
			pkg.DisplayImage(distance.Image.Path)
		}
	}
}
