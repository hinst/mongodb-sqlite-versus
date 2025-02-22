package main

import "os"

func ReadStringFromFile(path string) string {
	return string(assertResultError(os.ReadFile(path)))
}

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	var isMissing = os.IsNotExist(err)
	return !isMissing
}
