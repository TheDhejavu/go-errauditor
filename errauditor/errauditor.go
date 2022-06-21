package errauditor

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/fatih/color"
)

type ErrorType string

type AggregatedError struct {
	Func   string
	Errors []string
}

type Result struct {
	AggregatedErrors  []*AggregatedError
	WrappedErrorCount int64
	ConstErrorCount   int64
}

const (
	Error   ErrorType = "Error"
	Default ErrorType = "Default"
)

var (
	result = Result{}
)

// ExtractFuncType extracts and returns the func returned type
func ExtractFuncType(funcType *ast.FuncType) (ErrorType, int) {

	if funcType == nil {
		return Default, -1
	}
	// spew.Dump(funcType)
	if funcType.Results != nil {
		for idx, r := range funcType.Results.List {
			if etype, ok := r.Type.(*ast.Ident); ok {
				if etype.Name == "error" {
					return Error, idx
				}
			}
		}
	}
	return Default, -1
}

func ReportSelFromExpr(expr ast.Expr, arg string) string {
	if selExpr, ok := expr.(*ast.SelectorExpr); ok {
		return fmt.Sprintf("%s(%s)", selExpr.Sel.Name, arg)
	}
	return ""
}

// ExtractArgFromExpr extracts arguments from expression
func ExtarctArgFromExpr(expr []ast.Expr) string {
	var argsConcat string
	for _, v := range expr {
		if arg, ok := v.(*ast.BasicLit); ok {
			argsConcat += arg.Value + ","
		}
	}
	return argsConcat
}

// ExtractReturnedErrorFromStmt extracts all instance of returned errors and string.
func ExtractReturnedErrorFromStmt(etypePosIdx int, body *ast.BlockStmt, funcName string) *AggregatedError {
	var errors []string
	agError := AggregatedError{
		Func: funcName,
	}
	ast.Inspect(body, func(node ast.Node) bool {
		if rtrnStmt, ok := node.(*ast.ReturnStmt); ok {
			for _, expr := range rtrnStmt.Results {
				var errorString string
				// handle call expression or wrapped errors
				if callExpr, ok := expr.(*ast.CallExpr); ok {
					errorString = ReportSelFromExpr(
						callExpr.Fun,
						ExtarctArgFromExpr(callExpr.Args),
					)
				} else {
					// handle selExpr
					errorString = ReportSelFromExpr(expr, "")
				}

				if errorString != "" {
					errors = append(errors, errorString)
				}
			}
		}
		return true
	})
	if len(errors) > 0 {
		agError.Errors = errors
		return &agError
	}
	return nil
}

// WalkThroughExpr work through the file nodes
func WalkThroughExpr(file *ast.File, fset *token.FileSet) {
	var aggregatedErrors []*AggregatedError
	for _, d := range file.Decls {
		if funcCall, ok := d.(*ast.FuncDecl); ok {
			name := funcCall.Name.Name
			returnedType, posIdx := ExtractFuncType(funcCall.Type)
			posn := fset.Position(funcCall.Pos())

			// check the returned type and position index
			if returnedType == Error && posIdx != -1 {
				agError := ExtractReturnedErrorFromStmt(posIdx, funcCall.Body, name)
				if agError != nil {
					result.AggregatedErrors = append(result.AggregatedErrors, agError)
					color.White("%s:  %s", posn, agError.Func)
					for _, v := range agError.Errors {
						color.Red("---%s ", v)
					}
				}
			}
			// ignore if func return type is not an error.
		}
	}
	result.AggregatedErrors = aggregatedErrors
}

func Run(f *ast.File, fset *token.FileSet) error {

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			// walk though expression
			WalkThroughExpr(n, fset)
		}

		return true
	})
	return nil
}
