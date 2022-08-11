package config

import (
	"log"
	"strings"

	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/fs"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/logger"
	"github.com/1995parham-teaching/fandogh/internal/metric"
	telemetry "github.com/1995parham-teaching/fandogh/internal/telemetry/config"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "fandogh_"
)

type (
	// Config holds all configurations.
	Config struct {
		Database    db.Config        `koanf:"database"`
		FileStorage fs.Config        `koanf:"file_storage"`
		Monitoring  metric.Config    `koanf:"monitoring"`
		Logger      logger.Config    `koanf:"logger"`
		Telemetry   telemetry.Config `koanf:"telemetry"`
		JWT         jwt.Config       `koanf:"jwt"`
	}
)

// New reads configuration with viper.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}
