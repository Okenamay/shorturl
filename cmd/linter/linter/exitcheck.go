package linter

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

// ExitCheckAnalyzer - это анализатор, который проверяет код на
// использование panic, os.Exit и log.Fatal.
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "checks for calls to panic, os.Exit, and log.Fatal",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// isMainFuncInMainPkg проверяет, находимся ли мы в функции main пакета main
	isMainFuncInMainPkg := func(n ast.Node, file *ast.File) bool {
		if pass.Pkg.Name() != "main" {
			return false
		}

		// Ищем ближайшую родительскую функцию
		path, _ := astutil.PathEnclosingInterval(file, n.Pos(), n.End())
		if path == nil {
			return false
		}

		for _, ancestor := range path {
			if fd, ok := ancestor.(*ast.FuncDecl); ok {
				// Нашли. Это main?
				return fd.Name.Name == "main"
			}
		}
		return false
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем вызов функции
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// 1. Проверка на panic()
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				// Проверяем, что это встроенный panic, а не кастомная функция
				if obj, ok := pass.TypesInfo.Uses[ident]; ok {
					// Проверяем, что это действительно встроенная (Builtin) функция
					if builtin, ok := obj.(*types.Builtin); ok && builtin.Name() == "panic" {
						pass.Reportf(call.Pos(), "do not use panic")
					}
				}
				return true
			}

			// 2. Проверка на os.Exit() и log.Fatal()
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Ищем sel.X (имя пакета)
			pkgIdent, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			// Имя функции (sel.Sel.Name)
			funcName := sel.Sel.Name

			isOsExit := pkgIdent.Name == "os" && funcName == "Exit"
			isLogFatal := pkgIdent.Name == "log" && (funcName == "Fatal" || funcName == "Fatalf")

			if isOsExit || isLogFatal {
				// Правило 2: Разрешено использовать только в main.main
				if !isMainFuncInMainPkg(n, file) {
					pass.Reportf(call.Pos(), "os.Exit/log.Fatal call outside main.main")
				}
			}

			return true
		})
	}
	return nil, nil
}
