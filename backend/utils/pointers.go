package utils

// strPtr รับค่า string แล้วคืนค่าเป็น pointer ของ string นั้น
func StrPtr(s string) *string  { 
	return &s 
}

// intPtr รับค่า int แล้วคืนค่า pointer ของ int
func IntPtr(i int) *int        { 
	return &i 
}

// uintPtr รับค่า uint แล้วคืนค่า pointer ของ uint
func UintPtr(i uint) *uint     { 
	return &i 
}

// boolPtr รับค่า bool แล้วคืนค่า pointer ของ bool
func BoolPtr(b bool) *bool     { 
	return &b 
}

// floatPtr รับค่า float64 แล้วคืนค่า pointer ของ float64
func FloatPtr(f float64) *float64 { 
	return &f 
}
