package main

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func assertResultError[T any](result T, err error) T {
	assertError(err)
	return result
}

func assertCondition(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func use(v any) {
}
