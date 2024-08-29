package iox

type AnyWriter interface {
	WriteAny(v any) (n int, err error)
}
