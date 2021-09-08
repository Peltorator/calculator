package eval

import (
    "errors"
    "fmt"
    "strconv"
)

type ExpressionEvaluator interface {
	Evaluate(expr string) (string, error)
}

type IncrementalFakeEvaluator struct {
	Current int
}

func (ev *IncrementalFakeEvaluator) Evaluate(expr string) (string, error) {
	result := strconv.Itoa(ev.Current)
	ev.Current++
	return result, nil
}

type SmartEvaluator struct {}

var parseError = errors.New("Parse error")

func isOperator(ch byte) bool {
    return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func isDigit(ch byte) bool {
    return ('0' <= ch && ch <= '9') || ch == '.'
}

func (ev *SmartEvaluator) Evaluate(expr string) (string, error) {
    if len(expr) == 0 || expr == "\n" {
        return "", errors.New("Empty string")
    }
    if expr[len(expr) - 1] == '\n' {
        expr = expr[:len(expr) - 1]
    }
    balance := 0
    for i := 0; i < len(expr); i++ {
        ch := expr[i]
        if !isDigit(ch) && ch != '(' && ch != ')' && !isOperator(ch) {
            return "", parseError
        }
        if i > 0 && isOperator(ch) && isOperator(expr[i - 1]) {
            return "", parseError
        }
        if ch == '(' {
            if i > 0 && isDigit(expr[i - 1]) {
                return "", parseError
            }
            balance++
        }
        if ch == ')' {
            if i + 1 < len(expr) && (isDigit(expr[i + 1]) || expr[i + 1] == '(') {
                return "", parseError
            }
            balance--
            if balance < 0 {
                return "", parseError
            }
        }
    }
    if balance != 0 {
        return "", parseError
    }
    result, err := evaluateHelper(&expr, 0, len(expr))
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%.10g", result), nil
}

func sliceToInt(expr *string, leftBound int, rightBound int) (float64, error) {
    var strRep string = ""
    for i := leftBound; i < rightBound; i++ {
        strRep += string((*expr)[i])
        }
    result, err := strconv.ParseFloat(strRep, 64)
    if err != nil {
        return 0.0, parseError
    }
    return result, nil
}

func evaluateOperator(leftValue float64, op byte, rightValue float64) (float64, error) {
    if op == '+' {
        return leftValue + rightValue, nil
    } else if op == '-' {
        return leftValue - rightValue, nil
    } else if op == '*' {
        return leftValue * rightValue, nil
    } else if op == '/' {
        if rightValue == 0 {
            return 0.0, errors.New("Result is ambiguous")
        }
        return leftValue / rightValue, nil
    } else {
        return 0.0, parseError
    }
}

func basicEvaluation(vals []float64, ops []byte) (float64, error) {
    if len(ops) == 0 {
        return vals[0], nil
    }
    var ans float64
    for i := 0; i < len(vals); i++ {
        curVal := vals[i]
        fin := -1
        for j := i + 1; j < len(vals); j++ {
            if ops[j - 1] == '+' || ops[j - 1] == '-' {
                fin = j - 1
                break
            }
            result, err := evaluateOperator(curVal, ops[j - 1], vals[j])
            if err != nil {
                return 0.0, err
            }
            curVal = result
        }
        if i == 0 {
            ans = curVal
        } else {
            result, err := evaluateOperator(ans, ops[i - 1], curVal)
            if err != nil {
                return 0.0, err
            }
            ans = result
        }
        if fin == -1 {
            break
        }
        i = fin
    }
    return ans, nil
}

func evaluateHelper(expr *string, leftBound int, rightBound int) (float64, error) {
    if leftBound >= rightBound {
        return 0.0, parseError
    }
    if (isOperator((*expr)[0]) && (*expr)[0] != '-') || isOperator((*expr)[rightBound - 1]) {
        return 0.0, parseError
    }
    var ops []byte
    var vals []float64
    for i := leftBound; i < rightBound; i++ {
        j := i
        if (*expr)[i] == '-' {
            j++
        }
        var val float64
        midPoint := -1
        if (*expr)[j] == '(' {
            balance := 1
            for k := j + 1; k < rightBound; k++ {
                if (*expr)[k] == '(' {
                    balance++
                }
                if (*expr)[k] == ')' {
                    balance--
                    if balance == 0 {
                        midPoint = k
                        break
                    }
                }
            }
            if midPoint == -1 {
                return 0.0, parseError
            }
            if midPoint != rightBound - 1 && !isOperator((*expr)[midPoint + 1]) {
                return 0.0, parseError
            }
            curval, err := evaluateHelper(expr, j + 1, midPoint)
            if err != nil {
                return 0.0, err
            }
            val = curval
            midPoint++
        } else {
            for k := j; k < rightBound; k++ {
                if !isDigit((*expr)[k]) {
                    midPoint = k
                    break
                }
            }
            if midPoint == -1 {
                midPoint = rightBound
            }
            if midPoint != rightBound && !isOperator((*expr)[midPoint]) {
                return 0.0, parseError
            }
            var newVal float64
            newVal, err := sliceToInt(expr, j, midPoint)
            if err != nil {
                return 0.0, parseError
            }
            val = newVal
        }
        if (*expr)[i] == '-' {
            val = -val
        }
        vals = append(vals, val)
        if midPoint == rightBound {
            break
        }
        if !isOperator((*expr)[midPoint]) {
            return 0.0, parseError
        }
        ops = append(ops, (*expr)[midPoint])
        i = midPoint
    }
    if len(ops) + 1 != len(vals) {
        return 0.0, parseError
    }
    return basicEvaluation(vals, ops)
}
