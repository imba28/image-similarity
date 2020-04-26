package pkg

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

const (
	hueBins        = 12
	saturationBins = 12
	brightnessBins = 4
)

func FeatureVector(i Image) ([]float64, error) {
	img := gocv.IMRead(i.Path, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("error reading image from %q\n", i.Path)
	}

	// convert img to hsv color space
	img.ConvertTo(&img, gocv.ColorBGRToHSV)

	black := color.RGBA{0, 0, 0, 0}
	white := color.RGBA{255, 255, 255, 0}

	width, height := img.Size()[1], img.Size()[0]
	cx, cy := int(width/2), int(height/2)

	segments := [][]int{
		{0, 0, cx, cy},          // top left
		{cx, 0, width, cy},      // top right
		{0, cy, cx, height},     // bottom left
		{cx, cy, width, height}, // bottom right
	}

	axesX, axesY := int((float32(width)*0.75)/2), int((float32(height)*0.75)/2)
	ellipMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer ellipMask.Close()
	gocv.Ellipse(&ellipMask, image.Point{cx, cy}, image.Point{axesX, axesY}, 0, 0, 360, white, -1)

	segmentMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer segmentMask.Close()

	var features []float64

	segmentHistogram := gocv.NewMat()
	defer segmentHistogram.Close()
	for _, segment := range segments {
		// reset mask
		gocv.Rectangle(&segmentMask, image.Rect(0, 0, width, height), black, -1)

		// calculate intersection of current segment and elliptic mask
		gocv.Rectangle(&segmentMask, image.Rect(segment[0], segment[1], segment[2], segment[3]), white, -1)
		gocv.Subtract(segmentMask, ellipMask, &segmentMask)

		ShowMask(segmentMask)

		segmentFeatures, err := featuresInSegment(img, segmentMask, segmentHistogram)
		if err != nil {
			return nil, err
		}

		features = append(features, segmentFeatures...)
	}

	ShowMask(ellipMask)

	ellipFeatures, err := featuresInSegment(img, ellipMask, segmentHistogram)
	if err != nil {
		return nil, err
	}
	features = append(features, ellipFeatures...)

	return features, nil
}

func featuresInSegment(img gocv.Mat, mask gocv.Mat, hist gocv.Mat) ([]float64, error) {
	// h, s, v channel
	channels := []int{0, 1, 2}
	bins := []int{hueBins, saturationBins, brightnessBins}
	// hue range = 0-180, saturation range = 0-256, value/brightness range = 0-256
	ranges := []float64{0, 180, 0, 256, 0, 256}

	gocv.CalcHist([]gocv.Mat{img}, channels, mask, &hist, bins, ranges, false)
	gocv.Normalize(hist, &hist, 1, 0, gocv.NormL2)

	float64Hist := gocv.NewMat()
	defer float64Hist.Close()
	hist.ConvertTo(&float64Hist, gocv.MatTypeCV64F)

	return float64Hist.DataPtrFloat64()
}
