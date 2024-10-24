package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("\nWelcome to the Pokedex!!!\n\n")

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))

	commands := GetCommands()
	config := GetConfig()

	for true {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		cmd, ok := commands[input]
		if !ok {
			fmt.Printf("\n%s is not a Pokedex command.\n\n", input)
			continue
		}

		if err := cmd.callback(config); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}
