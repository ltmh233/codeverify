package analysis

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type FunctionInfo struct {
	Package   string // Package: 包名
	PkgPath   string // PkgPath: 包导入路径
	Symbol    string // Symbol: 唯一标志
	Name      string // Name: 函数名
	Receiver  string // Receiver: 方法接收者(成员函数)
	Signature string // Signature: 函数签名
	File      string // File: 相对仓库路径
	Line      int    // Line: 声明起始行
}

func getSignature(fset *token.FileSet, fn *ast.FuncDecl) (string, error) {
	var buf bytes.Buffer

	err := format.Node(&buf, fset, fn.Type)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getReceiver(fset *token.FileSet, fn *ast.FuncDecl) (string, error) {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", nil
	}

	var buf bytes.Buffer

	err := format.Node(&buf, fset, fn.Recv.List[0].Type)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SymbolID(fn *types.Func) string {
	return fn.Pkg().Name() + "::" + fn.FullName()
}

func ScanGoFunctions(repoAbs string) ([]FunctionInfo, *CallGraph, error) {

	functionInfos := []FunctionInfo{}

	graph := CallGraph{
		Edges: []CallEdge{},
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  repoAbs,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, nil, err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				pos := pkg.Fset.Position(fn.Pos())
				receiver, err := getReceiver(pkg.Fset, fn)
				if err != nil {
					return nil, nil, err
				}
				signature, err := getSignature(pkg.Fset, fn)
				if err != nil {
					return nil, nil, err
				}
				relativePath, err := filepath.Rel(repoAbs, pos.Filename)
				if err != nil {
					return nil, nil, err
				}

				obj := pkg.TypesInfo.Defs[fn.Name]

				def, ok := obj.(*types.Func)
				if !ok {
					continue
				}

				functionInfo := FunctionInfo{
					Name:      fn.Name.Name,
					Receiver:  receiver,
					File:      relativePath,
					Line:      pos.Line,
					Signature: signature,
					Symbol:    SymbolID(def),
					Package:   pkg.Name,
					PkgPath:   pkg.PkgPath,
				}

				functionInfos = append(functionInfos, functionInfo)

				ScanCalls(&graph, pkg, fn, def)
			}
		}
	}
	return functionInfos, &graph, nil
}

func ScanCalls(graph *CallGraph, pkg *packages.Package, fn *ast.FuncDecl, caller *types.Func) {
	if fn.Body == nil {
		return
	}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var ident *ast.Ident

		switch fun := call.Fun.(type) {
		case *ast.Ident:
			ident = fun
		case *ast.SelectorExpr:
			ident = fun.Sel
		default:
			return true
		}

		obj := pkg.TypesInfo.Uses[ident]

		callee, ok := obj.(*types.Func)
		if !ok {
			return true
		}

		pos := pkg.Fset.Position(callee.Pos())

		graph.Edges = append(graph.Edges, CallEdge{
			Caller:   SymbolID(caller),
			Callee:   SymbolID(callee),
			CallFile: pos.Filename,
			CallLine: pos.Line,
		})

		return true

	})
}
