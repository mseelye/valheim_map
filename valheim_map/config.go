package valheim_map

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type TextPlot struct {
	text string
	x    float64
	z    float64
}

type AppConfig struct {
	ForceReload             bool
	Verbose                 bool
	HelpFlag                bool
	IconsPath               string
	IconNamesFlag           bool
	ConfigName              string
	WriteConfigFlag         bool
	SkipMapGeneration       bool
	OutputOrphanPrefabNames bool
	Args                    []string
	Plots                   []Plot
	TextPlots               []TextPlot
}

//go:embed embedded-default-config.json
var defaultConfig string

func readEmbeddedConfig() error {
	err := viper.ReadConfig(strings.NewReader(defaultConfig))
	if err != nil {
		println(err.Error())
	}
	return err
}

func processViperConfig(config *AppConfig) error {
	viper.SetConfigName(config.ConfigName)
	viper.SetConfigType("json")
	if config.WriteConfigFlag {
		readEmbeddedConfig()
		viper.WriteConfigAs("default-config.json")
		os.Exit(0)
	}
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Config file '%s.json' error: (%s), using embedded default config.\n", config.ConfigName, err)
		err = readEmbeddedConfig()
	}
	return err
}

// Parse flags into AppConfig, reads in config
func InitConfig() (*AppConfig, error) {
	config := &AppConfig{}
	flag.BoolVar(&config.ForceReload, "forcereload", false, "Force the reload and generation of the tile bitmaps.")
	flag.BoolVar(&config.Verbose, "verbose", false, "Log more debug information.")
	flag.BoolVar(&config.HelpFlag, "help", false, "Show detailed help information.")
	flag.StringVar(&config.IconsPath, "iconsPath", "./icons", "Path to a directory containing png icons.")
	flag.BoolVar(&config.IconNamesFlag, "icons", false, "Show list of icons included with this project.")
	flag.StringVar(&config.ConfigName, "configName", "config", "Name of config file to use (omit json extension!)")
	flag.BoolVar(&config.WriteConfigFlag, "defaultConfig", false, "Outputs the embedded default config file to default-config.json.")
	flag.BoolVar(&config.SkipMapGeneration, "skipMapGeneration", false, "If set existing map.png will be used, tile pngs will not be reprocessed.")
	flag.BoolVar(&config.OutputOrphanPrefabNames, "prefabCheck", false, "Outputs any PrefabName values NOT represented in the current config.")
	flag.Parse()
	config.Args = flag.Args()
	err := processViperConfig(config)
	return config, err
}

func (af *AppConfig) BasePath() string {
	if len(af.Args) > 0 {
		return af.Args[0]
	} else {
		return ""
	}
}
