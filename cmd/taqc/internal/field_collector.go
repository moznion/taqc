package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	g "github.com/moznion/gowrtr/generator"
	"github.com/moznion/taqc"
	"github.com/moznion/taqc/internal"
)

// Field represents a field of the structure for a constructor to be generated.
type Field struct {
	// FieldName is a name of the field.
	FieldName string
	// FieldType is a type of the field.
	FieldType string
	// ParamName is a query parameter's name
	ParamName string

	TimeFormatterStmt g.Statement
}

func CollectQueryParameterFieldsFromAST(typeName string, astFiles []*ast.File) ([]*Field, error) {
	for _, astFile := range astFiles {
		for _, decl := range astFile.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structName := typeSpec.Name.Name
				if typeName != structName {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				fields, err := convertStructFieldsToQueryParamFields(structType.Fields.List)
				if err != nil {
					return nil, err
				}
				return fields, nil
			}
		}
	}

	return nil, fmt.Errorf("there is no suitable struct that matches given typeName [given=%s]", typeName)

}

func convertStructFieldsToQueryParamFields(fields []*ast.Field) ([]*Field, error) {
	fs := make([]*Field, 0)
	for _, field := range fields {
		if field.Tag == nil || len(field.Tag.Value) <= 0 {
			continue
		}

		customTag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		tag := customTag.Get(internal.TagName)
		if tag == "" {
			continue
		}

		splitTagValues := strings.Split(tag, ",")

		paramName := strings.TrimSpace(splitTagValues[0])
		timeLayout, unixTimeUnit := internal.ExtractTimeTag(splitTagValues[1:])
		timeFormatterStmtBase := g.NewAnonymousFunc(false, g.NewAnonymousFuncSignature().AddParameters(g.NewFuncParameter("t", "time.Time")).ReturnTypes("string"))

		timeFormatterStmt := timeFormatterStmtBase.Statements(g.NewReturnStatement(`fmt.Sprintf("%d", t.Unix())`))
		if unixTimeUnit != "" {
			switch unixTimeUnit {
			case "sec":
				// do nothing; using default stmt
			case "millisec":
				timeFormatterStmt = timeFormatterStmtBase.Statements(g.NewReturnStatement(`fmt.Sprintf("%d", t.UnixMilli())`))
			case "microsec":
				timeFormatterStmt = timeFormatterStmtBase.Statements(g.NewReturnStatement(`fmt.Sprintf("%d", t.UnixMicro())`))
			case "nanosec":
				timeFormatterStmt = timeFormatterStmtBase.Statements(g.NewReturnStatement(`fmt.Sprintf("%d", t.UnixNano())`))
			default:
				return nil, fmt.Errorf("%s is unsupported: %w", unixTimeUnit, taqc.ErrUnsupportedUnixTimeUnit)
			}
		}
		if timeLayout != "" { // higher priority
			timeFormatterStmt = timeFormatterStmtBase.Statements(g.NewReturnStatement(fmt.Sprintf(`t.Format("%s")`, timeLayout)))
		}

		fieldType := types.ExprString(field.Type)

		var fieldName string
		if len(field.Names) <= 0 {
			continue
		}
		fieldName = field.Names[0].Name

		fs = append(fs, &Field{
			FieldName:         fieldName,
			FieldType:         fieldType,
			ParamName:         paramName,
			TimeFormatterStmt: timeFormatterStmt,
		})
	}
	return fs, nil
}
