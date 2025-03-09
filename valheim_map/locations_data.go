package valheim_map

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

// Current list of prefab names from valheim
var PossiblePrefabNames = []string{
	"AbandonedLogCabin02", "AbandonedLogCabin03", "AbandonedLogCabin04",
	"AshlandRuins",
	"BogWitch_Camp",
	"Bonemass",
	"CharredFortress",
	"CharredRuins1", "CharredRuins2", "CharredRuins3", "CharredRuins4",
	"CharredStone_Spawner",
	"CharredTowerRuins1_dvergr", "CharredTowerRuins1", "CharredTowerRuins2", "CharredTowerRuins3",
	"Crypt2", "Crypt3", "Crypt4",
	"Dolmen01", "Dolmen02", "Dolmen03",
	"Dragonqueen",
	"DrakeLorestone",
	"DrakeNest01",
	"Eikthyrnir",
	"FaderLocation",
	"FireHole",
	"FortressRuins",
	"GDKing",
	"GoblinCamp2",
	"GoblinKing",
	"Grave1",
	"Greydwarf_camp1",
	"Hildir_camp", "Hildir_cave", "Hildir_crypt", "Hildir_plainsfortress",
	"InfestedTree01",
	"LeviathanLava",
	"Mistlands_DvergrBossEntrance1", "Mistlands_DvergrTownEntrance1", "Mistlands_DvergrTownEntrance2",
	"Mistlands_Excavation1", "Mistlands_Excavation2", "Mistlands_Excavation3",
	"Mistlands_Giant1", "Mistlands_Giant2",
	"Mistlands_GuardTower1_new", "Mistlands_GuardTower1_ruined_new", "Mistlands_GuardTower1_ruined_new2", "Mistlands_GuardTower2_new", "Mistlands_GuardTower3_new", "Mistlands_GuardTower3_ruined_new",
	"Mistlands_Harbour1",
	"Mistlands_Lighthouse1_new",
	"Mistlands_RoadPost1",
	"Mistlands_RockSpire1",
	"Mistlands_Statue1", "Mistlands_Statue2",
	"Mistlands_StatueGroup1",
	"Mistlands_Swords1", "Mistlands_Swords2", "Mistlands_Swords3",
	"Mistlands_Viaduct1", "Mistlands_Viaduct2",
	"MorgenHole1", "MorgenHole2", "MorgenHole3",
	"MountainCave02",
	"MountainGrave01",
	"MountainWell1",
	"PlaceofMystery1", "PlaceofMystery2", "PlaceofMystery3",
	"Ruin1", "Ruin2", "Ruin3",
	"Runestone_Ashlands", "Runestone_BlackForest", "Runestone_Boars", "Runestone_Draugr", "Runestone_Greydwarfs", "Runestone_Meadows", "Runestone_Mistlands", "Runestone_Mountains", "Runestone_Plains", "Runestone_Swamps",
	"ShipSetting01",
	"ShipWreck01", "ShipWreck02", "ShipWreck03", "ShipWreck04",
	"StartTemple",
	"StoneCircle",
	"StoneHenge1", "StoneHenge2", "StoneHenge3", "StoneHenge4", "StoneHenge5", "StoneHenge6",
	"StoneHouse3", "StoneHouse4",
	"StoneTower1", "StoneTower3",
	"StoneTowerRuins03", "StoneTowerRuins04", "StoneTowerRuins05", "StoneTowerRuins07", "StoneTowerRuins08", "StoneTowerRuins09", "StoneTowerRuins10",
	"SulfurArch",
	"SunkenCrypt4",
	"SwampHut1", "SwampHut2", "SwampHut3", "SwampHut4", "SwampHut5", "SwampRuin1",
	"SwampRuin2",
	"SwampWell1",
	"TarPit1", "TarPit2", "TarPit3",
	"TrollCave02",
	"Vendor_BlackForest",
	"VoltureNest",
	"Waymarker01", "Waymarker02",
	"WoodFarm1",
	"WoodHouse1", "WoodHouse10", "WoodHouse11", "WoodHouse12", "WoodHouse13", "WoodHouse2", "WoodHouse3", "WoodHouse4", "WoodHouse5", "WoodHouse6", "WoodHouse7", "WoodHouse8", "WoodHouse9",
	"WoodVillage1",
}

type LocationsData struct {
	basePath        string
	locationsObject []map[string]interface{}
}

func NewLocationsData(basePath string) *LocationsData {
	return &LocationsData{basePath: basePath, locationsObject: []map[string]interface{}{}}
}

func (ld *LocationsData) Path() string {
	return path.Join(ld.basePath, "locations.json")
}

func ReadLocationsData(basePath string) *LocationsData {
	locationsData := NewLocationsData(basePath)
	content, err := os.ReadFile(locationsData.Path())
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("You must specify the data folder as an argument, or run this from the data folder.")
		return nil
	} else {
		err := json.Unmarshal(content, &locationsData.locationsObject)
		if err != nil {
			fmt.Printf("Cannot parse json file, '%s'. Error: %v\n", locationsData.Path(), err)
		}
		return locationsData
	}
}

func OutputOrphanPrefabNames(config *AppConfig) {
	viper.UnmarshalKey("plots", &config.Plots)
	prefabNames := map[string]string{}
	for _, plot := range config.Plots {
		prefabNames[plot.PrefabName] = plot.PrefabName
	}
	var missingNames []string
	for _, name := range PossiblePrefabNames {
		found := false
		for usedName := range prefabNames {
			if strings.HasPrefix(name, usedName) {
				found = true
				break
			}
		}
		if !found {
			missingNames = append(missingNames, name)
		}
	}
	prefix := ""
	prefixLength := 6
	fmt.Printf("PrefabName values(%d) not found in config(%s):", len(missingNames), viper.ConfigFileUsed())
	for _, missingName := range missingNames {
		if len(missingName) < prefixLength || missingName[:prefixLength] != prefix {
			prefix = missingName[:min(prefixLength, len(missingName))]
			fmt.Print("\n\t")
		}
		fmt.Print(missingName, ", ")
	}
	fmt.Println()
}
