package main

import (
	"bufio"
	"fmt"
	"imba28/images/pkg"
	"os"
	"strings"
)

const (
	indexDirectory     string = "images"
	distanceThreshold         = 10
	maxResultSetLength        = 10
)

func main() {
	fmt.Println("Building index...")
	index := pkg.DirectoryIndex(indexDirectory)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Input image name: ")
		text, _ := reader.ReadString('\n')
		distances, err := pkg.CalculateDistances("images/"+strings.Trim(text, "\n"), index)
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
