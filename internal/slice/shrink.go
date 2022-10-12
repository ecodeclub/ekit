package slice

func CalCapacity(c, l int) int {
	if c <= 64 {
		return c
	}
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor))
	}
	if c <= 2048 && (c/l >= 4) {
		return c / 2
	}
	return c
}

func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n := CalCapacity(c, l)
	if n == c {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}
