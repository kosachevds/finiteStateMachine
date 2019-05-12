package machine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const commentMark = ";"
const finalStatePrefix = "f"

type stateCodesMap struct {
	codes map[string]int
}

type transitionRule struct {
	// TODO: remade with field "isFinal" for uniting states
	beginState   int
	symbol       rune
	nextState    int
	toFinalState bool
}

func (sc *stateCodesMap) getCode(stateName string) int {
	code, ok := sc.codes[stateName]
	if !ok {
		code = len(sc.codes) + 1
		sc.codes[stateName] = code
	}
	return code
}

func newStateCodesMap() *stateCodesMap {
	return &stateCodesMap{make(map[string]int)}
}

func parseTransitionRule(rule string, stateCodes *stateCodesMap) (transitionRule, error) {
	commaIndex := strings.IndexRune(rule, comma)
	equalsSignIndex := strings.IndexRune(rule, equalsSign)
	if commaIndex == -1 || equalsSignIndex == -1 {
		return transitionRule{}, fmt.Errorf("the string %s is not valid rule", rule)
	}

	beginStateName := rule[:commaIndex]
	beginStateCode := stateCodes.getCode(beginStateName)

	nextStateName := rule[equalsSignIndex+1:]
	nextStateCode := stateCodes.getCode(nextStateName)

	symbol, _ := utf8.DecodeRuneInString(rule[commaIndex+1 : commaIndex+2])

	return transitionRule{
		beginState:   beginStateCode,
		nextState:    nextStateCode,
		symbol:       symbol,
		toFinalState: strings.HasPrefix(nextStateName, finalStatePrefix),
	}, nil
}

func readRules(filename string) ([]transitionRule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make(chan string)
	go fillRuleLinesChannel(scanner, lines, func(line string) bool {
		return line != "" && !strings.HasPrefix(line, commentMark)
	})
	rules, err := createRules(lines)
	if err != nil {
		return nil, err
	}
	if err = scanner.Err(); err != nil {
		return nil, scanner.Err()
	}
	return rules, nil
}

func fillRuleLinesChannel(scanner *bufio.Scanner, lines chan string, predicate func(string) bool) {
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if predicate != nil && !predicate(line) {
			continue
		}
		lines <- line
	}
	close(lines)
}

func createRules(lines chan string) ([]transitionRule, error) {
	stateCodes := newStateCodesMap()
	var rules []transitionRule
	for line := range lines {
		newRule, err := parseTransitionRule(line, stateCodes)
		if err != nil {
			return nil, err
		}
		rules = append(rules, newRule)
	}
	return rules, nil
}
