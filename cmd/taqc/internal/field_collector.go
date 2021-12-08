package internal

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"

	"github.com/moznion/taqc"
)

// Field represents a field of the structure for a constructor to be generated.
type Field struct {
	// FieldName is a name of the field.
	FieldName string
	// FieldType is a type of the field.
	FieldType string
	// ParamName is a query parameter's name
	ParamName string
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

				return convertStructFieldsToQueryParamFields(structType.Fields.List), nil
			}
		}
	}

	return nil, fmt.Errorf("there is no suitable struct that matches given typeName [given=%s]", typeName)

}

func convertStructFieldsToQueryParamFields(fields []*ast.Field) []*Field {
	fs := make([]*Field, 0)
	for _, field := range fields {
		if field.Tag == nil || len(field.Tag.Value) <= 0 {
			continue
		}

		customTag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		paramName := customTag.Get(taqc.TagName)
		if paramName == "" {
			continue
		}

		fieldType := types.ExprString(field.Type)

		var fieldName string
		if len(field.Names) <= 0 {
			continue
		}
		fieldName = field.Names[0].Name

		fs = append(fs, &Field{
			FieldName: fieldName,
			FieldType: fieldType,
			ParamName: paramName,
		})
	}
	return fs
}
