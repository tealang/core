package runtime

import "testing"

func TestConstantException_Error(t *testing.T) {
	type fields struct {
		Alias string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{"Name"}, "ConstantException: Name can not be changed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConstantException{
				Alias: tt.fields.Alias,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("ConstantException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuntimeException_Error(t *testing.T) {
	type fields struct {
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{"Message"}, "RuntimeException: Message"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Exception{
				Message: tt.fields.Message,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("RuntimeException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueReferenceException_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Error string", "ValueReferenceException: Can not assign value to reference"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ValueReferenceException{}
			if got := c.Error(); got != tt.want {
				t.Errorf("ValueReferenceException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReferenceValueException_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Error string", "ReferenceValueException: Can not assign reference to value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ReferenceValueException{}
			if got := c.Error(); got != tt.want {
				t.Errorf("ReferenceValueException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchSpaceException_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Error string", "SearchSpaceException: Can not assign values from different search spaces"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := SearchSpaceException{}
			if got := c.Error(); got != tt.want {
				t.Errorf("SearchSpaceException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceException_Error(t *testing.T) {
	type fields struct {
		Alias string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{"Alias"}, "NamespaceException: Alias not found in search space"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NamespaceException{
				Alias: tt.fields.Alias,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("NamespaceException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgumentException_Error(t *testing.T) {
	type fields struct {
		Expected int
		Got      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{2, 1}, "ArgumentException: Unknown signature, expected 2 got 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ArgumentException{
				Expected: tt.fields.Expected,
				Got:      tt.fields.Got,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("ArgumentException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgumentCastException_Error(t *testing.T) {
	type fields struct {
		Expected *Datatype
		Got      *Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{&Datatype{Name: "dt1"}, &Datatype{Name: "dt2"}}, "ArgumentCastException: Unknown signature, expected dt1 got dt2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ArgumentCastException{
				Expected: tt.fields.Expected,
				Got:      tt.fields.Got,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("ArgumentCastException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionException_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Error string", "FunctionException: No matching signature found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FunctionException{}
			if got := c.Error(); got != tt.want {
				t.Errorf("FunctionException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastException_Error(t *testing.T) {
	type fields struct {
		From *Datatype
		To   *Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{&Datatype{Name: "dt1"}, &Datatype{Name: "dt2"}}, "CastException: dt1 can not be implicitly casted to dt2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CastException{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("CastException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignmentMismatchException_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Error string", "AssignmentMismatchException: Number of aliases must match number of values"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := AssignmentMismatchException{}
			if got := c.Error(); got != tt.want {
				t.Errorf("AssignmentMismatchException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnexpectedItemException_Error(t *testing.T) {
	type fields struct {
		Expected interface{}
		Got      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{1, "string"}, "UnexpectedItemException: Expected search item of type int, got type string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := UnexpectedItemException{
				Expected: tt.fields.Expected,
				Got:      tt.fields.Got,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("UnexpectedItemException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUncallableTypeException_Error(t *testing.T) {
	type fields struct {
		Type *Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{&Datatype{Name: "dt1"}}, "UncallableTypeException: Value of type dt1 can not be called"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := UncallableTypeException{
				Type: tt.fields.Type,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("UncallableTypeException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExplicitCastException_Error(t *testing.T) {
	type fields struct {
		From *Datatype
		To   *Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{&Datatype{Name: "dt1"}, &Datatype{Name: "dt2"}}, "ExplicitCastException: dt1 can not be explicitly casted to dt2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ExplicitCastException{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("ExplicitCastException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreException_Error(t *testing.T) {
	type fields struct {
		Item interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Error string", fields{Value{}}, "StoreException: Cannot store null in namespace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := StoreException{
				Item: tt.fields.Item,
			}
			if got := c.Error(); got != tt.want {
				t.Errorf("StoreException.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
