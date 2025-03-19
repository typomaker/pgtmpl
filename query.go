package pgtmpl

import "strings"

type Query struct {
	strings.Builder
	args []interface{}
}

func (q *Query) Reset() {
	q.Builder.Reset()
	q.args = nil
}
func (q Query) String() string {
	return q.Builder.String()
}
func (q Query) Args() []interface{} {
	return q.args
}
