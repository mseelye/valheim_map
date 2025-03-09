package valheim_map

import (
	"image/color"
)

type Plot struct {
	PrefabName     string
	IconPng        string
	IconColor      string
	HighlightColor string
	Text           string
}

func NewPlot(prefabName string, iconPngFileName string, iconColor string, highlightColor string, text string) Plot {
	return Plot{
		PrefabName:     prefabName,
		IconPng:        iconPngFileName,
		IconColor:      iconColor,
		HighlightColor: highlightColor,
		Text:           text,
	}
}

func (p *Plot) GetIconColor() color.RGBA {
	return GetColor(p.IconColor)
}

func (p *Plot) GetHighlightColor() color.RGBA {
	return GetColor(p.HighlightColor)
}
