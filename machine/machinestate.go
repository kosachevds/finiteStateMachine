package machine

import (
	"fmt"
	"sort"
)

type machineState struct {
	code int
	ways map[rune]*machineState
}

func (ms *machineState) isFinal() bool {
	return len(ms.ways) == 0
}

func newMachineState(code int) *machineState {
	return &machineState{code, make(map[rune]*machineState)}
}

func createStatesGraph(rules []transitionRule) (*machineState, error) {
	// Fix: this function
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].beginState < rules[j].beginState
	})
	addedStates := make(map[int]*machineState)
	for i := 0; i < len(rules); i++ {
		stateCode := rules[i].beginState
		j := i + 1
		for ; j < len(rules); j++ {
			if rules[i].beginState != rules[j].beginState {
				break
			}
		}
		state, ok := addedStates[stateCode]
		if !ok {
			state = newMachineState(stateCode)
			addedStates[stateCode] = state
		}
		for _, rule := range rules[i:j] {
			nextStateCode := rule.nextState
			nextState, ok := addedStates[nextStateCode]
			if !ok {
				nextState = newMachineState(nextStateCode)
				addedStates[rule.nextState] = nextState
			}
			if _, ok = state.ways[rule.symbol]; ok {
				return nil, fmt.Errorf("transition rules are not deterministic")
			}
			state.ways[rule.symbol] = nextState
		}
		i = j
	}
	return addedStates[rules[0].beginState], nil
}
