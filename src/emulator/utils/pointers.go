package utils

// 🟦 Develop pointers to primitives

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}
