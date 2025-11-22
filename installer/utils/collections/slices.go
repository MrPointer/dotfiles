package collections

import "errors"

// Last returns the last element of a slice. If the slice is empty, it returns an error.
// This function is generic and works with any type of slice.
// It is useful for getting the last element of a slice without needing to know the type in advance.
//
// Example usage:
//
//	slice := []int{1, 2, 3}
//	lastElement, err := Last(slice)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(lastElement) // Output: 3
func Last[E any](s []E) (E, error) {
	if len(s) == 0 {
		var zero E
		return zero, errors.New("slice is empty")
	}
	return s[len(s)-1], nil
}
