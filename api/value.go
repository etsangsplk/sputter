package api

const (
	// LessThan means left Value is less than right Value
	LessThan Comparison = -1

	// EqualTo means left Value is equal to right Value
	EqualTo Comparison = 0

	// GreaterThan means left Value is greater than right Value
	GreaterThan Comparison = 1
)

type (
	// Bool represents the values True or False
	Bool bool

	// Value is the generic interface for all 'Vector'
	Value interface {
		Str() Str
	}

	// ValueProcessor is the standard function interface for a func
	// that processes a Value against a Context (example: Emit)
	ValueProcessor func(Context, Value) Value

	// Comparison represents the result of a equality comparison
	Comparison int

	// Comparer is an interface for a Value capable of comparing.
	Comparer interface {
		Compare(Comparer) Comparison
	}

	// Name is a Variable name
	Name string

	// Names represents a set of Names
	Names []Name

	// Typed is the generic interface for Vector that are typed
	Typed interface {
		Type() Name
	}

	// Documented is the generic interface for Vector that are documented
	Documented interface {
		Documentation() Str
	}

	// Counted interfaces allow a Value to return a count of its items
	Counted interface {
		Count() int
	}

	// Mapped is the interface for Vector that have retrievable Properties
	Mapped interface {
		Get(Value) (Value, bool)
	}

	// Indexed is the interface for Vector that have indexed elements
	Indexed interface {
		ElementAt(int) (Value, bool)
	}

	nilValue struct{}
)

var (
	// True represents the boolean value of True
	True Bool = true

	// False represents the boolean value of false
	False Bool = false

	// Nil is a value that represents the absence of a Value
	Nil = &nilValue{}
)

// Name makes Name Named
func (n Name) Name() Name {
	return n
}

// Str converts this Value into a Str
func (n Name) Str() Str {
	return Str(n)
}

// Apply makes Bool Applicable
func (b Bool) Apply(_ Context, args Vector) Value {
	for _, f := range args {
		if f != b {
			return False
		}
	}
	return True
}

// Str converts this Value into a Str
func (b Bool) Str() Str {
	if bool(b) {
		return "true"
	}
	return "false"
}

func (n *nilValue) Apply(_ Context, args Vector) Value {
	for _, f := range args {
		if f != Nil {
			return False
		}
	}
	return True
}

func (n *nilValue) Str() Str {
	return "nil"
}

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	if v == False || v == Nil {
		return false
	}
	return true
}
