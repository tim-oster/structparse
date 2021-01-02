package structparse

import (
	"strings"
	"unicode"
)

type StrCase byte

const (
	// StrCaseCamel defines camel case with support for all caps words.
	// When converting from this casing, StrCasePascal is also accepted.
	// Supported: camelCase, camelWITHUpper
	StrCaseCamel StrCase = iota + 1

	// StrCasePascal defines pascal case with support for all caps words.
	// When converting from this casing, StrCaseCamel is also accepted.
	// Supported: PascalCase, PascalWITHUpper
	StrCasePascal

	// StrCaseUpperSnake defines snake casing with only uppercase chars.
	// Supported: SNAKE_CASE
	StrCaseUpperSnake

	// StrCaseLowerSnake defines snake casing with only lowercase chars.
	// Supported: snake_case
	StrCaseLowerSnake

	// StrCaseUpperKebab defines kebab casing with only uppercase chars.
	// Supported: KEBAB-CASE
	StrCaseUpperKebab

	// StrCaseLowerKebab defines kebab casing with only lowercase chars.
	// Supported: kebab-case
	StrCaseLowerKebab
)

var (
	casingNames = map[StrCase]string{
		StrCaseCamel:      "camel_case",
		StrCasePascal:     "pascal_case",
		StrCaseUpperSnake: "upper_snake_case",
		StrCaseLowerSnake: "lower_snake_case",
		StrCaseUpperKebab: "upper_kebab_case",
		StrCaseLowerKebab: "lower_kebab_case",
	}
	casingSeparator = map[StrCase]string{
		StrCaseUpperSnake: "_",
		StrCaseLowerSnake: "_",
		StrCaseUpperKebab: "-",
		StrCaseLowerKebab: "-",
	}
)

func JoinWithCasing(casing StrCase, words []string) string {
	switch casing {
	case StrCaseCamel, StrCasePascal:
		for i, word := range words {
			word = strings.ToLower(word)
			if casing == StrCasePascal || i > 0 {
				word = strings.Title(word)
			}
			words[i] = word
		}

	case StrCaseUpperSnake, StrCaseUpperKebab:
		for i, word := range words {
			words[i] = strings.Map(unicode.ToUpper, word)
		}

	case StrCaseLowerSnake, StrCaseLowerKebab:
		for i, word := range words {
			words[i] = strings.Map(unicode.ToLower, word)
		}

	default:
		panic("invalid casing")
	}

	return strings.Join(words, casingSeparator[casing])
}

func SplitWithCasing(casing StrCase, s string) []string {
	switch casing {
	case StrCaseCamel, StrCasePascal:
		var words []string

		for start := 0; start < len(s); {
			wordStart := unicode.IsUpper(rune(s[start]))

			// look how long the next camel-cased word is
			var (
				wordLen        int
				expectedCasing bool
			)
			for wordLen = 1; start+wordLen < len(s); wordLen++ {
				casing := unicode.IsUpper(rune(s[start+wordLen]))

				if wordLen == 1 {
					if wordStart {
						if casing {
							// supports camelWITHUpper and PascalWITHUpper
							expectedCasing = true
						} else {
							// supports PascalCase
							expectedCasing = false
						}
					} else {
						// supports camelCase
						expectedCasing = false
					}
				}

				if casing != expectedCasing {
					// complete uppercase words stops can only be detected on the start of the
					// following word (see IPAddress, where 'A' determines the end of IP but the
					// change in casing happens between 'A' and 'd') so the word length has to
					// be adjusted
					if expectedCasing {
						wordLen--
					}
					break
				}
			}

			words = append(words, s[start:start+wordLen])
			start += wordLen
		}

		return words

	case StrCaseUpperSnake, StrCaseLowerSnake, StrCaseUpperKebab, StrCaseLowerKebab:
		return strings.Split(s, casingSeparator[casing])

	default:
		panic("invalid casing")
	}
}

func SwitchCasing(from, to StrCase, s string) string {
	words := SplitWithCasing(from, s)
	return JoinWithCasing(to, words)
}

func CamelToUpperSnake(s string) string {
	return SwitchCasing(StrCaseCamel, StrCaseUpperSnake, s)
}

func CamelToLowerSnake(s string) string {
	return SwitchCasing(StrCaseCamel, StrCaseLowerSnake, s)
}

func CamelToUpperKebab(s string) string {
	return SwitchCasing(StrCaseCamel, StrCaseUpperKebab, s)
}

func CamelToLowerKebab(s string) string {
	return SwitchCasing(StrCaseCamel, StrCaseLowerKebab, s)
}
