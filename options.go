package json2go

type options struct {
	filename string
}

type OptionFunc func(*options)

func WithFilename(filename string) OptionFunc {
	return func(o *options) {
		o.filename = filename
	}
}
