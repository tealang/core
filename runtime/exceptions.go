package runtime

import "fmt"

type ConstantException struct {
	Alias string
}

func (c ConstantException) Error() string {
	return fmt.Sprintf("ConstantException: %s can not be changed", c.Alias)
}

type RuntimeException struct {
	Message string
}

func (c RuntimeException) Error() string {
	return fmt.Sprintf("RuntimeException: %s", c.Message)
}

type ValueReferenceException struct{}

func (c ValueReferenceException) Error() string {
	return fmt.Sprintf("ValueReferenceException: Can not assign value to reference")
}

type ReferenceValueException struct{}

func (c ReferenceValueException) Error() string {
	return fmt.Sprintf("ReferenceValueException: Can not assign reference to value")
}

type SearchSpaceException struct{}

func (c SearchSpaceException) Error() string {
	return fmt.Sprintf("SearchSpaceException: Can not assign values from different search spaces")
}

type NamespaceException struct {
	Alias string
}

func (c NamespaceException) Error() string {
	return fmt.Sprintf("NamespaceException: %s not found in search space", c.Alias)
}

type ArgumentException struct {
	Expected, Got int
}

func (c ArgumentException) Error() string {
	return fmt.Sprintf("ArgumentException: Unknown signature, expected %d got %d", c.Expected, c.Got)
}

type ArgumentCastException struct {
	Expected, Got *Datatype
}

func (c ArgumentCastException) Error() string {
	return fmt.Sprintf("ArgumentCastException: Unknown signature, expected %s got %s", c.Expected, c.Got)
}

type FunctionException struct {
}

func (c FunctionException) Error() string {
	return fmt.Sprintf("FunctionException: No matching signature found")
}

type CastException struct {
	From, To *Datatype
}

func (c CastException) Error() string {
	return fmt.Sprintf("CastException: %s can not be implicitly casted to %s", c.From, c.To)
}

type AssignmentMismatchException struct {
}

func (c AssignmentMismatchException) Error() string {
	return fmt.Sprintf("AssignmentMismatchException: Number of aliases must match number of values")
}

type UnexpectedItemException struct {
	Expected, Got interface{}
}

func (c UnexpectedItemException) Error() string {
	return fmt.Sprintf("UnexpectedItemException: Expected search item of type %T, got type %T", c.Expected, c.Got)
}

type UncallableTypeException struct {
	Type *Datatype
}

func (c UncallableTypeException) Error() string {
	return fmt.Sprintf("UncallableTypeException: Value of type %s can not be called", c.Type)
}

type ExplicitCastException struct {
	From, To *Datatype
}

func (c ExplicitCastException) Error() string {
	return fmt.Sprintf("ExplicitCastException: %s can not be explicitly casted to %s", c.From, c.To)
}

type StoreException struct {
	Item interface{}
}

func (c StoreException) Error() string {
	return fmt.Sprintf("StoreException: Cannot store %s in namespace", c.Item)
}
