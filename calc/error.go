package calc

import "errors"

var (
	createExpressError   = errors.New("不合法的表达式")
	illegalBracketError  = errors.New("不合法的括号")
	illegalNegativeError = errors.New("负数定义错误")
	illegalDecimalError  = errors.New("定义小数错误")
	illegalNumberError   = errors.New("定义数字错误")
	illegalSignalError   = errors.New("不合法的符号")
	illegalCompareError  = errors.New("不合法的比较表达式")
)
