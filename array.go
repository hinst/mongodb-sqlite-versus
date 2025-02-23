package main

// Divide array into N parts
func divideArray[Slice ~[]E, E any](array Slice, n int) []Slice {
	var outputs = make([]Slice, n)
	var outputIndex = 0
	for _, item := range array {
		outputs[outputIndex] = append(outputs[outputIndex], item)
		outputIndex++
		if outputIndex >= n {
			outputIndex = 0
		}
	}
	return outputs
}
