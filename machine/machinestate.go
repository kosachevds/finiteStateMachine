package machine

type machineState struct {
	code        int
	transitions map[rune]*machineState
}

func newMachineState(code int) machineState {
	return machineState{code, make(map[rune]*machineState)}
}
