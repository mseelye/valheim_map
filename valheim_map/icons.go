package valheim_map

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

func IconNames() {
	fmt.Println(`
All icons come from https://game-icons.net. Great site, thank you!

I have found that all icons work best with the following (ymmv):
	Background: none (different than default)
	Forground: 
		Gradient: plain, color: white
		Stroke: enabled, color: black, width: 8 (different than default)
	Size & Preset:
		Dimensions: 512px

You can obviously replace these with any icons you want to use and configure
config.json to use whatever you want.

aqueduct.png           https://game-icons.net/1x1/delapouite/aqueduct.html
bird-claw.png          https://game-icons.net/1x1/lorc/bird-claw.html
broken-skull.png       https://game-icons.net/1x1/lorc/broken-skull.html
cash.png               https://game-icons.net/1x1/lorc/cash.html
castle.png             https://game-icons.net/1x1/delapouite/castle.html
castle-ruins.png       https://game-icons.net/1x1/delapouite/castle-ruins.html
cave-entrance.png      https://game-icons.net/1x1/delapouite/cave-entrance.html
crypt-entrance.png     https://game-icons.net/1x1/delapouite/crypt-entrance.html
death-skull.png        https://game-icons.net/1x1/sbed/death-skull.html
dinosaur-bones.png     https://game-icons.net/1x1/lorc/dinosaur-bones.html
drakkar.png            https://game-icons.net/1x1/delapouite/drakkar.html
dungeon-gate.png       https://game-icons.net/1x1/delapouite/dungeon-gate.html
fire.png               https://game-icons.net/1x1/sbed/fire.html
fishing-lure.png       https://game-icons.net/1x1/delapouite/fishing-lure.html
goblin-camp.png        https://game-icons.net/1x1/delapouite/goblin-camp.html
gooey-daemon.png       https://game-icons.net/1x1/lorc/gooey-daemon.html
grapes.png             https://game-icons.net/1x1/lorc/grapes.html
half-body-crawling.png https://game-icons.net/1x1/delapouite/half-body-crawling.html
impact-point.png       https://game-icons.net/1x1/lorc/impact-point.html
jeweled-chalice.png    https://game-icons.net/1x1/lorc/jeweled-chalice.html
lighthouse.png         https://game-icons.net/1x1/delapouite/lighthouse.html
mining.png             https://game-icons.net/1x1/lorc/mining.html
nest-eggs.png          https://game-icons.net/1x1/delapouite/nest-eggs.html
oak-leaf.png           https://game-icons.net/1x1/delapouite/oak-leaf.html
ogre.png               https://game-icons.net/1x1/delapouite/ogre.html
open-chest.png         https://game-icons.net/1x1/skoll/open-chest.html
ore.png                https://game-icons.net/1x1/faithtoken/ore.html
pointy-hat.png         https://game-icons.net/1x1/lorc/pointy-hat.html
raise-skeleton.png     https://game-icons.net/1x1/skoll/raise-skeleton.html
ribcage.png            https://game-icons.net/1x1/lorc/ribcage.html
sewing-string.png      https://game-icons.net/1x1/delapouite/sewing-string.html
shard-sword.png        https://game-icons.net/1x1/lorc/shard-sword.html
stone-pile.png         https://game-icons.net/1x1/delapouite/stone-pile.html
stone-tower.png        https://game-icons.net/1x1/lorc/stone-tower.html
sword-spade.png        https://game-icons.net/1x1/lorc/sword-spade.html
tick.png               https://game-icons.net/1x1/lorc/tick.html
tower-fall.png         https://game-icons.net/1x1/lorc/tower-fall.html
tree-roots.png         https://game-icons.net/1x1/delapouite/tree-roots.html
tree-roots-vf.png    (same as above but flipped vertically, fg: Trans.: vertical)
trunk-mushroom.png     https://game-icons.net/1x1/delapouite/trunk-mushroom.html
tusks-flag.png         https://game-icons.net/1x1/delapouite/tusks-flag.html
viking-helmet.png      https://game-icons.net/1x1/delapouite/viking-helmet.html
village.png            https://game-icons.net/1x1/delapouite/village.html
wood-cabin.png         https://game-icons.net/1x1/delapouite/wood-cabin.html
wood-canoe.png         https://game-icons.net/1x1/delapouite/wood-canoe.html`)
}

// Embed all the icons in the icons/ directory
//
//go:embed icons/*.png
var iconsFS embed.FS

// Embed the "missing icon" from the hazard-sign.png
//
//go:embed hazard-sign.png
var missingIcon []byte

// map containing data for all the embedded icons, by png filename, with extension
var IconData = make(map[string][]byte)

func InitIcons() {
	files, _ := iconsFS.ReadDir("icons")
	for _, file := range files {
		// Skip directories and files not named .png
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".png") {
			continue
		}
		// Note: do not use `filepath.Join` with embedded, as embedded ALWAYS works like unix using /
		// filePath := filepath.Join("icons", file.Name())
		// Use `/`, NOT filepath.Join, embed works like unix, even on windowns
		filePath := "icons/" + file.Name()
		data, err := iconsFS.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}
		IconData[file.Name()] = data
	}
}

func OpenIcon(mapDescriptor *MapDescriptor, iconPng string) (image.Image, error) {
	pngPath := mapDescriptor.IconPath(iconPng)
	pngFile, err := os.Open(pngPath)
	if err == nil {
		// We found the icon png at the specified path in actual FS, use that.
		defer pngFile.Close()
		pngData, err := png.Decode(pngFile)
		if err != nil {
			pngData, _ := png.Decode(bytes.NewReader(missingIcon))
			return pngData, fmt.Errorf("could not decode png file '%s' (%s)", pngPath, err.Error())
		}
		return pngData, nil
	} else {
		// We did not find the icon png at the specified path in actual FS, try embedded.
		if data, exists := IconData[iconPng]; exists {
			pngData, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				pngData, _ := png.Decode(bytes.NewReader(missingIcon))
				return pngData, fmt.Errorf("could not decode embedded png '%s' (%s)", iconPng, err.Error())
			}
			return pngData, nil
		} else {
			// We did not find it again, use the "missing icon" icon but also return error.
			pngData, _ := png.Decode(bytes.NewReader(missingIcon))
			return pngData, fmt.Errorf("cannot find icon, '%s' in icon path '%s' nor in embedded icons", iconPng, mapDescriptor.IconPath(iconPng))
		}
	}
}
