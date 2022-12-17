package utils

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Numbers interface {
	Signed | Unsigned | Float
}

type Integer interface {
	Signed | Unsigned
}

type Ordered interface {
	Integer | Float | ~string
}

func PercentOf[T Numbers, E Numbers](part T, total E) float64 {
	return (float64(part) * float64(100)) / float64(total)
}

func Percent[T Numbers, E Numbers](percent T, all E) float64 {
	return ((float64(all) * float64(percent)) / float64(100))
}
