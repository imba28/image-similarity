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
	index, err := pkg.DirectoryIndex(*dir)
	if err != nil {
		fmt.Printf("could not open image directory %q\n", *dir)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Input image name: ")
		text, _ := reader.ReadString('\n')
		distances, err := pkg.CalculateDistances(*dir+"/"+strings.Trim(text, "\n"), index)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, distance := range distances {
			if distance.Distance > distanceThreshold || i > maxResultSetLength {
				break
			}

			fmt.Println(distance.Path, distance.Distance)
			pkg.DisplayImage(distance.Path)
		}
	}
}
