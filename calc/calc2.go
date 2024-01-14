package calc

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/shopspring/decimal"
	"sort"
	"strconv"
	"strings"
)

// CalcValue 使用自研的计算引擎(比较结果,差额)
func CalcValue(schema string, paramMap map[string]interface{}) (interface{}, float64, error) {
	paramMap = dealParam(paramMap)
	result, err := calc(schema, paramMap)
	if err != nil {
		return nil, 0, err
	}
	diffValue := 0.0
	if v, ok := result.(bool); ok == true && v == true {
		// 判断结果为true不需要算差额
	} else {
		if strings.Contains(schema, ">=") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, ">=", "-("))
		}
		if strings.Contains(schema, "<=") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, "<=", "-("))
		}
		if strings.Contains(schema, "<") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, "<", "-("))
		}
		if strings.Contains(schema, ">") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, ">", "-("))
		}
		if strings.Contains(schema, "=") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, "=", "-("))
		}
		if strings.Contains(schema, "!=") {
			schema = fmt.Sprintf("%s)", strings.ReplaceAll(schema, "!=", "-("))
		}
		diff, diffError := calc(schema, paramMap)
		if diffError != nil {
			return nil, 0, diffError
		}
		diffValue = gconv.Float64(diff)
	}
	return result, diffValue, nil
}

//参数处理原有逻辑
func dealParam(paramDict map[string]interface{}) map[string]interface{} {
	if paramDict == nil {
		return nil
	}
	for k, v := range paramDict {
		vStr := fmt.Sprintf("%s", v)
		if strings.Contains(vStr, "nil") {
			v = ""
		}
		floatValue, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if err == nil {
			delete(paramDict, k)
			paramDict[k] = fmt.Sprintf("%0.2f", floatValue)
		}
	}
	return paramDict
}

func calc(schema string, paramMap map[string]interface{}) (interface{}, error) {
	if paramMap != nil {
		// key值最长匹配
		lenMap := map[int][]string{}
		indexList := make([]int, 0)
		for key := range paramMap {
			length := len(key)
			if lenMap[length] == nil {
				lenMap[length] = make([]string, 0)
				indexList = append(indexList, length)
			}
			lenMap[length] = append(lenMap[length], key)
		}
		sort.Ints(indexList)
		for i := len(indexList) - 1; i >= 0; i-- {
			for _, key := range lenMap[indexList[i]] {
				schema = strings.ReplaceAll(schema, key, fmt.Sprintf("%v", paramMap[key]))
			}
		}
	}
	if len(schema) <= 0 {
		return schema, nil
	}
	token := ""
	if strings.Contains(schema, LTE) {
		token = LTE
	}
	if strings.Contains(schema, LT) && len(token) <= 0 {
		token = LT
	}
	if strings.Contains(schema, GTE) && len(token) <= 0 {
		token = GTE
	}
	if strings.Contains(schema, GT) && len(token) <= 0 {
		token = GT
	}
	if strings.Contains(schema, Not) && len(token) <= 0 {
		token = Not
	}
	if strings.Contains(schema, Equal) && len(token) <= 0 {
		token = Equal
	}
	if len(token) > 0 {
		schemaList := strings.Split(schema, token)
		if len(schemaList) != 2 {
			return false, illegalCompareError
		}
		v1Str, err := CalcStr(schemaList[0])
		if err != nil {
			return nil, err
		}
		v1, err := decimal.NewFromString(v1Str)
		if err != nil {
			return nil, err
		}
		v2Str, err := CalcStr(schemaList[1])
		if err != nil {
			return nil, err
		}
		v2, err := decimal.NewFromString(v2Str)
		if err != nil {
			return nil, err
		}
		return compare(v1, v2, token)
	} else {
		return CalcStr(schema)
	}
}

