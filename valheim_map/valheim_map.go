package valheim_map

/*
Valheim Tile To Image - v1.0.0
Mark Seelye - mseelye@yahoo.com

See: --help

TODO:
- Break up config to allow for plot types
- Add other config to be able to be read from config file
	- Paths
  - Flags
	- biome colors and such
- Add water color property for biome
- remove "GO" out of png paths
- update so locations can scan for ImportantContents
*/

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"golang.org/x/image/font/gofont/goregular"
)

func parseTile(mapDescriptor *MapDescriptor, gzReader io.Reader) *image.RGBA {
	var biomeData uint16
	var heightData float32
	var forestData float32
	imageData := image.NewRGBA(image.Rect(0, 0, mapDescriptor.tileRowCount, mapDescriptor.tileRowCount))
	for tx := 0; tx < mapDescriptor.tileRowCount; tx++ {
		for tz := 0; tz < mapDescriptor.tileRowCount; tz++ {
			binary.Read(gzReader, binary.LittleEndian, &biomeData)
			binary.Read(gzReader, binary.LittleEndian, &heightData)
			binary.Read(gzReader, binary.LittleEndian, &forestData)
			biome := Biome(biomeData)
			pixelX := tx
			pixelZ := (mapDescriptor.tileRowCount - tz - 1)
			normalizedHeight := Normalize(heightData, -2, 70)
			var pixelColor color.RGBA
			if heightData < 20 {
				pixelColor = Blue
			} else if heightData < 30 && biome == Swamp {
				pixelColor = color.RGBA{0x00, 0x32, 0x7f, 0xff} // ColorAverage(biome.Color(), Lightblue)
			} else if heightData < 30 {
				pixelColor = ColorAverage(biome.Color(), Blue)
			} else {
				pixelColor = ColorDarken(biome.Color(), normalizedHeight)
			}
			imageData.Set(pixelX, pixelZ, pixelColor)
		}
	}
	return imageData
}

// Mutex used for safely writing to big image file
var mu *sync.Mutex

// Note: Draw is very slow: draw.Draw(imageDataBig, image.Rect(destX, destZ, destX+tileRowCount, destZ+tileRowCount), imageData, image.Point{0, 0}, draw.Over)
// This uses a custom FastCopyRGBA to copy data directly (while locked with a mutex)
func writeTileImageToBigImage(mapDescriptor *MapDescriptor, x int, z int, imageDataBig *image.RGBA, imageData *image.RGBA) {
	destX := (x * mapDescriptor.tileRowCount)
	destZ := ((mapDescriptor.tileSideCount - z - 1) * mapDescriptor.tileRowCount)
	mu.Lock()
	FastCopyRGBA(imageDataBig, imageData, destX, destZ)
	mu.Unlock()
}

func processBitmap(mapDescriptor *MapDescriptor, x int, z int, imageDataBig *image.RGBA, forceReload bool) {
	pngData, err := os.Open(mapDescriptor.TileOutputPath(x, z))
	if err == nil && !forceReload {
		// png exists, skip processing for now
		// later load png and apply to bigmap
		imageData, err := png.Decode(pngData)
		if err != nil {
			panic(fmt.Sprintf("Could not decode existing png file at '%s' (%s)", mapDescriptor.TileOutputPath(x, z), err.Error()))
		}
		writeTileImageToBigImage(mapDescriptor, x, z, imageDataBig, imageData.(*image.RGBA))
	} else {
		// png does not exist yet, read gz, process, save to map pngs
		tileGz, err := os.Open(mapDescriptor.TileInputPath(x, z))
		if err != nil {
			panic(fmt.Sprintf("Could not open map tile file at '%s' (%s)", mapDescriptor.TileInputPath(x, z), err.Error()))
		}
		defer tileGz.Close()
		gzReader, err := gzip.NewReader(tileGz)
		if err != nil {
			panic(fmt.Sprintf("Could gz read gz file at '%s' (%s)", mapDescriptor.TileInputPath(x, z), err.Error()))
		}
		defer gzReader.Close()
		imageData := parseTile(mapDescriptor, gzReader)
		writeTileImageToBigImage(mapDescriptor, x, z, imageDataBig, imageData)
		SavePng(mapDescriptor.TileOutputPath(x, z), imageData, false)
	}
}

