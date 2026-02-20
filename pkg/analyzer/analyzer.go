package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
	"unicode"
)

const Doc = "Check log messages for rules compliance"

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			funcName := getFuncName(call)
			if !isLogFunction(funcName) {
				return true
			}
			for _, arg := range call.Args {
				if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind.String() == "STRING" {
					message := strings.Trim(lit.Value, "\"")
					checkLogMessage(pass, lit, message)
				}
			}
			return true
		})
	}
	return nil, nil
}

func getFuncName(call *ast.CallExpr) string {
	switch f := call.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := f.X.(*ast.Ident); ok {
			return ident.Name + "." + f.Sel.Name
		}
	case *ast.Ident:
		return f.Name
	}
	return ""
}

func isLogFunction(name string) bool {
	logFunctions := []string{
		"log.Info", "log.Debug", "log.Warn", "log.Error",
		"slog.Info", "slog.Debug", "slog.Warn", "slog.Error",
	}
	for _, f := range logFunctions {
		if strings.HasPrefix(name, f) {
			return true
		}
	}
	return false
}

func checkLogMessage(pass *analysis.Pass, node ast.Node, message string) {
	if message == "" {
		return
	}
	tmp := strings.TrimSpace(message)
	if len(tmp) > 0 {
		firstRune := []rune(tmp)[0]
		if unicode.IsUpper(firstRune) {
			pass.Reportf(node.Pos(), "log message should start with lowercase letter")
		}
	}
	for _, r := range tmp {
		if unicode.IsLetter(r) && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			pass.Reportf(node.Pos(), "use English only in log messages")
			break
		}
	}
	allowedSymbols := " .,:;-_()/=+"
	for _, r := range tmp {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune(allowedSymbols, r) {
			pass.Reportf(node.Pos(), "no special chars or emoji allowed")
			break
		}
	}
	if len(tmp) > 1000 {
		return
	}
	sensitiveWords := []string{"password", "token", "key", "secret", "api_key", "apikey"}
	lower := strings.ToLower(tmp)
	for _, word := range sensitiveWords {
		if strings.Contains(lower, word) {
			pass.Reportf(node.Pos(), "sensitive data detected")
			break
		}
	}
}
