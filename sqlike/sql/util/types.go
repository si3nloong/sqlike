package util

// StringSlice :
type StringSlice []string

// IndexOf :
func (slice StringSlice) IndexOf(search string) (idx int) {
	idx = -1
	length := len(slice)
	for i := 0; i < length; i++ {
		if slice[i] == search {
			idx = i
			break
		}
	}
	return
}

// Splice :
func (slice *StringSlice) Splice(idx int) {
	*slice = append((*slice)[:idx], (*slice)[idx+1:]...)
}