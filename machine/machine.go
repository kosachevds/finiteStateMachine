package machine

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const machineFileExtension = ".txt"

type FiniteStateMachine interface {
	IsCanHandle(string) bool
}

type statesGraph struct {
	root *machineState
}

func ReadFromFile(filename string) (FiniteStateMachine, error) {
	rules, err := readRules(filename)
	if err != nil {
		return nil, err
	}
	rules = determinateRules(rules)
	states, err := createStatesGraph(rules)
	if err != nil {
		return nil, fmt.Errorf("machine creation error %e", err)
	}
	return &statesGraph{
		root: states,
	}, nil

}

func ReadDirMachines(dirname string) ([]string, error) {
	fileInfoArray, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("reading directory error %v", err)
	}
	txtFiles := make([]string, 0, len(fileInfoArray))
	for _, item := range fileInfoArray {
		if !strings.HasSuffix(item.Name(), machineFileExtension) {
			continue
		}
		fullPath := path.Join(dirname, item.Name())
		txtFiles = append(txtFiles, fullPath)
	}
	return txtFiles, nil
}

func (sg *statesGraph) IsCanHandle(input string) bool {
	current := sg.root
	lastIndex := len(input) - 1
	for i, item := range input {
		nextState, ok := current.ways[item]
		if !ok {
			break
		}
		if i == lastIndex && nextState.isFinal {
			return true
		}
		current = nextState
	}
	return false
}
