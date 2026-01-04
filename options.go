package json2go

type options struct {
	filename string
}

type optionFunc func(*options)

func WithFilename(filename string) optionFunc {
	return func(o *options) {
		o.filename = filename
	}
}
