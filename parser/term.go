package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/types"

	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

func newTermParser() *termParser {
	tp := &termParser{
		output:    newItemStack(),
		operators: newItemStack(),
	}
	tp.handlers = map[*tokens.Type]func() error{
		tokens.Identifier:       tp.handleIdentifier,
		tokens.String:           tp.handleString,
		tokens.Number:           tp.handleNumber,
		tokens.Operator:         tp.handleOperator,
		tokens.LeftParentheses:  tp.handleLeftParentheses,
		tokens.RightParentheses: tp.handleRightParentheses,
	}
	return tp
}

type termParser struct {
	output, operators      *itemStack
	active, previous, next tokens.Token
	handlers               map[*tokens.Type]func() error
	index, size            int
	input                  []tokens.Token
}

func (tp *termParser) binding(item termItem) bool {
	switch item.Value.Value {
	case "^", "!":
		return true
	case "+", "-":
		if tp.isUnaryOperator(item) {
			return true
		}
	}
	return false
}

func (tp *termParser) argCount(item termItem) int {
	if tp.isUnaryOperator(item) {
		return 1
	}
	return 2
}

func (termParser) isUnaryOperator(item termItem) bool {
	switch item.Value.Value {
	case "+", "-":
		switch item.Previous.Type {
		case nil, tokens.Operator:
			return true
		}
	case "!":
		return true
	}
	return false
}

func (tp *termParser) priority(item termItem) int {
	switch item.Value.Value {
	case "&", "|":
		return 8
	case "!":
		return 7
	case "^":
		return 6
	case "*", "/":
		return 5
	case "+", "-":
		if tp.isUnaryOperator(item) {
			return 9
		}
		return 4
	case ":":
		return 4
	case "%":
		return 3
	case "<", ">", ">=", "<=", "=<", "!=", "==":
		return 2
	case "&&", "||", "^|":
		return 1
	case "=>":
		return 0
	default:
		return -1
	}
}

const (
	trueLiteral  = "true"
	falseLiteral = "false"
	nullLiteral  = "null"
)

func (tp *termParser) itemFromActive(node nodes.Node) termItem {
	return termItem{
		Value:    tp.active,
		Next:     tp.next,
		Previous: tp.previous,
		Node:     node,
	}
}

func (tp *termParser) handleIdentifier() error {
	switch tp.active.Value {
	case trueLiteral:
		tp.output.Push(tp.itemFromActive(nodes.NewLiteral(types.True)))
	case falseLiteral:
		tp.output.Push(tp.itemFromActive(nodes.NewLiteral(types.False)))
	case nullLiteral:
		tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{})))
	default:
		if tp.next.Type == tokens.LeftParentheses {
			alias := tp.active.Value
			tp.fetch(true)
			tp.operators.Push(tp.itemFromActive(nodes.NewFunctionCall(alias)))
			return nil
		}
		tp.output.Push(tp.itemFromActive(nodes.NewIdentifier(tp.active.Value)))
	}
	return nil
}

func (tp *termParser) handleString() error {
	tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
		Type:     types.String,
		Data:     strings.Trim(tp.active.Value, "\""),
		Constant: true,
	})))
	return nil
}

func (tp *termParser) handleNumber() error {
	if strings.Contains(tp.active.Value, ".") {
		f, err := strconv.ParseFloat(tp.active.Value, 64)
		if err != nil {
			return err
		}
		tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
			Type:     types.Float,
			Data:     f,
			Constant: true,
		})))
	} else {
		i, err := strconv.ParseInt(tp.active.Value, 10, 64)
		if err != nil {
			return err
		}
		tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
			Type:     types.Integer,
			Data:     i,
			Constant: true,
		})))
	}
	return nil
}

func (tp *termParser) handleOperator() error {
	item := tp.itemFromActive(nil)
	item.Node = nodes.NewOperation(tp.active.Value, tp.argCount(item))
	for !tp.operators.Empty() {
		top := tp.operators.Peek()
		if top.Value.Type != tokens.Operator {
			break
		}
		if tp.priority(top) < tp.priority(item) {
			break
		}
		if tp.binding(top) {
			break
		}
		tp.operators.Pop()
		for i := 0; i < tp.argCount(top); i++ {
			if tp.output.Empty() {
				return newMissingOperandsException(tp.argCount(top), i)
			}
			top.Node.AddFront(tp.output.Peek().Node)
			tp.output.Pop()
		}
		tp.output.Push(top)
	}
	tp.operators.Push(item)
	return nil
}

func (tp *termParser) handleLeftParentheses() error {
	tp.operators.Push(tp.itemFromActive(nil))
	return nil
}

