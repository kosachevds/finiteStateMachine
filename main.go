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
	machine, err := machine.ReadFromFile(files[index])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Welcome to the machine!")
	for {
		var input string
		fmt.Scanln(&input)
		fmt.Println(machine.IsCanHandle(input))
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
