// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package eval

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"reflect"
	"strconv"
)

var (
	ErrNilExprNode     = errors.New("nil expr node")
	ErrValNotValid     = errors.New("value not valid")
	ErrInvalidCallNode = errors.New("invalid call node")
)

type Node struct {
	Key string
	Val reflect.Value
}

type Context struct {
	This     any
	Expr     string
	ReadOnly bool
	Nodes    []*Node
}

func NewContext(this any) *Context {
	return &Context{
		This: this,
	}
}

func (c *Context) walkExprGet(this reflect.Value, astNode ast.Expr) (val reflect.Value, err error) {
	if !this.IsValid() {
		err = ErrValNotValid
		return
	}
	if astNode == nil {
		err = ErrNilExprNode
		return
	}

	var node = new(Node)
	switch expr := astNode.(type) {
	case *ast.Ident:
		err = c.evalIdentGet(this, expr, node) // A
	case *ast.IndexExpr:
		err = c.evalIndexGet(this, expr, node) // A[B]
	case *ast.CallExpr:
		err = c.evalCall(this, expr, node) // A.C()
	case *ast.SelectorExpr:
		err = c.evalSelectorGet(this, expr, node) // A.B
	default:
		return val, fmt.Errorf("unexpected expr node %T", expr)
	}
	if err == nil {
		c.Nodes = append(c.Nodes, node)
	}
	return node.Val, err
}

// 常量只能是寻址struct的field和map的key
func (c *Context) evalIdentGet(this reflect.Value, ident *ast.Ident, node *Node) error {
	if !this.IsValid() {
		return ErrValNotValid
	}
	if this.Kind() == reflect.Ptr {
		this = this.Elem()
	}
	node.Key = ident.Name
	switch this.Kind() {
	case reflect.Struct:
		node.Val = this.FieldByName(ident.Name)
	case reflect.Map:
		node.Val = this.MapIndex(reflect.ValueOf(ident.Name))
	default:
		return fmt.Errorf("unexpected kind %v with ident %s", this.Kind(), ident.Name)
	}
	return nil
}

func createMapKey(rv reflect.Value, value string) (val reflect.Value, err error) {
	s, err := strconv.Unquote(value)
	if err != nil {
		s = value
	}
	var keyType = rv.Type().Key()
	return ConvParamToType(keyType, s)
}

// 只有数组、切片、字符串和map支持`[]`操作符
func (c *Context) evalIndexGet(this reflect.Value, expr *ast.IndexExpr, node *Node) error {
	index, ok := expr.Index.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("index is not literal")
	}
	rv, err := c.walkExprGet(this, expr.X)
	if err != nil {
		return err
	}
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		i, er := strconv.Atoi(index.Value)
		if er != nil {
			return fmt.Errorf("cannot index by key %s: %w", index.Value, er)
		}
		if i >= 0 && i < rv.Len() {
			node.Key = index.Value
			node.Val = rv.Index(i)
			return nil
		} else {
			return fmt.Errorf("index out of range: %d", i)
		}
	case reflect.Map:
		key, er := createMapKey(rv, index.Value)
		if er != nil {
			return er
		}
		node.Key = index.Value
		node.Val = rv.MapIndex(key)
		return nil
	default:
		return fmt.Errorf("unexpected kind %v with ident %s", this.Kind(), expr.Index)
	}
}

// 支持单返回值和无返回值的带参函数调用
func doCallMethod(fn reflect.Value, funcName string, args []string) (val reflect.Value, err error) {
	var fnType = fn.Type()
	// 最多1个返回值
	if fnType.NumOut() > 1 {
		return val, fmt.Errorf("method %s signature not match", funcName)
	}
	input, er := ParseInputArgs(fn.Type(), args)
	if er != nil {
		err = fmt.Errorf("cannot parse input args %s: %w", funcName, er)
		return
	}
	var output = fn.Call(input)
	if len(output) > 0 {
		return output[0], nil
	}
	return
}

