// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor, err error)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Visitor, list []*Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmtList(v Visitor, list []Stmt) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkDeclList(v Visitor, list []Decl) {
	for _, x := range list {
		Walk(v, x)
	}
}

// TODO(gri): Investigate if providing a closure to Walk leads to
//            simpler use (and may help eliminate Inspect in turn).

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) error {
	var err error
	v, err = v.Visit(node)
	if err != nil {
		return err
	}
	if v == nil {
		return nil
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	// case *Comment:
	// nothing to do

	// case *CommentGroup:
	// 	for _, c := range n.List {
	// 		Walk(v, c)
	// 	}

		// TODO: обработка ошибок!!!
	case *Field:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		walkIdentList(v, n.Names)
		err=Walk(v, n.Type)
		if err!=nil{
			return err
		}
		if n.Tag != nil {
			err=Walk(v, n.Tag)
		}
		// if n.Comment != nil {
		// 	Walk(v, n.Comment)
		// }

	case *FieldList:
		for _, f := range n.List {
			err=Walk(v, f)
			if err!=nil{
				return err
			}
		
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	// case *Ellipsis:
	// 	if n.Elt != nil {
	// 		Walk(v, n.Elt)
	// 	}

	case *FuncLit:
		err=Walk(v, n.Type)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Body)

	case *CompositeLit:
		// if n.Type != nil {
		// 	Walk(v, n.Type)
		// }
		walkExprList(v, n.Elts)

	case *ParenExpr:
		err=Walk(v, n.X)

	case *TernaryExpr:
		err=Walk(v, n.Cond)
		if err!=nil{
			return err
		}
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Y)

	case *SelectorExpr:
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Sel)

	case *IndexExpr:
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Index)

	case *SliceExpr:
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		if n.Low != nil {
			err=Walk(v, n.Low)
		}
		if err!=nil{
			return err
		}
		if n.High != nil {
			err=Walk(v, n.High)
		}
		// if n.Max != nil {
		// 	Walk(v, n.Max)
		// }

	case *TypeAssertExpr:
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		if n.Type != nil {
			err=Walk(v, n.Type)
		}

	case *CallExpr:
		err=Walk(v, n.Fun)
		if err!=nil{
			return err
		}
		walkExprList(v, n.Args)

	// case *StarExpr:
	// 	Walk(v, n.X)

	case *UnaryExpr:
		err=Walk(v, n.X)

	case *BinaryExpr:
		err=Walk(v, n.X)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Y)

	case *KeyValueExpr:
		err=Walk(v, n.Key)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Value)

	// Types
	// case *ArrayType:
	// 	if n.Len != nil {
	// 		Walk(v, n.Len)
	// 	}
	// 	Walk(v, n.Elt)

	// case *StructType:
	// 	Walk(v, n.Fields)

	case *FuncType:
		if n.Params != nil {
			err=Walk(v, n.Params)
			if err!=nil{
				return err
			}
		}
		// if n.Results != nil {
		// 	Walk(v, n.Results)
		// }

	// case *InterfaceType:
	// 	Walk(v, n.Methods)

	// case *MapType:
	// 	Walk(v, n.Key)
	// 	Walk(v, n.Value)

	// case *ChanType:
	// 	Walk(v, n.Value)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		err=Walk(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *LabeledStmt:
		err=Walk(v, n.Label)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Stmt)

	case *ExprStmt:
		err=Walk(v, n.X)

	// case *SendStmt:
	// 	Walk(v, n.Chan)
	// 	Walk(v, n.Value)

	case *IncDecStmt:
		err=Walk(v, n.X)

	case *AssignStmt:
		walkExprList(v, n.Lhs)
		walkExprList(v, n.Rhs)

	case *GoStmt:
		err=Walk(v, n.Call)

	case *DeferStmt:
		err=Walk(v, n.Call)

	case *ReturnStmt:
		walkExprList(v, n.Results)

	case *BranchStmt:
		if n.Label != nil {
			err=Walk(v, n.Label)
		}

	case *BlockStmt:
		walkStmtList(v, n.List)

	case *IfStmt:
		// if n.Init != nil {
		// 	Walk(v, n.Init)
		// }
		err=Walk(v, n.Cond)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Body)
		if err!=nil{
			return err
		}
		if n.ElsIf != nil {
			walkStmtList(v, n.ElsIf)
		}
		if err!=nil{
			return err
		}
		if n.Else != nil {
			err=Walk(v, n.Else)
		}

	case *TryStmt:
		err=Walk(v, n.Body)
		if err!=nil{
			return err
		}
		if n.Except != nil {
			err=Walk(v, n.Except)
		}

	// case *CaseClause:
	// 	walkExprList(v, n.List)
	// 	walkStmtList(v, n.Body)

	// case *SwitchStmt:
	// 	if n.Init != nil {
	// 		Walk(v, n.Init)
	// 	}
	// 	if n.Tag != nil {
	// 		Walk(v, n.Tag)
	// 	}
	// 	Walk(v, n.Body)

	// case *TypeSwitchStmt:
	// 	if n.Init != nil {
	// 		Walk(v, n.Init)
	// 	}
	// 	Walk(v, n.Assign)
	// 	Walk(v, n.Body)

	// case *CommClause:
	// 	if n.Comm != nil {
	// 		Walk(v, n.Comm)
	// 	}
	// 	walkStmtList(v, n.Body)

	// case *SelectStmt:
	// 	Walk(v, n.Body)

	case *ForStmt:
		// if n.Init != nil {
		// 	Walk(v, n.Init)
		// }
		if n.Cond != nil {
			err=Walk(v, n.Cond)
		}
		// if n.Post != nil {
		// 	Walk(v, n.Post)
		// }
		if err!=nil{
			return err
		}
		err=Walk(v, n.Body)

	case *WhileStmt:
		if n.Cond != nil {
			err=Walk(v, n.Cond)
		}
		if err!=nil{
			return err
		}
		err=Walk(v, n.Body)

	// case *RangeStmt:
	// 	if n.Key != nil {
	// 		Walk(v, n.Key)
	// 	}
	// 	if n.Value != nil {
	// 		Walk(v, n.Value)
	// 	}
	// 	Walk(v, n.X)
	// 	Walk(v, n.Body)

	// Declarations
	case *ImportSpec:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		if n.Name != nil {
			err=Walk(v, n.Name)
		}
		if err!=nil{
			return err
		}
		err=Walk(v, n.Path)
		// if n.Comment != nil {
		// 	Walk(v, n.Comment)
		// }

	case *ValueSpec:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		walkIdentList(v, n.Names)
		if n.Type != nil {
			err=Walk(v, n.Type)
		}
		if err!=nil{
			return err
		}
		walkExprList(v, n.Values)
		// if n.Comment != nil {
		// 	Walk(v, n.Comment)
		// }

	// case *TypeSpec:
	// 	if n.Doc != nil {
	// 		Walk(v, n.Doc)
	// 	}
	// 	Walk(v, n.Name)
	// 	Walk(v, n.Type)
	// 	if n.Comment != nil {
	// 		Walk(v, n.Comment)
	// 	}

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		for _, s := range n.Specs {
			err=Walk(v, s)
			if err!=nil{
				return err
			}
		}

	case *FuncDecl:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		// if n.Recv != nil {
		// 	Walk(v, n.Recv)
		// }
		err=Walk(v, n.Name)
		if err!=nil{
			return err
		}
		err=Walk(v, n.Type)
		if err!=nil{
			return err
		}
		if n.Body != nil {
			err=Walk(v, n.Body)
		}

	// Files and packages
	case *File:
		// if n.Doc != nil {
		// 	Walk(v, n.Doc)
		// }
		err=Walk(v, n.Name)
		if err!=nil{
			return err
		}
		walkDeclList(v, n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for _, f := range n.Files {
			err=Walk(v, f)
			if err!=nil{
				return err
			}
		}

	default:
		return fmt.Errorf("ast.Walk: unexpected node type %T", n)
	}
	 if err!=nil{
		 return err
	 }
	_,err=v.Visit(nil)
	return err
}

type inspector func(Node) (bool, error)

func (f inspector) Visit(node Node) (Visitor, error) {
	res, err := f(node)
	if err != nil {
		return nil, err
	}
	if res {
		return f, nil
	}
	return nil, nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node Node, f func(Node) (bool, error)) error {
	return Walk(inspector(f), node)
}