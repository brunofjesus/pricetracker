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
