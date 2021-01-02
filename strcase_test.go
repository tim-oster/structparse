package structparse

import (
	"testing"
)

func TestJoinWithCasing(t *testing.T) {
	tests := []struct {
		casing   StrCase
		input    []string
		expected string
	}{
		{casing: StrCaseCamel, input: []string{"ONE", "two"}, expected: "oneTwo"},
		{casing: StrCasePascal, input: []string{"ONE", "two"}, expected: "OneTwo"},
		{casing: StrCaseUpperSnake, input: []string{"ONE", "two"}, expected: "ONE_TWO"},
		{casing: StrCaseLowerSnake, input: []string{"ONE", "two"}, expected: "one_two"},
		{casing: StrCaseUpperKebab, input: []string{"ONE", "two"}, expected: "ONE-TWO"},
		{casing: StrCaseLowerKebab, input: []string{"ONE", "two"}, expected: "one-two"},
	}
	for _, test := range tests {
		t.Run(casingNames[test.casing], func(t *testing.T) {
			result := JoinWithCasing(test.casing, test.input)
			assertEqual(t, test.expected, result)
		})
	}
}

func TestSplitWithCasing(t *testing.T) {
	tests := []struct {
		casing   StrCase
		input    string
		expected []string
	}{
		// StrCaseCamel and StrCasePascal are tested separately
		{casing: StrCaseUpperSnake, input: "ONE_TWO", expected: []string{"ONE", "TWO"}},
		{casing: StrCaseLowerSnake, input: "one_two", expected: []string{"one", "two"}},
		{casing: StrCaseUpperKebab, input: "ONE-TWO", expected: []string{"ONE", "TWO"}},
		{casing: StrCaseLowerKebab, input: "one-two", expected: []string{"one", "two"}},
	}
	for _, test := range tests {
		t.Run(casingNames[test.casing], func(t *testing.T) {
			result := SplitWithCasing(test.casing, test.input)
			assertEqual(t, test.expected, result)
		})
	}
}

func TestSplitWithCasing_camel(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "lowercase", expected: []string{"lowercase"}},
		{input: "UPPERCASE", expected: []string{"UPPERCASE"}},
		{input: "camelCase", expected: []string{"camel", "Case"}},
		{input: "PascalCase", expected: []string{"Pascal", "Case"}},
		{input: "UPPERStart", expected: []string{"UPPER", "Start"}},
		{input: "upperEND", expected: []string{"upper", "END"}},
		{input: "camelWITHUpper", expected: []string{"camel", "WITH", "Upper"}},
		{input: "PascalWITHUpper", expected: []string{"Pascal", "WITH", "Upper"}},
		{input: "ThisIsATest", expected: []string{"This", "Is", "A", "Test"}},
	}
	for _, test := range tests {
		result := SplitWithCasing(StrCaseCamel, test.input)
		assertEqual(t, test.expected, result)
		result = SplitWithCasing(StrCasePascal, test.input)
		assertEqual(t, test.expected, result)
	}
}

func TestSwitchCasing(t *testing.T) {
	result := SwitchCasing(StrCaseCamel, StrCaseLowerKebab, "ThisIsATest")
	assertEqual(t, "this-is-a-test", result)
}

func TestXxxToXxx(t *testing.T) {
	tests := []struct {
		fn       func(string) string
		input    string
		expected string
	}{
		{fn: CamelToUpperSnake, input: "aB", expected: "A_B"},
		{fn: CamelToLowerSnake, input: "aB", expected: "a_b"},
		{fn: CamelToUpperKebab, input: "aB", expected: "A-B"},
		{fn: CamelToLowerKebab, input: "aB", expected: "a-b"},
	}
	for _, test := range tests {
		assertEqual(t, test.expected, test.fn(test.input))
	}
}
