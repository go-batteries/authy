//simplest package ever to put the table name in place
// Usage
// bob.Table(&User{}).Build(`SELECT * FROM :table`)
package bob

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type Bob struct {
	tableName string
	partQuery string
}

// Table accepts a pointer to a struct
// the struct should have TableName() string
// method defined on it.
// It calls the TableName() method to get the tablename
// or defaults to snakecaseing
func Table(i interface{}) *Bob {
	value := reflect.ValueOf(i).Elem()
	tableName := reflect.TypeOf(i).Elem().Name()

	method := value.MethodByName("TableName")
	tableName = fmt.Sprintf("%ss", strings.ToLower(strcase.ToSnake(tableName)))
	if !method.IsValid() {
		return &Bob{tableName: tableName}
	}

	v, ok := method.Call([]reflect.Value{})[0].Interface().(string)
	if ok {
		tableName = v
	}

	return &Bob{tableName: tableName}
}

func (d *Bob) TableName() string {
	return d.tableName
}

// BuildQuery takes a string with :table template
// and replaces it in query
func (d *Bob) Build(query string, pattern ...string) string {
	pat := "table"
	if len(pattern) != 0 {
		pat = pattern[0]
	}

	re := regexp.MustCompile(fmt.Sprintf(`(?m):%s\b`, pat))
	return re.ReplaceAllString(query, d.tableName)
}

func (d *Bob) BuildWithClause(query string, pattern ...string) *Bob {
	q := d.Build(query, pattern...)
	return &Bob{tableName: d.tableName, partQuery: q}
}

func (d *Bob) Assemble(m map[string]interface{}) string {
	q := d.partQuery
	parts := []string{}

	for k := range m {
		parts = append(parts, fmt.Sprintf("%s=:%s", k, k))
	}

	return fmt.Sprintf(q, strings.Join(parts, ", "))
}
