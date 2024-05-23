// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modifications Copyright (c) qawatake 2023

package fsel

import (
	"go/token"
	"go/types"
	"maps"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const name = "fsel"
const doc = "flags field access with unverified nil errors."
const url = "https://pkg.go.dev/github.com/qawatake/fsel"

var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	URL:  url,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	for _, fn := range funcs {
		runFunc(pass, fn)
	}

	return nil, nil
}

func returnsPtrsAndErr(pass *analysis.Pass, id int, instr ssa.Instruction) []*ptrerr {
	defer func() { recover() }()
	call := instr.(*ssa.Call)
	numResults := call.Common().Signature().Results().Len()
	// -1 to omit error
	ptrerrs := make([]*ptrerr, 0, numResults-1)
	var errValue ssa.Value
	for j := id; j < len(instr.Block().Instrs); j++ {
		func() {
			defer func() { recover() }()
			extract := instr.Block().Instrs[j].(*ssa.Extract)
			if extract.Tuple == call {
				if types.Identical(extract.Type(), typesErr) {
					errValue = extract
					return
				}
				if extract.Type().Underlying().(*types.Pointer) != nil {
					ptrerrs = append(ptrerrs, &ptrerr{
						ptr: extract,
					})
					return
				}
			}
		}()
	}
	if len(ptrerrs) != 0 && errValue != nil {
		for _, p := range ptrerrs {
			p.err = errValue
		}
		return ptrerrs
	}
	return nil
}

type ptrerr struct {
	ptr ssa.Value
	err ssa.Value
}

var typesErr = types.Universe.Lookup("error").Type()

