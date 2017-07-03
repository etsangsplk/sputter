package api

// ExpectedBool is thrown when a Value is not a Bool
const ExpectedBool = "value is not a bool: %s"

// Comparison represents the result of a equality comparison
type Comparison int

// Name is a Variable name
type Name string

// Bool represents the values True or False
type Bool bool

// Names represents a set of Names
type Names []Name

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Value is the generic interface for all 'Values'
type Value interface {
	Str() Str
}

// ValueProcessor is the standard function interface for a func that
// processes a Value against a Context (example: Emit)
type ValueProcessor func(Context, Value) Value

// Comparer is an interface for a Value capable of comparing.
type Comparer interface {
	Compare(Comparer) Comparison
}

// Named is the generic interface for Values that are named
type Named interface {
	Name() Name
}

// Typed is the generic interface for Values that are typed
type Typed interface {
	Type() Name
}

// Documented is the generic interface for Values that are documented
type Documented interface {
	Documentation() Str
}

// Getter is the interface for Values that have retrievable Properties
type Getter interface {
	Get(Value) (Value, bool)
}

// Elementer is the interface for Values that have indexed elements
type Elementer interface {
	ElementAt(int) (Value, bool)
}

type nilValue struct{}

const (
	// LessThan means left Value is less than right Value
	LessThan Comparison = -1

	// EqualTo means left Value is equal to right Value
	EqualTo Comparison = 0

	// GreaterThan means left Value is greater than right Value
	GreaterThan Comparison = 1
)

var (
	// True represents the boolean value of True
	True = Bool(true)

	// False represents the boolean value of false
	False = Bool(false)

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
func (b Bool) Apply(c Context, args Sequence) Value {
	for i := args; i.IsSequence(); i = i.Rest() {
		if i.First() != b {
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

func (n *nilValue) Apply(c Context, args Sequence) Value {
	for i := args; i.IsSequence(); i = i.Rest() {
		if i.First() != Nil {
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
	switch {
	case v == False || v == Nil:
		return false
	default:
		return true
	}
}

// AssertBool will cast a Value into a Bool or explode violently
func AssertBool(v Value) Bool {
	if b, ok := v.(Bool); ok {
		return b
	}
	panic(Err(ExpectedBool, v))
}
