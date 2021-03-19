package gin_swagger

import "unsafe"

func stringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func containsString(array []string, val string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			return true
		}
	}
	return false
}
