package model

import (
	"errors"
	"fmt"
	"github.com/gpestana/kapacitor-unit/io"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TestInterface interface {
	Init(conf TestConf) error
	Run(kapacitor io.Kapacitor, influxdb io.Influxdb) error
	Validate() error
	GetName() string
	GetType() string
	SetTask(task Task)
	Passed() bool
}

type TestCollection []TestInterface

// NewTest is a factory responsible for instantiating the proper model type
func NewTest(t TestConf, config Config) (TestInterface, error) {

	if t.TestType == "" { //for backwards compatibility
		t.TestType = t.Type
	}

	task, err := NewTask(t, config)
	if err != nil {
		return nil, err
	}

	switch t.TestType {
	case "influxdbout":
		test := DbOutTest{}
		err := test.Init(t)
		test.Task = *task
		return &test, err

	case "stream", "batch":
		test := AlertTest{}
		err := test.Init(t)
		test.Task = *task
		return &test, err
	}

	return nil, errors.New(fmt.Sprintf("not supported configuration type: %s in configuration: %+v", t.Type, t))
}


func LoadTestCollectionFromYamlFile(fileName string, config Config) (TestCollection, error) {

	type conf struct {
		Tests []TestConf
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	c := conf{}

	err = yaml.Unmarshal(b, &c)

	if err != nil {
		return nil, err
	}

	tc := make(TestCollection, 0)

	for _, t := range c.Tests {
		newTest, err := NewTest(t, config)
		if err != nil {
			return nil, err
		}
		tc = append(tc, newTest)
	}

	return tc, nil

}

//Opens and parses model configuration file into a structure
func LoadTestsFromFileSystem(config Config) (TestCollection, error) {

	fileName := config.TestsPath

	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	if stat.IsDir() {
		filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
			if ext := filepath.Ext(path); ext == ".yml" || ext == ".yaml" {
				files = append(files, path)
			}
			return nil
		})

	} else {
		files = append(files, fileName)
	}

	tests := make(TestCollection, 0)

	for _, file := range files {
		fileTests, err := LoadTestCollectionFromYamlFile(file, config)
		if err != nil {
			return nil, err
		}
		tests = append(tests, fileTests...)
	}

	return tests, nil
}
