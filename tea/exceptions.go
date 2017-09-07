package tea

import "fmt"

func ConstantException(alias string) error {
	return fmt.Errorf("ConstantException: %s can not be changed", alias)
}

func RuntimeException(err string) error {
	return fmt.Errorf("RuntimeException: %s", err)
}

func ValueReferenceException() error {
	return fmt.Errorf("ValueReferenceException: Can not assign value to reference")
}

func ReferenceValueException() error {
	return fmt.Errorf("ReferenceValueException: Can not assign reference to value")
}

func SearchSpaceException() error {
	return fmt.Errorf("SearchSpaceException: Can not assign values from different search spaces")
}

func NamespaceException(alias string) error {
	return fmt.Errorf("NamespaceException: %s not found in search space", alias)
}

func ArgumentException(expected, got int) error {
	return fmt.Errorf("ArgumentException: Unknown signature, expected %d got %d", expected, got)
}

func ArgumentCastException(expected, got *Datatype) error {
	return fmt.Errorf("ArgumentCastException: Unknown signature, expected %s got %s", expected, got)
}

func FunctionException(alias string) error {
	return fmt.Errorf("FunctionException: No matching signature for %s", alias)
}

func CastException(from, to *Datatype) error {
	return fmt.Errorf("CastException: %s can not be implicitly casted to %s", from, to)
}

func StoreException(t interface{}) error {
	return fmt.Errorf("StoreException: Cannot store %s in namespace", t)
}