func tryCallMethod(this reflect.Value, call *ast.CallExpr, name string) (val reflect.Value, err error) {
	var kind = this.Kind()
	if kind == reflect.Interface {
		kind = this.Elem().Kind()
		this = this.Elem()
	}
	if kind == reflect.Ptr {
		kind = this.Elem().Kind()
	}
	if kind != reflect.Struct {
		err = fmt.Errorf("unexpected kind %v with selector %v", this.Kind(), call.Fun)
		return
	}

	var fn reflect.Value
	var isValid = false
	if this.Kind() == reflect.Ptr {
		fn = this.MethodByName(name)
		isValid = fn.IsValid()
	} else {
		// try pointer method first
		if this.CanAddr() {
			fn = this.Addr().MethodByName(name)
			isValid = fn.IsValid()
			if !isValid {
				fn = this.MethodByName(name)
				isValid = fn.IsValid()
			}
		}
	}

	if isValid {
		args, er := ParseCallExprArgs(call)
		if er != nil {
			err = fmt.Errorf("cannot parse input args %s: %w", name, er)
			return
		}
		val, err = doCallMethod(fn, name, args)
	} else {
		err = fmt.Errorf("method call %s() of %s not valid", name, this.Type().String())
	}
	return
}

func (c *Context) evalCall(this reflect.Value, call *ast.CallExpr, node *Node) error {
	if c.ReadOnly {
		return fmt.Errorf("cannot call method in read-only mode")
	}
	if !this.IsValid() {
		return ErrValNotValid
	}
	switch expr := call.Fun.(type) {
	case *ast.Ident:
		val, err := tryCallMethod(this, call, expr.Name)
		if err != nil {
			return err
		}
		node.Key = expr.Name + "()"
		node.Val = val
		return nil

	case *ast.SelectorExpr:
		obj, err := c.walkExprGet(this, expr.X)
		if err != nil {
			return err
		}
		val, err := tryCallMethod(obj, call, expr.Sel.Name)
		if err != nil {
			return err
		}
		node.Key = expr.Sel.Name + "()"
		node.Val = val
		return nil
	default:
		return fmt.Errorf("unexpect call expr %T", call.Fun)
	}
}

// 选择表达式
func (c *Context) evalSelectorGet(this reflect.Value, expr *ast.SelectorExpr, node *Node) error {
	var kind = this.Kind()
	if kind == reflect.Ptr {
		kind = this.Elem().Kind()
	}
	if kind != reflect.Struct {
		return fmt.Errorf("unexpected kind %v with ident %s", this.Kind(), expr.Sel.Name)
	}
	rv, err := c.walkExprGet(this, expr.X)
	if err != nil {
		return err
	}
	if !rv.IsValid() {
		return ErrValNotValid
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	node.Key = expr.Sel.Name
	node.Val = rv.FieldByName(expr.Sel.Name)
	return nil
}

// Eval 在`this`上，返回其`expr`对应的值
func (c *Context) Eval(expr string) (any, error) {
	c.Expr = expr
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return reflect.Value{}, err
	}
	rv, err := c.walkExprGet(reflect.ValueOf(c.This), node)
	if err != nil {
		return rv, err
	}
	if rv.IsValid() {
		return rv.Interface(), nil
	}
	return nil, ErrValNotValid
}

// set `src` to `dst`
func setValueTo(dst, src reflect.Value) error {
	if !dst.CanAddr() {
		return ErrValNotValid
	}
	srtType, dstType := src.Type(), dst.Type()
	if isPrimitiveKind(dstType.Kind()) {
		if srtType.ConvertibleTo(dstType) {
			var val = src.Convert(dstType)
			dst.Set(val)
			return nil
		}
	} else if srtType.Kind() == reflect.String {
		val, err := ConvParamToType(dstType, src.String())
		if err == nil {
			if dst.Kind() == reflect.Ptr {
				dst.Set(val)
			} else {
				dst.Set(val.Elem())
			}
			return nil
		}
	}
	return fmt.Errorf("type %s not convertible to %s", srtType.String(), dstType.String())
}

// a.X = b
func (c *Context) setIdent(lhv, rhv reflect.Value, ident *ast.Ident) error {
	var kind = lhv.Kind()
	if kind == reflect.Ptr {
		if lhv.Elem().Kind() == reflect.Struct {
			lhv = lhv.Elem()
		}
	}
	switch lhv.Kind() {
	case reflect.Struct:
		var field = lhv.FieldByName(ident.Name)
		if field.IsValid() && field.CanSet() {
			return setValueTo(field, rhv)
		}
		return ErrValNotValid

	default:
		return fmt.Errorf("unexpected kind %v with ident %s", lhv.Kind(), ident.Name)
	}
}

