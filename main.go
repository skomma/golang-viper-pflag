package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel   string      `mapstructure:"loglevel"`
	ListenPort int16       `mapstructure:"port"`
	TopString  string      `mapstructure:"top_string"`
	Child      ChildConfig `mapstructure:"child"`
}

type ChildConfig struct {
	Enabled    bool
	Strings    []string
	GrandChild GrandChildConfig `mapstructure:"grand_child"`
}

type GrandChildConfig struct {
	String string
}

func initViper() {
	// load from config/config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	// set environment variable processing
	viper.SetEnvPrefix("viper_pflag")
	// "_" "-" in environment variable's key is treated same as "." when viper.Get is used
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// do same as pflag.CommandLine
// This function exists because invoking SetNormalizeFunc
func parseFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	flagSet.String("loglevel", "info", "Log level")
	flagSet.Int16P("port", "p", 9000, "A port number that this program will listen on")
	flagSet.String("top-string", "", "Top-level string")

	flagSet.Bool("child.enabled", false, "Enable child")
	flagSet.StringSlice("child.strings", nil, "Child strings")

	flagSet.String("child.grand-child.string", "", "Grand-child string")

	// do same as pflags.Parse()
	flagSet.Parse(os.Args[1:])

	// all "-" options are converted to "_" to match viper.Get parameter
	// see: https://github.com/spf13/pflag#mutating-or-normalizing-flag-names
	wordSepNormalizeFunc := func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "-", "_"))
	}
	flagSet.SetNormalizeFunc(wordSepNormalizeFunc)

	return flagSet
}

func main() {
	initViper()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Warn("Cannot find config.yaml, ignoring.")
		} else {
			// Config file was found but another error was produced
		}
	}

	flagSet := parseFlags()
	viper.BindPFlags(flagSet)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Infof("cfg: %#v", cfg)
}
