package utils

// AssertEqual check if two things are equal and if not raise error
func AssertEqual(a interface{}, b interface{}) {
	if a != b {
		panic("Assertion Fail!")
	}
}

// AssertTrue check if something is true and if not raise error
func AssertTrue(a bool) {
	if a == false {
		panic("Assertion Fail!")
	}
}
