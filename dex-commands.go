package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DexCommand struct {
	name        string
	description string
	display     bool
	callback    func(*CommandConfig) error
}

type CommandConfig struct {
	baseUrl      string
	previousUrl  *string
	nextUrl      *string
	locationArea string
}

func GetConfig() *CommandConfig {
	return &CommandConfig{
		baseUrl:      "https://pokeapi.co/api/v2/",
		previousUrl:  nil,
		nextUrl:      nil,
		locationArea: "location-area/",
	}
}

func GetCommands() map[string]DexCommand {
	return map[string]DexCommand{
		"help": {
			name:        "help",
			description: "Displays this help message",
			display:     true,
			callback:    CommandHelp,
		},
		"exit": {
			name:        "exit/quit",
			description: "Exit the Pokedex",
			display:     true,
			callback:    CommandExit,
		},
		"quit": {
			name:        "exit/quit",
			description: "Exit the Pokedex",
			display:     false,
			callback:    CommandExit,
		},
		"map": {
			name:        "map",
			description: "List locations/next page",
			display:     true,
			callback:    CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List locations/previous page",
			display:     true,
			callback:    CommandMapb,
		},
	}
}

func CommandHelp(config *CommandConfig) error {
	fmt.Print("\n")
	commands := GetCommands()
	for cmdKey := range commands {
		cmd := commands[cmdKey]
		if !cmd.display {
			continue
		}

		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Print("\n")

	return nil
}

func CommandExit(config *CommandConfig) error {
	fmt.Println("\nClosing Pokedex...")
	os.Exit(0)
	return nil
}

func CommandMap(config *CommandConfig) error {
	var mapUrl string
	if config.nextUrl == nil {
		mapUrl = config.baseUrl + config.locationArea
	} else {
		mapUrl = *config.nextUrl
	}

	response, err := http.Get(mapUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fmt.Println(response.Body)

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//fmt.Printf("jsonData: %v\n", jsonData)

	locationArea := LocationArea{}
	//var locationArea map[string]interface{}
	if err = json.Unmarshal(jsonData, &locationArea); err != nil {
		return err
	}
	//fmt.Printf("locationArea: %v\n", locationArea)
	config.nextUrl = &locationArea.Next
	config.previousUrl = &locationArea.Previous

	//var results interface{} //[]map[string]string
	results := locationArea.Results
	//var rlocations []map[string]string
	//rlocations = results.([]map[string]string)
	//fmt.Printf("locationArea: %v\n", locationArea)
	fmt.Printf("results: %v\n", results)
	for _, result := range results {
		name := result["name"]
		fmt.Printf("%s\n", name)
	}

	return nil
}

func CommandMapb(config *CommandConfig) error {
	var mapUrl string
	if config.previousUrl == nil {
		return fmt.Errorf("There is no previous page of locations.\n")
	} else if *config.previousUrl == "" {
		//mapUrl = config.baseUrl + config.locationArea
		return fmt.Errorf("Already at first page of locations.\n")
	} else {
		mapUrl = *config.previousUrl
	}

	response, err := http.Get(mapUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fmt.Println(response.Body)

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//fmt.Printf("jsonData: %v\n", jsonData)

	locationArea := LocationArea{}
	//var locationArea map[string]interface{}
	if err = json.Unmarshal(jsonData, &locationArea); err != nil {
		return err
	}
	//fmt.Printf("locationArea: %v\n", locationArea)
	config.nextUrl = &locationArea.Next
	config.previousUrl = &locationArea.Previous

	//var results interface{} //[]map[string]string
	results := locationArea.Results
	//var rlocations []map[string]string
	//rlocations = results.([]map[string]string)
	//fmt.Printf("locationArea: %v\n", locationArea)
	fmt.Printf("results: %v\n", results)
	for _, result := range results {
		name := result["name"]
		fmt.Printf("%s\n", name)
	}

	return nil
}
