package eval

type ExpressionEvaluator interface {
	Evaluate(expression string) (string, error)
}

