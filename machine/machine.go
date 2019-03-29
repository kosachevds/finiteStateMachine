package machine

const (
	comma          = ','
	equalsSign     = '='
	finalSign      = '#'
	startStateName = "q0"
)

type finiteStateMachine struct {
	isDeterministic bool
	startState      machineState
}

func ReadFromFile(filename string) (*finiteStateMachine, error) {
	return &finiteStateMachine{}, nil
}

func (fsm *finiteStateMachine) IsDeterministic() bool {
	return fsm.isDeterministic
}

func (fsm *finiteStateMachine) IsCanHandle(input string) bool {
	return true
}

func determinateRules(rules []transitionRule) []transitionRule {
	badRules := make([]transitionRule, 0, len(rules))
	newRules := make([]transitionRule, 0, len(rules))
	rulesBuffer := make([]transitionRule, 0, len(rules))
	for {
		badRulesIndices := selectBadRules(rules)
		if badRulesIndices == nil || len(badRulesIndices) == 0 {
			break
		}
		badRules = badRules[:0]
		rulesBuffer = rulesBuffer[:0]
		for i, j := 0, 0; i < len(rules); i++ {
			if i != badRulesIndices[j] {
				rulesBuffer = append(rulesBuffer, rules[i])
			} else {
				badRules = append(badRules, rules[i])
				j++
			}
		}
		rules = rulesBuffer

	}
}

func selectBadRules(rules []transitionRule) []int {
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

func createDeterministicRules(rules []transitionRule) []transitionRule {

}