func parsePosition(positionString string) (float64, float64, float64) {
	posRE := regexp.MustCompile(`\(([^,]+), ([^,]+), ([^)]+)\)`)
	result := posRE.FindStringSubmatch(positionString)
	if len(result) > 0 {
		x, _ := strconv.ParseFloat(result[1], 64)
		y, _ := strconv.ParseFloat(result[2], 64)
		z, _ := strconv.ParseFloat(result[3], 64)
		return x, y, z
	} else {
		fmt.Printf("ERROR: Invalid Position Format: '%s'\n", positionString)
		return -1.0, -1.0, -1.0
	}
}

func plotLocations(mapDescriptor *MapDescriptor, imageDataBig *image.RGBA, locationsData *LocationsData, plot *Plot) []TextPlot {
	scaleFactor := mapDescriptor.ScaleFactor()
	imageWidth := mapDescriptor.ImageWidth()
	iconPngData, err := OpenIcon(mapDescriptor, plot.IconPng)
	if err == nil {
		iconPngData = TintImage(iconPngData, plot.GetIconColor(), 100)
	} else {
		fmt.Println("\n", err.Error())
	}
	iconPngDataSizedTinted := resize.Resize(uint(mapDescriptor.iconSize), uint(mapDescriptor.iconSize), iconPngData, resize.Bicubic)
	// Set up gg context
	dc := gg.NewContextForRGBA(imageDataBig)
	positions := []image.Point{}
	textPlots := []TextPlot{}
	for _ /*index*/, item := range locationsData.locationsObject {
		if strings.HasPrefix(item["PrefabName"].(string), plot.PrefabName) {
			posString := item["Position"].(string)
			xpos, ypos, zpos := parsePosition(posString)
			imageXPos := int((float64(imageWidth/2) + (xpos / scaleFactor)) - float64(mapDescriptor.iconSize/2))
			imageZPos := int((float64(imageWidth/2) - (zpos / scaleFactor)) - float64(mapDescriptor.iconSize/2))
			Debug(fmt.Sprintf("%s %s %f, %f, %f ==> (%f) {%d}:{%d}, '%s'(%v)", plot.PrefabName, posString, xpos, ypos, zpos, scaleFactor, imageXPos, imageZPos, plot.IconColor, plot.GetIconColor()))
			if plot.HighlightColor != "" {
				dc.SetColor(Black)
				dc.DrawEllipse(
					float64(imageXPos)+float64(mapDescriptor.iconSize/2)+2,
					float64(imageZPos)+float64(mapDescriptor.iconSize/2)+2,
					(float64(mapDescriptor.iconSize)/2)+5,
					(float64(mapDescriptor.iconSize)/2)+5)
				dc.Fill()
				dc.SetColor(plot.GetHighlightColor())
				dc.DrawEllipse(
					float64(imageXPos)+float64(mapDescriptor.iconSize/2),
					float64(imageZPos)+float64(mapDescriptor.iconSize/2),
					(float64(mapDescriptor.iconSize)/2)+5,
					(float64(mapDescriptor.iconSize)/2)+5)
				dc.Fill()
			}
			if plot.Text != "" {
				textPlot := TextPlot{
					text: plot.Text,
					x:    float64(imageXPos), // - (mapDescriptor.iconSize / 2)),
					z:    float64(imageZPos), // + (mapDescriptor.iconSize * 2)),
				}
				textPlots = append(textPlots, textPlot)
			}
			positions = append(positions, image.Point{X: imageXPos, Y: imageZPos})
		}
	}

	for _, pos := range positions {
		iconRect := image.Rectangle{Min: pos, Max: pos.Add(iconPngDataSizedTinted.Bounds().Size())}
		draw.Draw(imageDataBig, iconRect, iconPngDataSizedTinted, image.Point{}, draw.Over)
	}
	return textPlots
}