func (tp *termParser) handleRightParentheses() error {
	for !tp.operators.Empty() && tp.operators.Peek().Value.Type != tokens.LeftParentheses {
		top := tp.operators.Peek()
		tp.operators.Pop()
		for i := 0; i < tp.argCount(top); i++ {
			if tp.output.Empty() {
				return ParseException{"Missing operands"}
			}
			top.Node.AddFront(tp.output.Peek().Node)
			tp.output.Pop()
		}
		tp.output.Push(top)
	}
	if tp.operators.Empty() {
		return ParseException{"Missing closing bracket"}
	} else if tp.operators.Peek().Node != nil {
		top := tp.operators.Peek()
		if tp.previous.Type != tokens.LeftParentheses && !tp.output.Empty() {
			top.Node.AddBack(tp.output.Peek().Node)
			tp.output.Pop()
		}
		tp.output.Push(top)
	}
	tp.operators.Pop()
	return nil
}

func (tp *termParser) fetch(incr bool) {
	if incr {
		tp.index++
	}
	tp.active = tp.input[tp.index]
	if tp.index > 0 {
		tp.previous = tp.input[tp.index-1]
	} else {
		tp.previous = tokens.Token{}
	}
	if tp.index < tp.size-1 {
		tp.next = tp.input[tp.index+1]
	} else {
		tp.next = tokens.Token{}
	}
}

func (tp *termParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	tp.operators.Clear()
	tp.output.Clear()

	tp.input = input
	tp.index, tp.size = 0, len(input)
parser:
	for ; tp.index < tp.size; tp.index++ {
		tp.fetch(false)

		/*
			fmt.Println("------------------------------")
			fmt.Println("ACTIVE", tp.active, "PREVIOUS", tp.previous, "NEXT", tp.next)
			fmt.Println("OPERATORS", tp.operators)
			fmt.Println("OUTPUT", tp.output)
			fmt.Println("------------------------------")
		*/

		switch tp.active.Type {
		case tokens.Statement, tokens.RightBlock, tokens.LeftBlock:
			break parser
		case tokens.Separator:
			for !tp.operators.Empty() && tp.operators.Peek().Value.Type != tokens.LeftParentheses {
				top := tp.operators.Peek()
				tp.operators.Pop()
				for i := 0; i < tp.argCount(top); i++ {
					if tp.output.Empty() {
						return nil, 0, ParseException{"Missing operands"}
					}
					top.Node.AddFront(tp.output.Peek().Node)
					tp.output.Pop()
				}
				tp.output.Push(top)
			}
			if tp.operators.Empty() {
				return tp.output.Peek().Node, tp.index, nil
			} else if tp.operators.Peek().Node != nil {
				top := tp.operators.Peek()
				top.Node.AddBack(tp.output.Peek().Node)
				tp.output.Pop()
			}
		default:
			handler, ok := tp.handlers[tp.active.Type]
			if !ok {
				return nil, 0, ParseException{"Unexpected token type"}
			}
			if err := handler(); err != nil {
				return nil, 0, err
			}
		}
	}

	for !tp.operators.Empty() {
		top := tp.operators.Peek()
		tp.operators.Pop()
		switch top.Value.Type {
		case tokens.Operator:
			for i := 0; i < tp.argCount(top); i++ {
				if tp.output.Empty() {
					return nil, 0, ParseException{"Missing operands"}
				}
				top.Node.AddFront(tp.output.Peek().Node)
				tp.output.Pop()
			}
			tp.output.Push(top)
		default:
			return nil, 0, ParseException{"Expected closing bracket"}
		}
	}

	/*
		fmt.Println("------------------------------")
		fmt.Println("END OF TERM")
		fmt.Println("OPERATORS", tp.operators)
		fmt.Println("OUTPUT", tp.output)
		fmt.Println("------------------------------")
	*/

	return tp.output.Peek().Node, tp.index, nil
}

type termItem struct {
	Value, Previous, Next tokens.Token
	Node                  nodes.Node
}

func (item termItem) String() string {
	if item.Node == nil {
		return fmt.Sprintf("[=%s, %s <- -> %s]", item.Value, item.Previous.Value, item.Next.Value)
	}
	return fmt.Sprintf("[=%s, %s <- -> %s, n = %v]", item.Value, item.Previous.Value, item.Next.Value, item.Node.Name())
}

func newItemStack() *itemStack {
	return &itemStack{
		items: make([]termItem, 0),
		size:  0,
	}
}

type itemStack struct {
	items []termItem
	size  int
}

func (stack itemStack) String() string {
	return fmt.Sprintf("stack%v", stack.items)
}

func (stack *itemStack) Clear() {
	stack.items = make([]termItem, 0)
	stack.size = 0
}

func (stack *itemStack) Empty() bool {
	return stack.size < 1
}

func (stack *itemStack) Push(item termItem) {
	stack.items = append(stack.items, item)
	stack.size++
}

func (stack *itemStack) Peek() termItem {
	if stack.Empty() {
		return termItem{}
	}
	return stack.items[stack.size-1]
}

func (stack *itemStack) Pop() {
	if stack.Empty() {
		return
	}
	stack.items = stack.items[:stack.size-1]
	stack.size--
}

func (stack *itemStack) Size() int {
	return stack.size
}
