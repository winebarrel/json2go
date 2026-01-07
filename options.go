package json2go

type options struct {
	filename string
	flat     bool
}

type OptFn func(*options)

func OptionFilename(filename string) OptFn {
	return func(o *options) {
		o.filename = filename
	}
}

func OptionFlat(flat bool) OptFn {
	return func(o *options) {
		o.flat = flat
	}
}
