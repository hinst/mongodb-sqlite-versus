package main

import (
	"fmt"
	"os"
	"time"
)

func testJson(users []*User) {
	const JSON_FILE_PATH = "jsonFile.json"
	var beginning = time.Now()
	writeJsonToFile(JSON_FILE_PATH, false, users)
	var elapsed = time.Since(beginning)
	var fileInfo = assertResultError(os.Stat(JSON_FILE_PATH))
	fmt.Printf("JSON time: %v, file size: %v\n", elapsed, formatFileSize(fileInfo.Size()))
}
