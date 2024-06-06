package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

type AstHelper struct {
	posArr  []token.Pos
	nodeMap map[token.Pos][]ast.Node
	tfSet   *token.FileSet
	tf      *token.File
	f       *ast.File
}

func NewAstHelper(filePath string) (*AstHelper, error) {
	tfSet := token.NewFileSet()
	f, err := parser.ParseFile(tfSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	posArr := []token.Pos{}
	nodeMap := map[token.Pos][]ast.Node{}
	ast.Inspect(f, func(n ast.Node) bool {
		if n != nil && n.Pos() != n.End() {
			pos := n.Pos()
			posArr = append(posArr, pos)
			nodeMap[pos] = append(nodeMap[pos], n)
		}
		return true
	})

	sort.Slice(posArr, func(i, j int) bool {
		return posArr[i] < posArr[j]
	})

	return &AstHelper{
		posArr:  posArr,
		nodeMap: nodeMap,
		tfSet:   tfSet,
		tf:      tfSet.File(f.Pos()),
		f:       f,
	}, nil
}

func (h *AstHelper) GetNode(p token.Pos) []ast.Node {
	return h.nodeMap[p]
}

func (h *AstHelper) Iter(fn func(int, ast.Node) bool) {
	for i, p := range h.posArr {
		for _, n := range h.nodeMap[p] {
			if !fn(i, n) {
				return
			}
		}
	}
}

func (h *AstHelper) PrintAst() {
	_ = ast.Print(h.tfSet, h.f)
}

func GetNodeType(n any) string {
	// Find Function Call Statements
	switch n := n.(type) {
	case *ast.Comment:
		println("Comment", n.Pos(), "~", n.End())
		return "Comment"
	case *ast.CommentGroup:
		println("CommentGroup", n.Pos(), "~", n.End())
		return "CommentGroup"

		// Expressions
	case *ast.Field:
		println("Field", n.Pos(), "~", n.End())
		return "Field"
	case *ast.FieldList:
		println("FieldList", n.Pos(), "~", n.End())
		return "FieldList"
	case *ast.BadExpr:
		println("BadExpr", n.Pos(), "~", n.End())
		return "BadExpr"
	case *ast.Ident:
		println("Ident", n.Pos(), "~", n.End())
		return "Ident"
	case *ast.Ellipsis:
		println("Ellipsis", n.Pos(), "~", n.End())
		return "Ellipsis"
	case *ast.BasicLit:
		println("BasicLit", n.Pos(), "~", n.End())
		return "BasicLit"
	case *ast.FuncLit:
		println("FuncLit", n.Pos(), "~", n.End())
		return "FuncLit"
	case *ast.CompositeLit:
		println("CompositeLit", n.Pos(), "~", n.End())
		return "CompositeLit"

		// Expression
	case *ast.ParenExpr:
		println("ParenExpr", n.Pos(), "~", n.End())
		return "ParenExpr"
	case *ast.SelectorExpr:
		println("SelectorExpr", n.Pos(), "~", n.End())
		return "SelectorExpr"
	case *ast.IndexExpr:
		println("IndexExpr", n.Pos(), "~", n.End())
		return "IndexExpr"
	case *ast.IndexListExpr:
		println("IndexListExpr", n.Pos(), "~", n.End())
		return "IndexListExpr"
	case *ast.SliceExpr:
		println("SliceExpr", n.Pos(), "~", n.End())
		return "SliceExpr"
	case *ast.TypeAssertExpr:
		println("TypeAssertExpr", n.Pos(), "~", n.End())
		return "TypeAssertExpr"
	case *ast.CallExpr:
		println("CallExpr", n.Pos(), "~", n.End())
		return "CallExpr"
	case *ast.StarExpr:
		println("StarExpr", n.Pos(), "~", n.End())
		return "StarExpr"
	case *ast.UnaryExpr:
		println("UnaryExpr", n.Pos(), "~", n.End())
		return "UnaryExpr"
	case *ast.BinaryExpr:
		println("BinaryExpr", n.Pos(), "~", n.End())
		return "BinaryExpr"
	case *ast.KeyValueExpr:
		println("KeyValueExpr", n.Pos(), "~", n.End())
		return "KeyValueExpr"

		// TYPE
	case *ast.ArrayType:
		println("ArrayType", n.Pos(), "~", n.End())
		return "ArrayType"
	case *ast.StructType:
		println("StructType", n.Pos(), "~", n.End())
		return "StructType"
	case *ast.FuncType:
		println("FuncType", n.Pos(), "~", n.End())
		return "FuncType"
	case *ast.InterfaceType:
		println("InterfaceType", n.Pos(), "~", n.End())
		return "InterfaceType"
	case *ast.MapType:
		println("MapType", n.Pos(), "~", n.End())
		return "MapType"
	case *ast.ChanType:
		println("ChanType", n.Pos(), "~", n.End())
		return "ChanType"

		// Statements
	case *ast.BadStmt:
		println("BadStmt", n.Pos(), "~", n.End())
		return "BadStmt"
	case *ast.DeclStmt:
		println("DeclStmt", n.Pos(), "~", n.End())
		return "DeclStmt"
	case *ast.EmptyStmt:
		println("EmptyStmt", n.Pos(), "~", n.End())
		return "EmptyStmt"
	case *ast.LabeledStmt:
		println("LabeledStmt", n.Pos(), "~", n.End())
		return "LabeledStmt"
	case *ast.ExprStmt:
		println("ExprStmt", n.Pos(), "~", n.End())
		return "ExprStmt"
	case *ast.SendStmt:
		println("SendStmt", n.Pos(), "~", n.End())
		return "SendStmt"
	case *ast.IncDecStmt:
		println("IncDecStmt", n.Pos(), "~", n.End())
		return "IncDecStmt"
	case *ast.AssignStmt:
		println("AssignStmt", n.Pos(), "~", n.End())
		return "AssignStmt"
	case *ast.GoStmt:
		println("GoStmt", n.Pos(), "~", n.End())
		return "GoStmt"
	case *ast.DeferStmt:
		println("DeferStmt", n.Pos(), "~", n.End())
		return "DeferStmt"
	case *ast.ReturnStmt:
		println("ReturnStmt", n.Pos(), "~", n.End())
		return "ReturnStmt"
	case *ast.BranchStmt:
		println("BranchStmt", n.Pos(), "~", n.End())
		return "BranchStmt"
	case *ast.BlockStmt:
		println("BlockStmt", n.Pos(), "~", n.End())
		return "BlockStmt"
	case *ast.IfStmt:
		println("IfStmt", n.Pos(), "~", n.End())
		return "IfStmt"
	case *ast.CaseClause:
		println("CaseClause", n.Pos(), "~", n.End())
		return "CaseClause"
	case *ast.SwitchStmt:
		println("SwitchStmt", n.Pos(), "~", n.End())
		return "SwitchStmt"
	case *ast.TypeSwitchStmt:
		println("TypeSwitchStmt", n.Pos(), "~", n.End())
		return "TypeSwitchStmt"
	case *ast.CommClause:
		println("CommClause", n.Pos(), "~", n.End())
		return "CommClause"
	case *ast.SelectStmt:
		println("SelectStmt", n.Pos(), "~", n.End())
		return "SelectStmt"
	case *ast.ForStmt:
		println("ForStmt", n.Pos(), "~", n.End())
		return "ForStmt"
	case *ast.RangeStmt:
		println("RangeStmt", n.Pos(), "~", n.End())
		return "RangeStmt"

		// Declarations
	case *ast.ImportSpec:
		println("ImportSpec", n.Pos(), "~", n.End())
		return "ImportSpec"
	case *ast.ValueSpec:
		println("ValueSpec", n.Pos(), "~", n.End())
		return "ValueSpec"
	case *ast.TypeSpec:
		println("TypeSpec", n.Pos(), "~", n.End())
		return "TypeSpec"
	case *ast.BadDecl:
		println("BadDecl", n.Pos(), "~", n.End())
		return "BadDecl"
	case *ast.GenDecl:
		println("GenDecl", n.Pos(), "~", n.End())
		return "GenDecl"
	case *ast.FuncDecl:
		println("FuncDecl", n.Pos(), "~", n.End())
		return "FuncDecl"
		// File
	case *ast.File:
		println("File", n.Pos(), "~", n.End())
		return "File"
	case *ast.Package:
		println("Package", n.Pos(), "~", n.End())
		return "Package"

		// Interface
	case ast.Expr:
		println("Expr", n.Pos(), "~", n.End())
		return "Expr"
	case ast.Stmt:
		println("Stmt", n.Pos(), "~", n.End())
		return "Stmt"
	case ast.Decl:
		println("Decl", n.Pos(), "~", n.End())
		return "Decl"
	case ast.Node:
		println("Node", n.Pos(), "~", n.End())
		return "Node"
	case nil:
		println("nil")
		return "nil"
	default:
		return "【】"
	}
}
