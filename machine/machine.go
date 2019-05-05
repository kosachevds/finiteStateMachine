package machine

import (
	"fmt"
	"sort"
)

const (
	comma      = ','
	equalsSign = '='
)

type finiteStateMachine struct {
	isDeterministic bool
	// TODO: try without struct machineState
	states *machineState
}

func ReadFromFile(filename string) (*finiteStateMachine, error) {
	rules, err := readRules(filename)
	if err != nil {
		return nil, err
	}
	oldLen := len(rules)
	rules = determinateRules(rules)
	states, err := createStatesGraph(rules)
	if err != nil {
		return nil, fmt.Errorf("machine creation error %e", err)
	}
	return &finiteStateMachine{
		isDeterministic: oldLen == len(rules),
		states:          states,
	}, nil

}

func (fsm *finiteStateMachine) IsCanHandle(input string) bool {
	current := fsm.states
	lastIndex := len(input) - 1
	for i, item := range input {
		nextState, ok := current.ways[item]
		if !ok {
			break
		}
		if i == lastIndex && !nextState.isFinal() {
			return true
		}
		current = nextState
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func determinateRules(rules []transitionRule) []transitionRule {
	badRules := make([]transitionRule, 0, len(rules))
	otherRules := make([]transitionRule, 0, len(rules))
	for {
		badRulesIndices := findBadRules(rules)
		if badRulesIndices == nil || len(badRulesIndices) == 0 {
			break
		}
		badRules = badRules[:0]
		otherRules = otherRules[:0]
		j := 0
		for i := 0; i < len(rules); i++ {
			if j >= len(badRulesIndices) || i != badRulesIndices[j] {
				otherRules = append(otherRules, rules[i])
			} else {
				badRules = append(badRules, rules[i])
				j++
			}
		}
		rules = determinateBadRules(badRules, otherRules)
	}
	return rules
}

func findBadRules(rules []transitionRule) []int {
	ruleIndices := make([]int, 0, len(rules))
	for refIndex := 0; refIndex < len(rules)-1; refIndex++ {
		refRule := rules[refIndex]
		beginIndex := refIndex + 1
		for i, rule := range rules[beginIndex:] {
			if rule.beginState == refRule.beginState && rule.symbol == refRule.symbol {
				ruleIndices = append(ruleIndices, i+beginIndex)
			}
		}
		if len(ruleIndices) > 0 {
			ruleIndices = append(ruleIndices, refIndex)
			sort.Ints(ruleIndices)
			return ruleIndices
		}
	}
	return nil
}

// TODO badRules as struct {begin, symbol, ends}
func determinateBadRules(badRules, otherRules []transitionRule) []transitionRule {
	oldStates := getEndStates(badRules)
	sort.Ints(oldStates)
	newRules := make([]transitionRule, 0, len(badRules)+len(otherRules))
	newRules = append(newRules, uniteBadRules(badRules))
	newState := newRules[0].nextState
	for _, rule := range otherRules {
		if !isBadRuleNextState(rule.beginState, badRules) {
			newRules = append(newRules, rule)
		} else {
			newRules = append(newRules, transitionRule{
				beginState: newState,
				symbol:     rule.symbol,
				nextState:  rule.nextState,
			})
		}
	}
	return newRules
}

func isBadRuleNextState(state int, badRules []transitionRule) bool {
	for _, badRule := range badRules {
		if badRule.nextState == state {
			return true
		}
	}
	return false
}

func uniteBadRules(rules []transitionRule) transitionRule {
	// TODO: remade with count of states
	statesToUnite := getEndStates(rules)
	sort.Ints(statesToUnite)
	newState := 0
	for _, oldState := range statesToUnite {
		newState *= 10
		newState += oldState
	}
	return transitionRule{rules[0].beginState, rules[0].symbol, newState}
}

func getEndStates(rules []transitionRule) []int {
	states := make([]int, 0, cap(rules))
	for _, rule := range rules {
		states = appendIntWithoutRepeats(states, rule.nextState)
	}
	return states
}

func appendIntWithoutRepeats(buffer []int, newItem int) []int {
	for presentItem := range buffer {
		if presentItem == newItem {
			return buffer
		}
	}
	return append(buffer, newItem)
}
