package json2go

import (
	"bytes"
	"container/list"
	"iter"
)

type buffers struct {
	l *list.List
}

func newBuffers() *buffers {
	return &buffers{l: list.New()}
}

func (bs *buffers) Push(v *bytes.Buffer) {
	bs.l.PushBack(v)
}

func (bs *buffers) Pop() *bytes.Buffer {
	if e := bs.l.Back(); e != nil {
		return bs.l.Remove(e).(*bytes.Buffer)
	} else {
		return nil
	}
}

func (bs *buffers) Write(p []byte) (n int, err error) {
	b := bs.l.Back().Value.(*bytes.Buffer)
	return b.Write(p)
}

func (bs *buffers) Iter() iter.Seq[*bytes.Buffer] {
	return func(yield func(*bytes.Buffer) bool) {
		for e := bs.l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value.(*bytes.Buffer)) {
				break
			}
		}
	}
}
