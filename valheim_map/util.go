package valheim_map

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var Verbosity = 0

func Debug(text string) {
	if Verbosity > 0 {
		fmt.Printf("DEBUG: %s\n", text)
	}
}

// Clamp restricts a value to a specified range [min, max]
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func Normalize(value float32, min float32, max float32) float32 {
	if min >= max {
		panic(fmt.Sprintf("For Normalize, min (%f) must be less than max (%f)!", min, max))
	}
	return (Clamp(value, min, max) - min) / (max - min)
}

func FastCopyRGBA(dst, src *image.RGBA, xOffset, yOffset int) {
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()

	// Ensure we don't copy out of bounds
	width := min(srcBounds.Dx(), dstBounds.Dx()-xOffset)
	height := min(srcBounds.Dy(), dstBounds.Dy()-yOffset)
	for y := 0; y < height; y++ {
		// Compute the index of the start of the row in both images
		srcStart := y * src.Stride
		dstStart := (yOffset+y)*dst.Stride + xOffset*4
		// Copy one full row at a time (avoids per-pixel overhead)
		copy(dst.Pix[dstStart:dstStart+width*4], src.Pix[srcStart:srcStart+width*4])
	}
}

var fileWriteSem = make(chan struct{}, 4) // limit to 4 conncurrent writes
func SavePng(pngPath string, imageData *image.RGBA, showBar bool) {
	// Write PNG file
	fileWriteSem <- struct{}{}        // Acquire write slot
	defer func() { <-fileWriteSem }() // Defer release
	pngFile, err := os.Create(pngPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create png file at '%s' (%s)", pngPath, err.Error()))
	}
	defer pngFile.Close()
	bar := progressbar.NewOptions(-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(fmt.Sprintf("Saving png '%s'", pngPath)),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetVisibility(showBar),
	)
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	err = encoder.Encode(io.MultiWriter(pngFile, bar), imageData)
	if err != nil {
		panic(fmt.Sprintf("Could not encode png file at '%s' (%s)", pngPath, err.Error()))
	}
	pngFile.Close()
	bar.Close()
}
