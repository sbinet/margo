package rot13

import "io"

// STARTROT13-FUNC OMIT

func rot13(b byte) byte {
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}
	return (b-a+13)%(z-a+1) + a
}

// ENDROT13-FUNC OMIT

// STARTROT13 OMIT

type reader struct {
	r io.Reader
}

// NewReader returns a new io.Reader whose data will be encrypted via rot13
func NewReader(r io.Reader) io.Reader {
	return reader{r}
}

// ENDROT13 OMIT

func (r reader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = rot13(p[i])
	}
	return n, err
}
