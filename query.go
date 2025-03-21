package pgtmpl

import (
	"strings"
	"sync"
)

type Query struct {
	strings.Builder
	args []any
	name string
}

var queryPool = sync.Pool{New: func() any { return new(Query) }}

func NewQuery() *Query {
	var q, _ = queryPool.Get().(*Query)
	return q
}
func (q *Query) Cap() int {
	return q.Builder.Cap()
}
func (q *Query) Len() int {
	return q.Builder.Len()
}
func (q *Query) Grow(n int) {
	q.Builder.Grow(n)
}
func (q *Query) Write(p []byte) (int, error) {
	return q.Builder.Write(p)
}
func (q *Query) WriteByte(c byte) error {
	return q.Builder.WriteByte(c)
}
func (q *Query) WriteRune(r rune) (int, error) {
	return q.Builder.WriteRune(r)
}
func (q *Query) WriteString(s string) (int, error) {
	return q.Builder.WriteString(s)
}
func (q *Query) Reset() {
	if q == nil {
		return
	}
	q.Builder.Reset()
	q.name = ""
	if q.args != nil {
		q.args = q.args[:0]
	}
}
func (q *Query) Close() {
	if q == nil {
		return
	}
	q.Reset()
	queryPool.Put(q)
}
func (q *Query) String() string {
	if q == nil {
		return ""
	}
	return q.Builder.String()
}
func (q *Query) Args() []any {
	if q == nil {
		return nil
	}
	return q.args
}
func (q *Query) Name() string {
	if q == nil {
		return ""
	}
	return q.name
}
