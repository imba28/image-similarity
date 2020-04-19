package main

import (
	"bufio"
	"flag"
	"fmt"
	"imba28/images/pkg"
	"os"
	"strings"
)

const (
	distanceThreshold  = 10
	maxResultSetLength = 10
)

func main() {
	dir := flag.String("directory", "images", "Directory that contains the images set")
	flag.Parse()

	fmt.Println("Building index...")
	index, err := pkg.NewIndex(*dir)
	if err != nil {
		fmt.Printf("could not open image directory %q\n", *dir)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Input image name: ")
		text, _ := reader.ReadString('\n')

		distances, err := index.Search(text, distanceThreshold, maxResultSetLength)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, distance := range distances {
			fmt.Println(distance.Path, distance.Distance)
			pkg.DisplayImage(distance.Path)
		}
	}
}
