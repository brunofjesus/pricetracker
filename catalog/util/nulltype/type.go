package nulltype

type NullBoolean int

const (
	UndefinedValue NullBoolean = iota - 1
	FalseValue
	TrueValue
)

func IsTrue(n NullBoolean) bool {
	return n == TrueValue
}

func IsFalse(n NullBoolean) bool {
	return n == FalseValue
}

func IsUndefined(n NullBoolean) bool {
	return n == UndefinedValue
}

func FromInt(value int) NullBoolean {
	switch value {
	case int(TrueValue):
		return NullBoolean(TrueValue)
	case int(FalseValue):
		return NullBoolean(FalseValue)
	default:
		return NullBoolean(UndefinedValue)
	}
}

func ToString(n NullBoolean) string {
	switch n {
	case TrueValue:
		return "true"
	case FalseValue:
		return "false"
	default:
		return "undefined"
	}
}
