package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gpestana/kapacitor-unit/io"
	"github.com/gpestana/kapacitor-unit/model"
	"log"
	"strings"
)


func main() {
	fmt.Println(renderWelcome())

	f := model.LoadAppConfig()
	kapacitor := io.NewKapacitor(f.KapacitorHost)
	influxdb := io.NewInfluxdb(f.InfluxdbHost)

	tests, err := model.LoadTestsFromFileSystem(*f)
	if err != nil {
		log.Fatal("Error loading model configurations: ", err)
	}

	// Validates, runs tests in series and print results
	for _, t := range tests {
		if err := t.Validate(); err != nil {
			log.Printf("Error validating model: Error: %v Test: %+v  ", err, t)
			continue
		}

		// Runs model
		err = t.Run(kapacitor, influxdb)
		if err != nil {
			log.Println("Error running model: ", t, " Error: ", err)
			continue
		}
		//Prints model output
		setColor(t)
		log.Println(t)
		color.Unset()
	}
}

// Sets output color based on model results
func setColor(t model.TestInterface) {
	if t.Passed() == true {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgRed)
	}
}

func renderWelcome() string {
	logo := make([]string, 9)
	logo[0] = "  _                          _ _                                _ _            "
	logo[1] = " | |                        (_) |                              (_) |           "
	logo[2] = " | | ____ _ _ __   __ _  ___ _| |_ ___  _ __ ______ _   _ _ __  _| |_          "
	logo[3] = " | |/ / _` | '_ \\ / _` |/ __| | __/ _ \\| '__|______| | | | '_ \\| | __|      "
	logo[4] = " |   < (_| | |_) | (_| | (__| | || (_) | |         | |_| | | | | | |_          "
	logo[5] = " |_|\\_\\__,_| .__/ \\__,_|\\___|_|\\__\\___/|_|          \\__,_|_| |_|_|\\__| "
	logo[6] = "           | |                                                                 "
	logo[7] = "           |_|                                                        		      "
	logo[8] = "The unit model framework for TICK scripts (v0.8)\n"
	return strings.Join(logo, "\n")
}