func runFunc(pass *analysis.Pass, fn *ssa.Function) {
	// visit visits reachable blocks of the CFG in dominance order,
	// maintaining a stack of dominating nilness facts.
	//
	// By traversing the dom tree, we can pop facts off the stack as
	// soon as we've visited a subtree.  Had we traversed the CFG,
	// we would need to retain the set of facts for each block.
	seen := make([]bool, len(fn.Blocks)) // seen[i] means visit should ignore block i
	// allocated: 同じポインタt0に複数の値が代入される場合、最後の代入値を記録する。
	// *t0のnilnessは最後の代入値のnilnessにもなる。
	// 例はtestdata/src/a/a.goのf6関数
	var visit func(b *ssa.BasicBlock, stack []fact, allocated map[*ssa.Alloc]ssa.Value)
	ptrerrs := make([]*ptrerr, 0, 10)
	visit = func(b *ssa.BasicBlock, stack []fact, latestAllocated map[*ssa.Alloc]ssa.Value) {
		if seen[b.Index] {
			return
		}
		seen[b.Index] = true

		// Report nil dereferences.
		for i, instr := range b.Instrs {
			if p := returnsPtrsAndErr(pass, i, instr); len(p) > 0 {
				ptrerrs = append(ptrerrs, p...)
			}
			switch instr := instr.(type) {
			case *ssa.FieldAddr:
				for _, ptrerr := range ptrerrs {
					if instr.X == ptrerr.ptr {
						if nilnessOf(stack, ptrerr.err) != isnil && nilnessOf(stack, ptrerr.ptr) != isnonnil {
							pass.Reportf(instr.Pos(), "field address without checking nilness of err")
						}
					} else if alloc := refOfAllocated(instr.X); alloc != nil {
						if latestAllocated[alloc] == ptrerr.ptr {
							if nilnessOf(stack, ptrerr.err) != isnil && nilnessOf(stack, ptrerr.ptr) != isnonnil {
								pass.Reportf(instr.Pos(), "field address without checking nilness of err")
							}
						}
					}
				}
			case *ssa.Store:
				if alloc, ok := instr.Addr.(*ssa.Alloc); ok {
					latestAllocated[alloc] = instr.Val
				}
			}
		}

		// For nil comparison blocks, report an error if the condition
		// is degenerate, and push a nilness fact on the stack when
		// visiting its true and false successor blocks.
		if binop, tsucc, fsucc := eq(b); binop != nil {
			xnil := nilnessOf(stack, binop.X)
			ynil := nilnessOf(stack, binop.Y)

			if ynil != unknown && xnil != unknown && (xnil == isnil || ynil == isnil) {
				// If tsucc's or fsucc's sole incoming edge is impossible,
				// it is unreachable.  Prune traversal of it and
				// all the blocks it dominates.
				// (We could be more precise with full dataflow
				// analysis of control-flow joins.)
				var skip *ssa.BasicBlock
				if xnil == ynil {
					skip = fsucc
				} else {
					skip = tsucc
				}
				for _, d := range b.Dominees() {
					if d == skip && len(d.Preds) == 1 {
						continue
					}
					visit(d, stack, maps.Clone(latestAllocated))
				}
				return
			}

			// "if x == nil" or "if nil == y" condition; x, y are unknown.
			if xnil == isnil || ynil == isnil {
				var newFacts facts
				if xnil == isnil {
					// x is nil, y is unknown:
					// t successor learns y is nil.
					newFacts = expandFacts(fact{binop.Y, isnil})
					if alloc := refOfAllocated(binop.Y); alloc != nil {
						if val, ok := latestAllocated[alloc]; ok {
							newFacts = append(newFacts, expandFacts(fact{val, isnil})...)
						}
					}
				} else {
					// x is nil, y is unknown:
					// t successor learns x is nil.
					newFacts = expandFacts(fact{binop.X, isnil})
					if alloc := refOfAllocated(binop.X); alloc != nil {
						if val, ok := latestAllocated[alloc]; ok {
							newFacts = append(newFacts, expandFacts(fact{val, isnil})...)
						}
					}
				}

				for _, d := range b.Dominees() {
					// Successor blocks learn a fact
					// only at non-critical edges.
					// (We could do be more precise with full dataflow
					// analysis of control-flow joins.)
					s := stack
					if len(d.Preds) == 1 {
						if d == tsucc {
							s = append(s, newFacts...)
						} else if d == fsucc {
							s = append(s, newFacts.negate()...)
						}
					}
					allocated := maps.Clone(latestAllocated)
					visit(d, s, allocated)
				}
				return
			}
		}

		for _, d := range b.Dominees() {
			allocated := maps.Clone(latestAllocated)
			visit(d, stack, allocated)
		}
	}

	// Visit the entry block.  No need to visit fn.Recover.
	if fn.Blocks != nil {
		visit(fn.Blocks[0], make([]fact, 0, 20), make(map[*ssa.Alloc]ssa.Value)) // 20 is plenty
	}
}

// *t0 (t0 is a *ssa.Alloc) -> t0
// otherwise returns nil
func refOfAllocated(v ssa.Value) *ssa.Alloc {
	if unop, ok := v.(*ssa.UnOp); ok {
		if unop.Op == token.MUL {
			if alloc, ok := unop.X.(*ssa.Alloc); ok {
				return alloc
			}
		}
	}
	return nil
}

// A fact records that a block is dominated
// by the condition v == nil or v != nil.
type fact struct {
	value   ssa.Value
	nilness nilness
}

func (f fact) negate() fact { return fact{f.value, -f.nilness} }

type nilness int

const (
	isnonnil         = -1
	unknown  nilness = 0
	isnil            = 1
)

var nilnessStrings = []string{"non-nil", "unknown", "nil"}

func (n nilness) String() string { return nilnessStrings[n+1] }

