package util

func NewTruePtr() *bool {
	val := true

	return &val
}

func NewFalsePtr() *bool {
	val := false

	return &val
}
