package main

import (
	"finiteStateMachine/machine"
	"fmt"
	"strconv"
)

const testMachineIndex = 2

func main() {
	files, err := machine.ReadDirMachines("./data")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = printStrings(files, "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	//index := 2
	fmt.Println("Write index:")
	index, err := scanIndex()
	if err != nil {
		fmt.Println(err)
		return
	}
	fsm, err := machine.ReadFromFile(files[index])
	if err != nil {
		fmt.Println(err)
		return
	}
	if index == testMachineIndex {
		if checkPredefinedMachine(fsm) {
			fmt.Println("Good predefined machine!")
		} else {
			fmt.Println("Bad predefined machine!")
		}
	}
	fmt.Println("Welcome to the machine!")
	for {
		var input string
		fmt.Scanln(&input)
		fmt.Println(fsm.IsCanHandle(input))
	}
}

func scanIndex() (int, error) {
	var input string
	fmt.Scanln(&input)
	return strconv.Atoi(input)
}

func printStrings(filenames []string, delimiter string) error {
	for i, name := range filenames {
		_, err := fmt.Printf("%d: %s%s", i, name, delimiter)
		if err != nil {
			return fmt.Errorf("writing machines error: %v", err)
		}
	}
	return nil
}

func checkPredefinedMachine(fsm machine.FiniteStateMachine) bool {
	// TODO: unit-tests
	inputs := [...]string{
		"acdaf",
		"acdd",
		"ae",
		"aeb",
		"ab",
		"acf",
		"adaf",
		"add",
		"abf",
	}
	for _, input := range inputs {
		if !fsm.IsCanHandle(input) {
			return false
		}
	}
	return true
}
