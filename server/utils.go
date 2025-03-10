package main

func RandomString(length int) string {
	result := ""
	for range length {
		randomIndex := randRange(0, len(StringGenerationAlphabet)-1)
		result += string(StringGenerationAlphabet[randomIndex])
	}
	return result
}
