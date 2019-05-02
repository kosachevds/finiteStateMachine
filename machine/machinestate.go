package machine

import (
	"fmt"
	"sort"
)

type machineState struct {
	code        int
	transitions map[rune]*machineState
}

func newMachineState(code int) *machineState {
	return &machineState{code, make(map[rune]*machineState)}
}

func createStatesTree(rules []transitionRule) (*machineState, error) {
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].beginState < rules[j].beginState
	})
	addedStates := make(map[int]*machineState)
	for i := 0; i < len(rules); i++ {
		stateCode := rules[i].beginState
		var j int
		for j := i + 1; j < len(rules); j++ {
			if rules[i].beginState != rules[j].beginState {
				break
			}
		}
		state, ok := addedStates[stateCode]
		if !ok {
			state = newMachineState(stateCode)
			addedStates[stateCode] = state
		}
		for _, rule := range rules[i:(j - i)] {
			nextState, ok := addedStates[rule.nextState]
			if !ok {
				nextState = newMachineState(stateCode)
				addedStates[rule.nextState] = nextState
			}
			if _, ok = state.transitions[rule.symbol]; ok {
				return nil, fmt.Errorf("transition rules are not deterministic")
			}
			state.transitions[rule.symbol] = nextState
		}
		i = j
	}
	return addedStates[0], nil
}
