package main

import (
	"finiteStateMachine/machine"
	"fmt"
	"os"
	"strconv"
)

func main() {
	handler, err := machine.ReadDirMachines("./data")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = handler.WriteMachinesList(os.Stdout)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Write index:")
	index, err := scanIndex()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = handler.InitMachineWithIndex(index)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var input string
		fmt.Scanln(&input)
		fmt.Println(handler.CheckMachineInput(input))
	}
}

func scanIndex() (int, error) {
	var input string
	fmt.Scanln(&input)
	return strconv.Atoi(input)
}
