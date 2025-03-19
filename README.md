# pgtmpl - golang template engine for building postgresql queries

## What is the difference from text/template?

There is only one difference. It's support for replacing values with placeholders.

## Usage

`/sql.tpl`
```sql
{{define "select_author"}}
    SELECT * FROM author WHERE id = {{hold $v}};
{{end}}
```
`/main.go`
```go
package main
import (
    "github.com/typomaker/pgtmpl"
    "embed"
    "fmt"
)

//go:embed sql.tpl
var sqltpl embed.FS
var tpl = pgtmpl.Must(pgtmpl.ParseFS(sqltpl, "*"))
func main() {
    var query pgtmpl.Query{}
    var userID = 777
    if err := tpl.ExecuteTemplate(&query, "select_author", userID); err != nil {
        panic(err)
    }
    fmt.Println(query.String())
    fmt.Println(query.Args())
    // Output:
    // SELECT * FROM author WHERE id = $1
    // 777
}
```