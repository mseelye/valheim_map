package valheim_map

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"

	"github.com/spf13/viper"
)

// https://pkg.go.dev/golang.org/x/image/colornames
var (
	NilColor    = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgb(0, 0, 0, 0)
	Black       = color.RGBA{0x00, 0x00, 0x00, 0xff} // rgb(0, 0, 0)
	Green       = color.RGBA{0x00, 0x80, 0x00, 0xff} // rgb(0, 128, 0)
	Seagreen    = color.RGBA{0x2e, 0x8b, 0x57, 0xff} // rgb(46, 139, 87)
	Gray        = color.RGBA{0x33, 0x33, 0x33, 0xff} // rgb(51, 51, 51)
	Lightgray   = color.RGBA{0xd3, 0xd3, 0xd3, 0xff} // rgb(211, 211, 211)
	Darkgreen   = color.RGBA{0x00, 0x64, 0x00, 0xff} // rgb(0, 100, 0)
	Yellow      = color.RGBA{0xff, 0xff, 0x00, 0xff} // rgb(255, 255, 0)
	Red         = color.RGBA{0xff, 0x00, 0x00, 0xff} // rgb(255, 0, 0)
	White       = color.RGBA{0xff, 0xff, 0xff, 0xff} // rgb(255, 255, 255)
	Cyan        = color.RGBA{0x00, 0xff, 0xff, 0xff} // rgb(0, 255, 255)
	Darkblue    = color.RGBA{0x00, 0x00, 0x8b, 0xff} // rgb(0, 0, 139)
	Purple      = color.RGBA{0x80, 0x00, 0x80, 0xff} // rgb(128, 0, 128)
	Blue        = color.RGBA{0x00, 0x00, 0xff, 0xff} // rgb(0, 0, 255)
	Lime        = color.RGBA{0x00, 0xff, 0x00, 0xff} // rgb(0, 255, 0)
	Darkgrey    = color.RGBA{0xa9, 0xa9, 0xa9, 0xff} // rgb(169, 169, 169)
	Gold        = color.RGBA{0xff, 0xd7, 0x00, 0xff} // rgb(255, 215, 0)
	Darkbrown   = color.RGBA{0x2b, 0x1d, 0x0e, 0xff} // rgb(43, 29, 14)
	Brown       = color.RGBA{0xa5, 0x2a, 0x2a, 0xff} // rgb(165, 42, 42)
	Saddlebrown = color.RGBA{0x8b, 0x45, 0x13, 0xff} // rgb(139, 69, 19)
	Lightblue   = color.RGBA{0xad, 0xd8, 0xe6, 0xff} // rgb(173, 216, 230)
)

var DefaultColors = map[string]color.RGBA{
	"NilColor":    NilColor,
	"Black":       Black,
	"Green":       Green,
	"Seagreen":    Seagreen,
	"Gray":        Gray,
	"Lightgray":   Lightgray,
	"Darkgreen":   Darkgreen,
	"Yellow":      Yellow,
	"Red":         Red,
	"White":       White,
	"Cyan":        Cyan,
	"Darkblue":    Darkblue,
	"Purple":      Purple,
	"Blue":        Blue,
	"Lime":        Lime,
	"Darkgrey":    Darkgrey,
	"Gold":        Gold,
	"Darkbrown":   Darkbrown,
	"Brown":       Brown,
	"Saddlebrown": Saddlebrown,
	"Lightblue":   Lightblue,
}

func ColorAverage(c1 color.RGBA, c2 color.RGBA) color.RGBA {
	return color.RGBA{
		(c1.R + c2.R) / 2,
		(c1.G + c2.G) / 2,
		(c1.B + c2.B) / 2,
		255}
}

func ColorDarken(c color.RGBA, level float32) color.RGBA {
	// 'level' is how dark to make the color.
	// Each RGB component is divided by 'level'.
	return color.RGBA{
		uint8(Clamp(float32(c.R)/level, 0, 255)),
		uint8(Clamp(float32(c.G)/level, 0, 255)),
		uint8(Clamp(float32(c.B)/level, 0, 255)),
		uint8(c.A)}
}

// TintImage tints the white-ish areas of an image with a given color.
func TintImage(img image.Image, tint color.RGBA, threshold float64) *image.RGBA {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	draw.Draw(out, bounds, img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			origColor := img.At(x, y)
			r, g, b, a := origColor.RGBA()

			// Convert to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Calculate brightness (perceived luminance)
			brightness := math.Sqrt(0.299*float64(r8)*float64(r8) +
				0.587*float64(g8)*float64(g8) +
				0.114*float64(b8)*float64(b8))

			// If the brightness is above the threshold, apply tint
			if brightness > threshold {
				out.Set(x, y, BlendColor(r8, g8, b8, a8, tint))
			}
		}
	}

	return out
}

// BlendColor blends an original color with a tint based on alpha blending.
func BlendColor(r, g, b, a uint8, tint color.RGBA) color.RGBA {
	// Blend factor: the more transparent, the more it gets tinted
	alphaFactor := float64(a) / 255.0
	newR := uint8(float64(r)*(1-alphaFactor) + float64(tint.R)*alphaFactor)
	newG := uint8(float64(g)*(1-alphaFactor) + float64(tint.G)*alphaFactor)
	newB := uint8(float64(b)*(1-alphaFactor) + float64(tint.B)*alphaFactor)

	return color.RGBA{newR, newG, newB, a}
}

func ParseHexColor(hexString string) (color.RGBA, error) {
	hexString = strings.TrimPrefix(hexString, "#")
	var rgba color.RGBA
	switch len(hexString) {
	case 6: // "#RRGGBB"
		rgba.A = 255 // Default to fully opaque
		_, err := fmt.Sscanf(hexString, "%02x%02x%02x", &rgba.R, &rgba.G, &rgba.B)
		return rgba, err
	case 8: // "#RRGGBBAA"
		_, err := fmt.Sscanf(hexString, "%02x%02x%02x%02x", &rgba.R, &rgba.G, &rgba.B, &rgba.A)
		return rgba, err
	default:
		return color.RGBA{}, fmt.Errorf("invalid color format: %s", hexString)
	}
}

// note: inf loop possible like this, wee
func ParseNamedColor(nameString string) (color.RGBA, error) {
	if nameString == "" {
		return White, nil
	}

	var namedColors = viper.GetStringMapString("colors")
	if colorValue, ok := namedColors[strings.ToLower(nameString)]; ok {
		return GetColor(colorValue), nil
	}

	defaultColor, exists := DefaultColors[nameString]
	if exists {
		return defaultColor, nil
	}

	fmt.Printf("unknown color name: '%s'\n", nameString)
	return color.RGBA{}, fmt.Errorf("unknown color name: '%s'", nameString)
}

func GetColor(val string) color.RGBA {
	var rgba = color.RGBA{}
	if strings.HasPrefix(val, "#") {
		rgba, _ = ParseHexColor(val)
	} else {
		rgba, _ = ParseNamedColor(val)
	}
	return rgba
}
