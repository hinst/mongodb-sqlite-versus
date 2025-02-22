package main

func AssertError(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertResultError[T any](result T, err error) T {
	AssertError(err)
	return result
}

func AssertCondition(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
