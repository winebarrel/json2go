package json2go_test

func ptr[T any](v T) *T {
	return &v
}