// 数组或者map的下标赋值，a[X] = b
func (c *Context) setIndexExpr(lhv, rhv reflect.Value, expr *ast.IndexExpr) error {
	index, ok := expr.Index.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("index is not literal")
	}
	rv, err := c.walkExprGet(lhv, expr.X)
	if err != nil {
		return err
	}

	if !rv.CanSet() {
		return ErrValNotValid
	}
	switch rv.Kind() {
	case reflect.Slice:
		idx, er := strconv.Atoi(index.Value)
		if er != nil {
			return fmt.Errorf("cannot index by key %s: %w", index.Value, er)
		}
		return setValueTo(rv.Index(idx), rhv)

	case reflect.Map:
		var valType = rv.Type().Elem()
		if !rhv.Type().ConvertibleTo(valType) {
			return fmt.Errorf("cannot convert %v to %v", rhv.Type().Kind(), valType.Kind())
		}
		if key, er := createMapKey(rv, index.Value); er != nil {
			return fmt.Errorf("cannot index map key %s: %w", index.Value, er)
		} else {
			rv.SetMapIndex(key, rhv.Convert(valType))
		}
	default:
		return fmt.Errorf("unexpected kind %v with ident %s", lhv.Kind(), expr.Index)
	}
	return nil
}

func (c *Context) setSelector(lhv, rhv reflect.Value, expr *ast.SelectorExpr) error {
	if lhv.Kind() == reflect.Ptr {
		if lhv.Elem().Kind() == reflect.Struct {
			lhv = lhv.Elem()
		}
	}
	if lhv.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected kind %v with ident %s", lhv.Kind(), expr.Sel.Name)
	}
	rv, err := c.walkExprGet(lhv, expr.X)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return ErrValNotValid
	}
	var field = rv.FieldByName(expr.Sel.Name)
	if !field.IsValid() || !field.CanAddr() {
		return ErrValNotValid
	}
	return setValueTo(field, rhv)
}

// Set 在`this`上，设置v到对应`expr`
func (c *Context) Set(expr string, v any) error {
	c.Expr = expr
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return err
	}
	var lhv = reflect.ValueOf(c.This)
	var rhv = reflect.ValueOf(v)
	switch n := node.(type) {
	case *ast.Ident:
		return c.setIdent(lhv, rhv, n)
	case *ast.IndexExpr:
		return c.setIndexExpr(lhv, rhv, n)
	case *ast.SelectorExpr:
		return c.setSelector(lhv, rhv, n)
	default:
		return fmt.Errorf("unexpected expr node %T", node)
	}
}

// 删除操作，slice, map
func (c *Context) removeIndex(this reflect.Value, expr *ast.IndexExpr) error {
	index, ok := expr.Index.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("index is not literal")
	}
	rv, err := c.walkExprGet(this, expr.X)
	if err != nil {
		return err
	}
	if !rv.CanAddr() {
		return ErrValNotValid
	}
	switch rv.Kind() {
	case reflect.Slice:
		idx, er := strconv.Atoi(index.Value)
		if er != nil {
			return fmt.Errorf("cannot index by key %s: %w", index.Value, er)
		}
		var sliceLen = rv.Len()
		if idx < 0 || idx >= sliceLen {
			return fmt.Errorf("slice index %s out of range", index.Value)
		}
		// 删除slice是通过先new一个新slice，然后把老的值赋到新slice里
		var newSlice = reflect.MakeSlice(rv.Type(), sliceLen-1, rv.Cap())
		var j = 0
		for i := 0; i < sliceLen; i++ {
			if i != idx {
				newSlice.Index(j).Set(rv.Index(i))
				j++
			}
		}
		rv.Set(newSlice)

	case reflect.Map:
		if key, er := createMapKey(rv, index.Value); er != nil {
			return fmt.Errorf("cannot index map key %s: %w", index.Value, er)
		} else {
			rv.SetMapIndex(key, reflect.Value{}) // do deletion
		}
	default:
		return fmt.Errorf("unexpected kind %v with ident %s", this.Kind(), expr.Index)
	}
	return nil
}

// Delete 删除 a[X]
func (c *Context) Delete(expr string) error {
	c.Expr = expr
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return err
	}
	if node == nil {
		return ErrNilExprNode
	}
	var rv = reflect.ValueOf(c.This)
	if !rv.IsValid() {
		return ErrValNotValid
	}
	switch n := node.(type) {
	case *ast.IndexExpr:
		return c.removeIndex(rv, n)
	default:
		return fmt.Errorf("unexpected expr node %T", node)
	}
}

func (c *Context) PrintNodes(w io.Writer) {
	for _, node := range c.Nodes {
		fmt.Fprintf(w, "%s -> %s\n", node.Key, node.Val.Type().String())
	}
	w.Write([]byte("\n"))
}
