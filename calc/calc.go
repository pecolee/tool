package calc

import (
	"fmt"
	"math"
	"strconv"
)

type token struct {
	kind  rune
	value float64
}

// Evaluate Arithmetic return result
func Evaluate(input string) (float64, error) {
	tokens, err := lex(input)
	if err != nil {
		return 0, err
	}
	result, restTokens, err := parseExpression(tokens)
	if err != nil {
		return 0, err
	}
	if len(restTokens) > 0 {
		return 0, fmt.Errorf("含了未知的 token 或错误的字符")
	}
	return result, nil
}

func lex(input string) ([]token, error) {
	var tokens []token
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '+':
			tokens = append(tokens, token{kind: '+', value: 0})
		case '-':
			tokens = append(tokens, token{kind: '-', value: 0})
		case '*':
			tokens = append(tokens, token{kind: '*', value: 0})
		case '/':
			tokens = append(tokens, token{kind: '/', value: 0})
		case '(':
			tokens = append(tokens, token{kind: '(', value: 0})
		case ')':
			tokens = append(tokens, token{kind: ')', value: 0})
		default:
			j := i
			for ; j < len(input); j++ {
				if input[j] == '+' || input[j] == '-' || input[j] == '*' || input[j] == '/' || input[j] == '(' || input[j] == ')' {
					break
				}
			}
			value, err := strconv.ParseFloat(input[i:j], 64)
			if err != nil {
				return nil, fmt.Errorf("无效输入: %v", input[i:j])
			}
			tokens = append(tokens, token{kind: 'n', value: value})
			i = j - 1
		}
	}
	return tokens, nil
}

func parseFactor(tokens []token) (float64, []token, error) {
	if tokens[0].kind == '(' {
		result, restTokens, err := parseExpression(tokens[1:])
		if err != nil {
			return 0, nil, err
		}
		if restTokens[0].kind != ')' {
			return 0, nil, fmt.Errorf("缺少 )")
		}
		return result, restTokens[1:], nil
	}
	if tokens[0].kind == 45 {
		return tokens[0].value, tokens, nil
	}
	if tokens[0].kind == 'n' {
		return tokens[0].value, tokens[1:], nil
	}
	return 0, nil, fmt.Errorf("无法识别符号: %v", tokens[0].kind)
}

func parseTerm(tokens []token) (float64, []token, error) {
	factor1, restTokens, err := parseFactor(tokens)
	if err != nil {
		return 0, nil, err
	}
	for len(restTokens) > 0 && (restTokens[0].kind == '*' || restTokens[0].kind == '/') {
		factor2, restTokens2, err := parseFactor(restTokens[1:])
		if err != nil {
			return 0, nil, err
		}
		if restTokens[0].kind == '*' {
			factor1 *= factor2
		} else {
			if factor2 == 0 {
				return 0, nil, fmt.Errorf(" 除数不能为0")
			}
			factor1 /= factor2
		}
		restTokens = restTokens2
	}
	return factor1, restTokens, nil
}

func parseExpression(tokens []token) (float64, []token, error) {
	term1, restTokens, err := parseTerm(tokens)
	if err != nil {
		return 0, nil, err
	}
	result := term1
	for len(restTokens) > 0 && (restTokens[0].kind == '+' || restTokens[0].kind == '-') {
		term2, restTokens2, err := parseTerm(restTokens[1:])
		if err != nil {
			return 0, nil, err
		}
		if restTokens[0].kind == '+' {
			result += term2
		} else {
			result -= term2
		}
		restTokens = restTokens2
	}
	return result, restTokens, nil
}

// FormatFloat Solve the rounding defect of go language (decimal also can do this)
func FormatFloat(f float64, precision int) float64 {
	if precision == 0 {
		return math.Round(f)
	}
	p := math.Pow10(precision)
	if precision < 0 {
		return math.Round(f*p) * math.Pow10(-precision)
	}
	return math.Round(f*p) / p
}
