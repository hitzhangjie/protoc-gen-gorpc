package gorpc

import (
	"strings"

	"github.com/iancoleman/strcase"
)

var FuncMap = map[string]interface{}{
	"simplify":   PBSimplifyGoType,
	"gopkg":      PBGoPackage,
	"gotype":     PBGoType,
	"export":     GoExport,
	"gofulltype": GoFullyQualifiedType,
	"title":      Title,
	"untitle":    UnTitle,
	"trimright":  TrimRight,
	"splitList":  SplitList,
	"last":       Last,
	"hasprefix":  HasPrefix,
	"camelcase":  strcase.ToCamel,
	"snakecase":  strcase.ToSnake,
	"contains":   strings.Contains,
	"add":        Add,
}
