package machine

const (
	comma      = ','
	equalsSign = '='
	finalSign  = '#'
)

type finiteStateMachine struct {
	isDeterministic bool
	startState      machineState
}

func ReadFromFile(filename string) (*finiteStateMachine, error) {

}

func (fsm *finiteStateMachine) IsDeterministic() bool {
	return fsm.isDeterministic
}

func (fsm *finiteStateMachine) IsCanHandle(input string) bool {

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
		for i, j := 0, 0; i < len(rules); i++ {
			if i != badRulesIndices[j] {
				otherRules = append(otherRules, rules[i])
			} else {
				badRules = append(badRules, rules[i])
				j++
			}
		}
		rules = otherRules
		otherRules = determinateBadRules(badRules, otherRules)
		rules = append(rules, otherRules...)
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
			return append(ruleIndices, refIndex)
		}
	}
	return nil
}

func determinateBadRules(badRules, otherRules []transitionRule) []transitionRule {

}

func uniteStates(states []int) int {
	// TODO: remade with count of states
	newState := 0
	for _, state := range states {
		newState *= 10
		newState += state
	}
	return newState
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
