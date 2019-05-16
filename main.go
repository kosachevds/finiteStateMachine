package main

import (
	"finiteStateMachine/machine"
	"fmt"
	"strconv"
)

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
	index := 2
	//fmt.Println("Write index:")
	//index, err := scanIndex()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	fsm, err := machine.ReadFromFile(files[index])
	if err != nil {
		fmt.Println(err)
		return
	}
	if index == 2 {
		if checkIndex2(fsm) {
			fmt.Println("Good 2nd machine!")
		} else {
			fmt.Println("Bad 2nd machine!")
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

func checkIndex2(fsm machine.FiniteStateMachine) bool {
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
