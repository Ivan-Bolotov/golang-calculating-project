package functions

// MinMax служит для нахождения минимального и максимального значения в слайсе []int
func MinMax(array []int) (int, int) {
	var maximum int = array[0]
	var minimum int = array[0]
	for _, value := range array {
		if maximum < value {
			maximum = value
		}
		if minimum > value {
			minimum = value
		}
	}
	return minimum, maximum
}
