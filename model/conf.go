package model

import (
	"flag"
	"log"
)

type Config struct {
	//Path for model definitions YAML file
	TestsPath string
	// Path for directory where TICKscripts are
	ScriptsDir    string
	InfluxdbHost  string
	KapacitorHost string
}

func LoadAppConfig() *Config {
	influxdbHost := flag.String("influxdb", "http://localhost:8086",
		"InfluxDB host")
	kapacitorHost := flag.String("kapacitor", "http://localhost:9092",
		"Kapacitor host")
	testsPath := flag.String("tests", "", "Tests definition file")
	scriptsDir := flag.String("dir", "", "TICKscripts directory")

	flag.Parse()

	if *testsPath == "" {
		log.Fatal("ERROR: Path for tests definitions (--tests) must be defined")
	}

	if *scriptsDir == "" {
		log.Fatal("ERROR: Path for where TICKscripts directory (--dir) must be defined")
	}

	config := Config{*testsPath, *scriptsDir, *influxdbHost, *kapacitorHost}

	return &config
}



// TestConf is an intermediary struct used for decoding model configurations
// before constructing concrete Test implementations.
// This prevents a lot of additional type assertion from having to use a
// general interface for our json / yaml data. An alternative could be something like:
// https://github.com/mitchellh/mapstructure
type TestConf struct {
	Name       string
	Db         string
	Rp         string
	SeedDb     string       `yaml:"seed_db,omitempty"`
	SeedRp     string       `yaml:"seed_rp,omitempty"`
	AssertDb   string       `yaml:"assert_db,omitempty"`
	AssertRp   string       `yaml:"assert_rp,omitempty"`
	SeedData   []string     `yaml:"seed_data,omitempty"`
	AssertData []string     `yaml:"assert_data,omitempty"`
	TaskName   string       `yaml:"task_name,omitempty"`
	Expects    Result `yaml:"expects"`
	Data       []string
	Type       string `yaml:"type"`
	TestType   string `yaml:"test_type"`
}

