package cascade

import (
	"flag"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `check cascades calls in nested structs

This analyzer finds struct methods and verifies that all struct fields 
whose types also have the same method are called.

Use -method flag to specify the method name (default: WithDefaults).

Example of code that would be flagged:

	type Outer struct {
		Inner Inner
	}

	func (s Outer) WithDefaults() Outer {
		// Missing: s.Inner = s.Inner.WithDefaults()
		return s
	}

	type Inner struct {
		Value int
	}

	func (s Inner) WithDefaults() Inner { return s }
`

var methodName string

var Analyzer = &analysis.Analyzer{
	Name:     "cascade",
	Doc:      doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func init() {
	Analyzer.Flags.Init("cascade", flag.ExitOnError)
	Analyzer.Flags.StringVar(&methodName, "method", "WithDefaults", "method name to check for cascading calls")
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)

		if fn.Recv == nil || fn.Name.Name != methodName {
			return
		}

		recvType := getReceiverType(pass, fn)
		if recvType == nil {
			return
		}

		structType, ok := recvType.Underlying().(*types.Struct)
		if !ok {
			return
		}

		fieldsWithMethod := findFieldsWithMethod(structType, methodName)
		if len(fieldsWithMethod) == 0 {
			return
		}

		calledFields := findCalledMethod(fn, methodName)

		for _, field := range fieldsWithMethod {
			if !calledFields[field] {
				pass.Reportf(fn.Pos(), "%s.%s() does not call %s.%s()", recvType.Obj().Name(), methodName, field, methodName)
			}
		}
	})

	return nil, nil
}

func getReceiverType(pass *analysis.Pass, fn *ast.FuncDecl) *types.Named {
	if len(fn.Recv.List) == 0 {
		return nil
	}

	recvExpr := fn.Recv.List[0].Type

	if star, ok := recvExpr.(*ast.StarExpr); ok {
		recvExpr = star.X
	}

	recvTypeInfo := pass.TypesInfo.TypeOf(recvExpr)
	if recvTypeInfo == nil {
		return nil
	}

	named, _ := recvTypeInfo.(*types.Named)
	return named
}

func findFieldsWithMethod(structType *types.Struct, method string) []string {
	var fields []string

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldType := field.Type()

		if ptr, ok := fieldType.(*types.Pointer); ok {
			fieldType = ptr.Elem()
		}

		if hasMethod(fieldType, method) {
			fields = append(fields, field.Name())
		}
	}

	return fields
}

func hasMethod(t types.Type, method string) bool {
	if mset := types.NewMethodSet(t); mset.Lookup(nil, method) != nil {
		return true
	}

	if _, ok := t.(*types.Pointer); !ok {
		ptrMset := types.NewMethodSet(types.NewPointer(t))
		if ptrMset.Lookup(nil, method) != nil {
			return true
		}
	}

	return false
}

func findCalledMethod(fn *ast.FuncDecl, method string) map[string]bool {
	called := make(map[string]bool)

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// Look for calls like: s.Field.Method() or s.Field = s.Field.Method()
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if it's a selector expression (method call)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != method {
			return true
		}

		// Get the field name from the selector chain
		// e.g., s.Field.Method() -> Field
		fieldName := extractFieldName(sel.X)
		if fieldName != "" {
			called[fieldName] = true
		}

		return true
	})

	return called
}

func extractFieldName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		// This is s.Field - return Field
		return e.Sel.Name
	case *ast.Ident:
		// This might be just the receiver (s), skip
		return ""
	}
	return ""
}
