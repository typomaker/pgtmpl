// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pgtmpl_test

import (
	"log"

	"github.com/typomaker/pgtmpl"
)

func ExampleTemplate() {
	const tpl = `SELECT * FROM author WHERE id IN(
					{{- range $i, $v := . }} 
						{{- if $i}},{{end}} 
						{{- hold $v}} 
					{{- end -}}
				)`

	var (
		authorIDs = []int{100, 12, 334}
		query     = pgtmpl.Query{}
	)

	t := pgtmpl.Must(pgtmpl.New("author_by_id").Parse(tpl))

	err := t.Execute(&query, authorIDs)
	if err != nil {
		log.Println("executing template:", err)
	}

	printQuery(query)

	// Output:
	// TEXT: SELECT * FROM author WHERE id IN($1,$2,$3)
	// ARGUMENT: 100 12 334
}
