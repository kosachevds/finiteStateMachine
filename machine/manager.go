package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Manager struct {
	fsm       *finiteStateMachine
	filenames []string
}

const machineFileExtension string = ".txt"

func ReadDirMachines(dirname string) (*Manager, error) {
	fileInfoArray, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("reading directory error %v", err)
	}
	txtFiles := make([]string, 0, len(fileInfoArray))
	for _, item := range fileInfoArray {
		if strings.HasSuffix(item.Name(), machineFileExtension) {
			txtFiles = append(txtFiles, item.Name())
		}
	}
	return &Manager{nil, txtFiles}, nil
}

func (mh *Manager) WriteMachinesList(writer io.Writer) error {
	if mh == nil || len(mh.filenames) == 0 {
		return fmt.Errorf("empty machine handler or nil")
	}
	for i, name := range mh.filenames {
		_, err := fmt.Fprintf(writer, "%d: %s", i, name)
		if err != nil {
			return fmt.Errorf("writing machines error: %v", err)
		}
	}
	return nil
}

func (mh *Manager) InitMachineWithIndex(index int) error {
	fsm, err := ReadFromFile(mh.filenames[index])
	if err != nil {
		return fmt.Errorf("creating machine error: %v", err)
	}
	mh.fsm = fsm
	return nil
}

func (mh *Manager) CheckMachineInput(input string) bool {
	return mh.fsm.IsCanHandle(input)
}
