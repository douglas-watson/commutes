package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config stores configuration parameters laoded from a YAML file
type Config struct {
	AccessToken string `yaml:"Access token"`

	Distance float64 `yaml:"Distance"`
	Duration int     `yaml:"Duration"`
	GearID   string  `yaml:"Gear ID"`

	StartDateStr    string     `yaml:"Start date"`
	EndDateStr      string     `yaml:"End date"`
	ExcludeDatesStr [][]string `yaml:"Exclude dates"`
	OnlyWeekdays    bool       `yaml:"Only weekdays"`

	StartDate    time.Time
	EndDate      time.Time
	ExcludeDates [][2]time.Time
}

// ConfigFromYaml loads a YAML config file. See `config.yaml` for example.
func ConfigFromYaml(path string) (Config, error) {
	var config Config

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return config, err
	}

	// Parse dates
	config.StartDate, err = time.Parse("2006-01-02 15:04", config.StartDateStr)
	config.EndDate, err = time.Parse("2006-01-02", config.EndDateStr)

	for _, pair := range config.ExcludeDatesStr {
		start, err := time.Parse("2006-01-02", pair[0])
		if err != nil {
			return config, err
		}
		end, err := time.Parse("2006-01-02", pair[1])
		if err != nil {
			return config, err
		}

		config.ExcludeDates = append(config.ExcludeDates, [...]time.Time{start, end})
	}

	return config, err
}