// nilnessOf reports whether v is definitely nil, definitely not nil,
// or unknown given the dominating stack of facts.
func nilnessOf(stack []fact, v ssa.Value) nilness {
	switch v := v.(type) {
	// unwrap ChangeInterface and Slice values recursively, to detect if underlying
	// values have any facts recorded or are otherwise known with regard to nilness.
	//
	// This work must be in addition to expanding facts about
	// ChangeInterfaces during inference/fact gathering because this covers
	// cases where the nilness of a value is intrinsic, rather than based
	// on inferred facts, such as a zero value interface variable. That
	// said, this work alone would only inform us when facts are about
	// underlying values, rather than outer values, when the analysis is
	// transitive in both directions.
	case *ssa.ChangeInterface:
		if underlying := nilnessOf(stack, v.X); underlying != unknown {
			return underlying
		}
	case *ssa.Slice:
		if underlying := nilnessOf(stack, v.X); underlying != unknown {
			return underlying
		}
	case *ssa.SliceToArrayPointer:
		nn := nilnessOf(stack, v.X)
		if slice2ArrayPtrLen(v) > 0 {
			if nn == isnil {
				// We know that *(*[1]byte)(nil) is going to panic because of the
				// conversion. So return unknown to the caller, prevent useless
				// nil deference reporting due to * operator.
				return unknown
			}
			// Otherwise, the conversion will yield a non-nil pointer to array.
			// Note that the instruction can still panic if array length greater
			// than slice length. If the value is used by another instruction,
			// that instruction can assume the panic did not happen when that
			// instruction is reached.
			return isnonnil
		}
		// In case array length is zero, the conversion result depends on nilness of the slice.
		if nn != unknown {
			return nn
		}
	}

	// Is value intrinsically nil or non-nil?
	switch v := v.(type) {
	case *ssa.Alloc,
		*ssa.FieldAddr,
		*ssa.FreeVar,
		*ssa.Function,
		*ssa.Global,
		*ssa.IndexAddr,
		*ssa.MakeChan,
		*ssa.MakeClosure,
		*ssa.MakeInterface,
		*ssa.MakeMap,
		*ssa.MakeSlice:
		return isnonnil
	case *ssa.Const:
		if v.IsNil() {
			return isnil // nil or zero value of a pointer-like type
		} else {
			return unknown // non-pointer
		}
	}

	// Search dominating control-flow facts.
	for _, f := range stack {
		if f.value == v {
			return f.nilness
		}
	}
	return unknown
}

func slice2ArrayPtrLen(v *ssa.SliceToArrayPointer) int64 {
	return v.Type().(*types.Pointer).Elem().Underlying().(*types.Array).Len()
}

// If b ends with an equality comparison, eq returns the operation and
// its true (equal) and false (not equal) successors.
func eq(b *ssa.BasicBlock) (op *ssa.BinOp, tsucc, fsucc *ssa.BasicBlock) {
	if If, ok := b.Instrs[len(b.Instrs)-1].(*ssa.If); ok {
		if binop, ok := If.Cond.(*ssa.BinOp); ok {
			switch binop.Op {
			case token.EQL:
				return binop, b.Succs[0], b.Succs[1]
			case token.NEQ:
				return binop, b.Succs[1], b.Succs[0]
			}
		}
	}
	return nil, nil, nil
}

// expandFacts takes a single fact and returns the set of facts that can be
// known about it or any of its related values. Some operations, like
// ChangeInterface, have transitive nilness, such that if you know the
// underlying value is nil, you also know the value itself is nil, and vice
// versa. This operation allows callers to match on any of the related values
// in analyses, rather than just the one form of the value that happened to
// appear in a comparison.
//
// This work must be in addition to unwrapping values within nilnessOf because
// while this work helps give facts about transitively known values based on
// inferred facts, the recursive check within nilnessOf covers cases where
// nilness facts are intrinsic to the underlying value, such as a zero value
// interface variables.
//
// ChangeInterface is the only expansion currently supported, but others, like
// Slice, could be added. At this time, this tool does not check slice
// operations in a way this expansion could help. See
// https://play.golang.org/p/mGqXEp7w4fR for an example.
func expandFacts(f fact) []fact {
	ff := []fact{f}

Loop:
	for {
		switch v := f.value.(type) {
		case *ssa.ChangeInterface:
			f = fact{v.X, f.nilness}
			ff = append(ff, f)
		default:
			break Loop
		}
	}

	return ff
}

type facts []fact

func (ff facts) negate() facts {
	nn := make([]fact, len(ff))
	for i, f := range ff {
		nn[i] = f.negate()
	}
	return nn
}
