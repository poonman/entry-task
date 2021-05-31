package helper

func NewString(b byte) string {
	str := make([]byte, 2000, 2000)
	for i := range str {
		str[i] = b
	}

	return string(str)
}
