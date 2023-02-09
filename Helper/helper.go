package Helper

import "strconv"

func StringToInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	value := int32(i)
	return value
}