func CalcStr(schema string) (string, error) {
	if len(schema) <= 0 {
		return schema, nil
	}
	// 首先将本次表达式的括号的值算出来 替换位置
	isEnd := false
	for true {
		for pos, value := range schema {
			isEnd = len(schema)-1 == pos
			if value == LeftBracket {
				// 如果是左括号 那么需要找到该位置下一个右括号
				newStr := schema[pos:]
				rightBracketPos := findRightBracket(newStr)
				if rightBracketPos < 0 {
					return "", illegalBracketError
				}
				innerValue, err := CalcStr(newStr[1:rightBracketPos])
				if err != nil {
					return "", err
				}
				schema = fmt.Sprintf("%s%s%s", schema[0:pos], innerValue, newStr[rightBracketPos+1:])
				break
			} else if value == RightBracket {
				// 不可能仅仅存在右括号的情况
				return "", illegalBracketError
			}
		}
		if isEnd {
			break
		}
	}

	// 扁平化计算
	// 下一个肯定是数字
	numberList := make([]string, 0)
	tokenList := make([]string, 0)
	nextMustBeNumber := true
	isNeg := false
	isDecimal := false
	str := ""
	for _, char := range schema {
		if char == Space {
			continue
		}
		if char == Mut || char == Add || char == Sub || char == Div {
			if nextMustBeNumber == false && len(str) > 0 {
				if len(str) > 0 {
					numberList = append(numberList, str)
				}
				tokenList = append(tokenList, fmt.Sprintf("%c", char))
				isNeg = false
				isDecimal = false
				nextMustBeNumber = true
				str = ""
				continue
			} else if nextMustBeNumber == false {
				return "", illegalSignalError
			}
		}
		if char == Neg {
			if char == Neg && isNeg == false && nextMustBeNumber {
				nextMustBeNumber = false
				isNeg = true
				str = fmt.Sprintf("%s%c", str, char)
			} else {
				return "", illegalNegativeError
			}
		}
		if char == Dot {
			if isDecimal == false {
				isDecimal = true
				nextMustBeNumber = false
				str = fmt.Sprintf("%s%c", str, char)
			} else {
				return "", illegalDecimalError
			}
		}
		if N0 <= char && N9 >= char {
			nextMustBeNumber = false
			str = fmt.Sprintf("%s%c", str, char)
		}
	}
	if nextMustBeNumber == true || len(str) == 0 {
		return "", illegalSignalError
	}
	// 先计算乘除 再计算 +-
	numberList = append(numberList, str)
	for index, token := range tokenList {
		if token == AddStr || token == SubStr {
			continue
		}
		number1, err := decimal.NewFromString(numberList[index])
		if err != nil {
			return "", err
		}
		number2, err := decimal.NewFromString(numberList[index+1])
		if err != nil {
			return "", err
		}
		value, err := switchDecimal(number1, number2, token)
		if err != nil {
			return "", err
		}
		numberList[index] = "0"
		numberList[index+1] = value.String()
		tokenList[index] = AddStr
	}

	// 计算加 减
	for index, token := range tokenList {
		number1, err := decimal.NewFromString(numberList[index])
		if err != nil {
			return "", err
		}
		number2, err := decimal.NewFromString(numberList[index+1])
		if err != nil {
			return "", err
		}
		value, err := switchDecimal(number1, number2, token)
		if err != nil {
			return "", err
		}
		numberList[index+1] = value.String()
		tokenList[index] = AddStr
	}
	return numberList[len(numberList)-1], nil
}

func findRightBracket(schema string) int {
	index := 0
	start := false
	for pos, v := range schema {
		if v == LeftBracket {
			index++
			start = true
		}
		if v == RightBracket {
			index--
		}
		if index == 0 && start {
			return pos
		}
	}
	return -1
}

func switchDecimal(d1, d2 decimal.Decimal, token string) (decimal.Decimal, error) {
	switch token {
	case MutStr:
		return d1.Mul(d2), nil
	case AddStr:
		return d1.Add(d2), nil
	case SubStr:
		return d1.Sub(d2), nil
	case DivStr:
		if d2.IsZero() {
			return d1, errors.New("除数不可为零")
		}
		return d1.Div(d2), nil
	}
	return decimal.NewFromFloat(0), nil
}

func compare(v1, v2 decimal.Decimal, signal string) (bool, error) {
	switch signal {
	case LT:
		return v1.LessThan(v2), nil
	case LTE:
		return v1.LessThanOrEqual(v2), nil
	case GT:
		return v1.GreaterThan(v2), nil
	case GTE:
		return v1.GreaterThanOrEqual(v2), nil
	case Not:
		return !v1.Equal(v2), nil
	case Equal:
		return v1.Equal(v2), nil
	default:
		return false, illegalCompareError
	}
}
