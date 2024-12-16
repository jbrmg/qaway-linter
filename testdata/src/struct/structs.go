package _struct

// Test is a struct
type Test struct { // want `Attribute 'MissingComment' is missing required comment`
	MissingComment string
}

type TestWithoutComment struct { // want `Struct 'TestWithoutComment' is missing required headline comment`
	// Comment
	MethodWithComment string
}
