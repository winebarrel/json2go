package json2go

type options struct {
	filename string
}

type optFn func(*options)

func OptionFilename(filename string) optFn {
	return func(o *options) {
		o.filename = filename
	}
}
