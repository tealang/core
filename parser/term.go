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

type functionCallParser struct{}

func (functionCallParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	return nil, 0, nil
}

func newTermParser() *termParser {
	return &termParser{
		output:    newItemStack(),
		operators: newItemStack(),
	}
}

type termParser struct {
	output, operators      *itemStack
	active, previous, next tokens.Token
}

func (termParser) binding(item termItem) bool {
	switch item.Value.Value {
	case "^":
		return true
	default:
		return false
	}
}

func (termParser) argCount(item termItem) int {
	return 2
}

func (termParser) priority(item termItem) int {
	switch item.Value.Value {
	case "&", "|":
		return 8
	case "!":
		return 7
	case "^":
		return 6
	case "*", "/":
		return 5
	case "+", "-", ":":
		return 4
	case "%":
		return 3
	case "<", ">", ">=", "<=", "=>", "=<", "!=", "==":
		return 2
	case "&&", "||", "^|":
		return 1
	default:
		return 0
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

func (tp *termParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	tp.operators.Clear()
	tp.output.Clear()

	var (
		index, size = 0, len(input)
	)
parser:
	for ; index < size; index++ {
		tp.active = input[index]
		if index > 0 {
			tp.previous = input[index-1]
		}
		if index < size-1 {
			tp.next = input[index+1]
		}

		fmt.Println("------------------------------")
		fmt.Println("ACTIVE", tp.active, "PREVIOUS", tp.previous, "NEXT", tp.next)
		fmt.Println("OPERATORS", tp.operators)
		fmt.Println("OUTPUT", tp.output)
		fmt.Println("------------------------------")

		switch tp.active.Type {
		case tokens.Statement:
			break parser
		case tokens.Identifier:
			switch tp.active.Value {
			case trueLiteral:
				tp.output.Push(tp.itemFromActive(nodes.NewLiteral(types.True)))
			case falseLiteral:
				tp.output.Push(tp.itemFromActive(nodes.NewLiteral(types.False)))
			case nullLiteral:
				tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{})))
			default:
				tp.output.Push(tp.itemFromActive(nodes.NewIdentifier(tp.active.Value)))
			}
		case tokens.String:
			tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
				Type:     types.String,
				Data:     tp.active.Value,
				Constant: true,
			})))
		case tokens.Number:
			if strings.Contains(tp.active.Value, ".") {
				f, err := strconv.ParseFloat(tp.active.Value, 64)
				if err != nil {
					return nil, 0, err
				}
				tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
					Type:     types.Float,
					Data:     f,
					Constant: true,
				})))
			} else {
				i, err := strconv.ParseInt(tp.active.Value, 10, 64)
				if err != nil {
					return nil, 0, err
				}
				tp.output.Push(tp.itemFromActive(nodes.NewLiteral(runtime.Value{
					Type:     types.Integer,
					Data:     i,
					Constant: true,
				})))
			}
		case tokens.Operator:
			item := tp.itemFromActive(nil)
			item.Node = nodes.NewOperation(tp.active.Value, tp.argCount(item))
			for !tp.operators.Empty() {
				top := tp.operators.Peek()
				if top.Value.Type != tokens.Operator {
					break
				}
				if tp.priority(top) >= tp.priority(item) {
					break
				}
				if tp.binding(top) {
					break
				}
				for i := 0; i < tp.argCount(top); i++ {
					if tp.output.Empty() {
						return nil, 0, ParseException{"Missing operands"}
					}
					top.Node.AddFront(tp.output.Peek().Node)
					tp.output.Pop()
				}
				tp.output.Push(top)
			}
			tp.operators.Push(item)
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

	return tp.output.Peek().Node, index, nil
}

type termItem struct {
	Value, Previous, Next tokens.Token
	Node                  nodes.Node
}

func (item termItem) String() string {
	return fmt.Sprintf("[=%s, previous %s, next %s, node = %v]", item.Value, item.Previous, item.Next, item.Node.Name())
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
