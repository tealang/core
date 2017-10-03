package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

func GenerateFunctionCall(input []tokens.Token) (nodes.Node, int, error) {
	return nil, 0, nil
}

type EmptyOperatorStackException struct{}

func (EmptyOperatorStackException) Error() string {
	return "EmptyOperatorStackException: The operator stack is empty"
}

type MissingOperatorArgsException struct {
}

func (MissingOperatorArgsException) Error() string {
	return "MissingOperatorArgsException: The operation is missing arguments"
}

type UnexpectedNodeTypeException struct{}

func (UnexpectedNodeTypeException) Error() string {
	return "UnexpectedNodeTypeException: Unexpected node type, cannot be handled correctly"
}

type InvalidExpressionException struct{}

func (InvalidExpressionException) Error() string {
	return "InvalidExpressionException: The expression is invalid"
}

func GenerateTerm(input []tokens.Token) (nodes.Node, int, error) {
	operands, operators := make([]nodes.Node, 0), make([]*nodes.Operation, 0)

	popOperator := func() error {
		if len(operators) < 1 {
			return EmptyOperatorStackException{}
		}
		var op *nodes.Operation
		op, operators = operators[len(operators)-1], operators[:len(operators)-1]
		if len(operands) < op.ArgCount {
			return MissingOperatorArgsException{}
		}
		for i := 0; i < op.ArgCount; i++ {
			var value nodes.Node
			value, operands = operands[len(operands)-1], operands[:len(operands)-1]
			op.AddFront(value)
		}
		operands = append(operands, op)
		return nil
	}

	popOperand := func() nodes.Node {
		value := operands[len(operands)-1]
		operands = operands[:len(operands)-1]
		return value
	}

	var (
		index       int
		level       int
		activeToken tokens.Token
	)

	for index = 0; index < len(input); index++ {
		//lastToken = activeToken
		activeToken = input[index]

		switch activeToken.Type {
		case tokens.Separator:
			for len(operators) > 0 {
				err := popOperator()
				if err != nil {
					return nil, 0, err
				}
			}
			if level > 0 {
				return operands[0], index, nil
			}
		}
	}

	if len(operands) != 1 {
		return nil, 0, InvalidExpressionException{}
	}

	return operands[0], index, nil
}