func processTiles(config *AppConfig, mapDescriptor *MapDescriptor, imageDataBig *image.RGBA) {
	if config.SkipMapGeneration {
		pngData, _ := os.Open(mapDescriptor.MapOutputPath())
		imageData, _ := png.Decode(pngData)
		// TODO: err handling
		draw.Draw(imageDataBig, imageData.Bounds(), imageData, image.Point{}, draw.Over)
		return
	}

	// Process Bitmaps Multi-threaded
	runtime.GOMAXPROCS(runtime.NumCPU()) // Ensure max parallelism
	var wg sync.WaitGroup
	wg.Add(mapDescriptor.tileSideCount * mapDescriptor.tileSideCount)
	mu = &sync.Mutex{}

	bar := progressbar.NewOptions(mapDescriptor.tileSideCount*mapDescriptor.tileSideCount,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(fmt.Sprintf("Processing %d Bitmaps", mapDescriptor.tileSideCount*mapDescriptor.tileSideCount)),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)
	for x := 0; x < mapDescriptor.tileSideCount; x++ {
		for z := 0; z < mapDescriptor.tileSideCount; z++ {
			go func(x int, z int) {
				defer wg.Done()
				processBitmap(mapDescriptor, x, z, imageDataBig, config.ForceReload)
				bar.Add(1)
			}(x, z)
		}
	}
	wg.Wait()
	bar.Close()
	bar.Clear()

	// Write Big PNG map
	SavePng(mapDescriptor.MapOutputPath(), imageDataBig, true)
}

func processLocations(config *AppConfig, mapDescriptor *MapDescriptor, imageDataBig *image.RGBA, locationsData *LocationsData) {
	// Plot Locations Icons
	InitIcons()
	viper.UnmarshalKey("plots", &config.Plots)
	bar := progressbar.NewOptions(len(config.Plots),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(fmt.Sprintf("Processing %d Plots over %d Locations", len(config.Plots), len(locationsData.locationsObject))),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)
	var textPlots = []TextPlot{}
	for _ /*index*/, plot := range config.Plots {
		tps := plotLocations(mapDescriptor, imageDataBig, locationsData, &plot)
		textPlots = append(textPlots, tps...)
		bar.Add(1)
	}
	bar.Close()

	// Draw Gathered Up Text Plots
	fmt.Printf("Plotting %d text plots.\n", len(textPlots))
	dc := gg.NewContextForRGBA(imageDataBig)
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{Size: 20})
	dc.SetFontFace(face)
	for _, textPlot := range textPlots {
		dc.SetColor(Black)
		dc.DrawString(textPlot.text, textPlot.x-10, textPlot.z+2)
		dc.SetColor(White)
		dc.DrawString(textPlot.text, textPlot.x-12, textPlot.z+0)
	}
}

func Process(config *AppConfig) { //basePath string, iconsPath string, forceReload bool, plots []Plot) int {
	mapDescriptor := ReadMapDescriptor(config.BasePath(), config.IconsPath)
	if mapDescriptor == nil {
		return
	}
	locationsData := ReadLocationsData(config.BasePath())
	Debug(fmt.Sprintf("MapDescriptor: %v\nLocationsData Count: %d", mapDescriptor, len(locationsData.locationsObject)))

	// Create Big Blank Image
	wh := int(mapDescriptor.ImageWidth())
	imageDataBig := image.NewRGBA(image.Rect(0, 0, wh, wh))

	processTiles(config, mapDescriptor, imageDataBig)
	processLocations(config, mapDescriptor, imageDataBig, locationsData)

	// Write map with icons
	SavePng(mapDescriptor.MapWithIconsOutputPath(), imageDataBig, true)
	fmt.Printf("Done!\nMap output to: %s\nMap with icons output to: %s\n", mapDescriptor.MapOutputPath(), mapDescriptor.MapWithIconsOutputPath())
}
