package machine

import (
	"fmt"
	"sort"
	"strings"
)

type unitedState string

type unitedStatesCodes struct {
	newCode     int
	existsCodes map[unitedState]int
}

func newUnitedStatesCodes(firstCode int) *unitedStatesCodes {
	return &unitedStatesCodes{
		newCode:     firstCode,
		existsCodes: make(map[unitedState]int),
	}
}

func (codes *unitedStatesCodes) get(state unitedState) (int, bool) {
	item, ok := codes.existsCodes[state]
	return item, ok
}

func (codes *unitedStatesCodes) isContains(state unitedState) bool {
	_, ok := codes.existsCodes[state]
	return ok
}

func (codes *unitedStatesCodes) add(state unitedState) {
	if _, ok := codes.get(state); ok {
		return
	}
	codes.existsCodes[state] = codes.newCode
	codes.newCode++
}

func determinateRules(rules []transitionRule) []transitionRule {
	maxCode := maxStateCode(rules)
	badRules := make([]transitionRule, 0, len(rules))
	otherRules := make([]transitionRule, 0, len(rules))
	newStatesCodes := newUnitedStatesCodes(maxCode + 1)
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
		maxCode++
		rules = determinateBadRules(badRules, otherRules, newStatesCodes)
	}
	return rules
}

func maxStateCode(rules []transitionRule) int {
	maxCode := 0
	for _, rule := range rules {
		if rule.nextState > maxCode {
			maxCode = rule.nextState
		}
	}
	return maxCode
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
func determinateBadRules(badRules, otherRules []transitionRule, codes *unitedStatesCodes) []transitionRule {
	unitedRule := uniteBadRules(badRules)
	newStateName := uniteEndStates(badRules)
	newRules := otherRules
	newRules = append(newRules, unitedRule)
	if codes.isContains(newStateName) {
		return newRules
	}
	codes.add(newStateName)
	for _, rule := range otherRules {
		if isBadRuleNextState(rule.beginState, badRules) {
			newRules = append(newRules, transitionRule{
				beginState:   unitedRule.nextState,
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

func uniteEndStates(rules []transitionRule) unitedState {
	endStates := getEndStates(rules)
	sort.Ints(endStates)
	unitedName := fmt.Sprint(endStates)
	unitedName = strings.Join(strings.Fields(unitedName), "")
	unitedName = strings.Trim(unitedName, "[]")
	return unitedState(unitedName)
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
