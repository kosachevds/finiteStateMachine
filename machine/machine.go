package machine

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func determinateRules(rules []transitionRule) []transitionRule {
	badRules := make([]transitionRule, 0, len(rules))
	otherRules := make([]transitionRule, 0, len(rules))
	newStates := make(map[string]bool)
	for {
		badRulesIndices := findBadRules(rules)
		if badRulesIndices == nil || len(badRulesIndices) == 0 {
			break
		}
		badRules = badRules[:0]
		otherRules = otherRules[:0]
		j := 0
		for i := range rules {
			if j >= len(badRulesIndices) || i != badRulesIndices[j] {
				otherRules = append(otherRules, rules[i])
			} else {
				badRules = append(badRules, rules[i])
				j++
			}
		}
		rules = determinateBadRules(badRules, otherRules, newStates)
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
func determinateBadRules(badRules, otherRules []transitionRule, newStates map[string]bool) []transitionRule {
	unitedRule := uniteBadRules(badRules)
	newStateName := uniteEndStateNames(badRules)
	newState := unitedRule.nextState
	newRules := otherRules
	newRules = append(newRules, unitedRule)
	if _, ok := newStates[newStateName]; ok {
		return newRules
	}
	newStates[newStateName] = true
	for _, rule := range otherRules {
		if isBadRuleNextState(rule.beginState, badRules) {
			newRules = append(newRules, transitionRule{
				beginState:   newState,
				symbol:       rule.symbol,
				nextState:    rule.nextState,
				toFinalState: rule.toFinalState,
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

func uniteEndStateNames(rules []transitionRule) string {
	endStates := getEndStates(rules)
	sort.Ints(endStates)
	unitedName := fmt.Sprint(endStates)
	unitedName = strings.Join(strings.Fields(unitedName), "")
	unitedName = strings.Trim(unitedName, "[]")
	return unitedName
}

func uniteBadRules(rules []transitionRule) transitionRule {
	// TODO: remade with count of root (but it will fail new root check)
	statesToUnite := getEndStates(rules)
	sort.Ints(statesToUnite)
	newState := 0
	for _, oldState := range statesToUnite {
		newState *= 10
		newState += oldState
	}

	return transitionRule{
		beginState:   rules[0].beginState,
		symbol:       rules[0].symbol,
		nextState:    newState,
		toFinalState: containsFinalRule(rules),
	}
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

func containsFinalRule(rules []transitionRule) bool {
	for _, r := range rules {
		if r.toFinalState {
			return true
		}
	}
	return false
}

func countStates(rules []transitionRule) int {
	codes := make(map[int]bool)
	for _, rule := range rules {
		codes[rule.beginState] = true
	}
	return len(codes)
}
