package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	humanize "github.com/dustin/go-humanize"
)

var executablePath = filepath.Dir(assertResultError(os.Executable()))

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

const TAB = "\t"

func readStringFromFile(path string) string {
	return string(assertResultError(os.ReadFile(path)))
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	var isMissing = os.IsNotExist(err)
	return !isMissing
}

func writeStringToFile(path string, text string) {
	assertError(os.WriteFile(path, []byte(text), OS_USER_RW))
}

func writeJsonToFile(path string, indent bool, obj any) {
	var bytes []byte
	if indent {
		bytes = assertResultError(json.MarshalIndent(obj, "", "\t"))
	} else {
		bytes = assertResultError(json.Marshal(obj))
	}
	assertError(os.WriteFile(path, bytes, OS_USER_RW))
}

func formatFileSize(size int64) string {
	return humanize.IBytes(uint64(size))
}
