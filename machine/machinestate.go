package machine

import (
	"fmt"
	"sort"
)

type machineState struct {
	code int
	ways map[rune]*machineState
}

type allStatesSet struct {
	states map[int]*machineState
}

func newMachineState(code int) *machineState {
	return &machineState{code, make(map[rune]*machineState)}
}

func (ms *machineState) isFinal() bool {
	return len(ms.ways) == 0
}

func newStatesSet() *allStatesSet {
	return &allStatesSet{states: make(map[int]*machineState)}
}

func (all *allStatesSet) get(stateCode int) *machineState {
	state, ok := all.states[stateCode]
	if !ok {
		state = newMachineState(stateCode)
		all.states[stateCode] = state
	}
	return state
}

func createStatesGraph(rules []transitionRule) (*machineState, error) {
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].beginState < rules[j].beginState
	})
	addedStates := newStatesSet()
	for i := 0; i < len(rules); i++ {
		stateCode := rules[i].beginState
		j := i + 1
		for ; j < len(rules); j++ {
			if rules[i].beginState != rules[j].beginState {
				break
			}
		}
		state := addedStates.get(stateCode)
		for _, rule := range rules[i:j] {
			nextStateCode := rule.nextState
			nextState := addedStates.get(nextStateCode)
			if _, ok := state.ways[rule.symbol]; ok {
				return nil, fmt.Errorf("transition rules are not deterministic")
			}
			state.ways[rule.symbol] = nextState
		}
		i = j
	}
	return addedStates.get(rules[0].beginState), nil
}
