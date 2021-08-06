// Copyright 2021 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package solver

import (
	"fmt"
	"strings"

	"github.com/irfansharif/solver/internal/pb"
)

// Constraint is what a model attempts to satisfy when deciding on a solution.
type Constraint interface {
	// OnlyEnforceIf enforces the constraint iff all literals listed are true. If
	// not explicitly called, of if the list is empty, then the constraint will
	// always be enforced.
	//
	// NB: Only a few constraints support enforcement:
	// - NewBooleanOrConstraint
	// - NewBooleanAndConstraint
	// - NewLinearConstraint
	//
	// Intervals support enforcement too, but only with a single literal.
	OnlyEnforceIf(literals ...Literal) Constraint

	// Stringer provides a printable format representation for the constraint.
	fmt.Stringer

	// XXX: Constraints can also be named, how do we support that? Would be nice
	// to have a unified way to name things (model, intvar, literal, intervals,
	// constraints).

	// protos returns the underlying CP-SAT constraint protobuf representations.
	protos() []*pb.ConstraintProto
}

type constraint struct {
	pb          *pb.ConstraintProto
	enforcement []Literal
	str         string
}

// String is part of the Constraint interface.
func (c *constraint) String() string {
	if c.str == "" {
		return fmt.Sprintf("<unimplemented stringer>: %s", c.pb.String())
	}

	var b strings.Builder
	b.WriteString(c.str)
	if len(c.enforcement) != 0 {
		b.WriteString(" iff (")
		for i, l := range c.enforcement {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(l.name())
		}
		b.WriteString(")")
	}
	return b.String()
}

var _ Constraint = &constraint{}

// NewAllDifferentConstraint forces all variables to take different values.
func NewAllDifferentConstraint(vars ...IntVar) Constraint {
	var b strings.Builder
	b.WriteString("all-different: [")
	for i, v := range vars {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(v.name())
	}
	b.WriteString("]")

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_AllDiff{
				AllDiff: &pb.AllDifferentConstraintProto{
					Vars: intVarList(vars).indexes(),
				},
			},
		},
		str: b.String(),
	}
}

// NewAllSameConstraint forces all variables to take the same values.
func NewAllSameConstraint(vars ...IntVar) Constraint {
	var b strings.Builder
	b.WriteString("all-same: [")
	for i, v := range vars {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(v.name())
	}
	b.WriteString("]")

	var cs []Constraint
	for i := range vars {
		if i == 0 {
			continue
		}
		cs = append(cs, NewMaximumConstraint(vars[i-1], vars[i]))
	}
	return constraints{cs: cs, str: b.String()}
}

// NewAtMostKConstraint ensures that no more than k literals are true.
func NewAtMostKConstraint(k int, literals ...Literal) Constraint {
	if k == 1 {
		return newAtMostOneConstraint(literals...)
	}

	lb, ub := int64(0), int64(k)
	return NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))
}

// NewAtLeastKConstraint ensures that at least k literals are true.
func NewAtLeastKConstraint(k int, literals ...Literal) Constraint {
	if k == 1 {
		return NewBooleanOrConstraint(literals...)
	}

	lb, ub := int64(k), int64(len(literals))
	return NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))
}

// NewExactlyKConstraint ensures that exactly k literals are true.
func NewExactlyKConstraint(k int, literals ...Literal) Constraint {
	lb, ub := int64(k), int64(k)
	c := NewLinearConstraint(Sum(asIntVars(literals)...), NewDomain(lb, ub))

	var b strings.Builder
	b.WriteString(fmt.Sprintf("exactly-k: k=%d ", k))
	for i, l := range literals {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(l.name())
	}

	c.(*constraint).str = b.String()
	return c
}

// NewBooleanAndConstraint ensures that all literals are true.
func NewBooleanAndConstraint(literals ...Literal) Constraint {
	var b strings.Builder
	b.WriteString("boolean-and: ")
	for i, l := range literals {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(l.name())
	}

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_BoolAnd{
				BoolAnd: &pb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
		str: b.String(),
	}
}

// NewBooleanOrConstraint ensures that at least one literal is true. It can be
// thought of as a special case of NewAtLeastKConstraint, but one that uses a
// more efficient internal encoding.
func NewBooleanOrConstraint(literals ...Literal) Constraint {
	var b strings.Builder
	b.WriteString("boolean-or: ")
	for i, l := range literals {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(l.name())
	}

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_BoolOr{
				BoolOr: &pb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
		str: b.String(),
	}
}

// NewBooleanXorConstraint ensures that an odd number of the literals are true.
func NewBooleanXorConstraint(literals ...Literal) Constraint {
	var b strings.Builder
	b.WriteString("boolean-xor: ")
	for i, l := range literals {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(l.name())
	}

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_BoolXor{
				BoolXor: &pb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
		str: b.String(),
	}
}

// NewImplicationConstraint ensures that the first literal implies the second.
func NewImplicationConstraint(a, b Literal) Constraint {
	return NewBooleanOrConstraint(a.Not(), b)
}

// NewAllowedLiteralAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given literals is one of the listed n-tuple assignments.
func NewAllowedLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	return newLiteralAssignmentsConstraintInternal(literals, assignments)
}

// NewForbiddenLiteralAssignmentsConstraint ensures that the values of the
// n-tuple formed by the given literals is not one of the listed n-tuple
// assignments.
func NewForbiddenLiteralAssignmentsConstraint(literals []Literal, assignments [][]bool) Constraint {
	constraint := newLiteralAssignmentsConstraintInternal(literals, assignments)
	constraint.pb.GetTable().Negated = true
	return constraint
}

// NewDivisionConstraint ensures that the target is to equal to
// numerator/denominator.
func NewDivisionConstraint(target, numerator, denominator IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_IntDiv{
				IntDiv: &pb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVarList([]IntVar{numerator, denominator}).indexes(),
				},
			},
		},
	}
}

// NewProductConstraint ensures that the target to equal to the product of all
// multiplicands.
func NewProductConstraint(target IntVar, multiplicands ...IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_IntProd{
				IntProd: &pb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVarList(multiplicands).indexes(),
				},
			},
		},
	}
}

// NewMaximumConstraint ensures that the target is equal to the maximum of all
// variables.
func NewMaximumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_IntMax{
				IntMax: &pb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVarList(vars).indexes(),
				},
			},
		},
	}
}

// NewMinimumConstraint ensures that the target is equal to the minimum of all
// variables.
func NewMinimumConstraint(target IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_IntMin{
				IntMin: &pb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVarList(vars).indexes(),
				},
			},
		},
	}
}

// NewModuloConstraint ensures that the target to equal to dividend%divisor.
func NewModuloConstraint(target, dividend, divisor IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_IntMod{
				IntMod: &pb.IntegerArgumentProto{
					Target: target.index(),
					Vars:   intVarList([]IntVar{dividend, divisor}).indexes(),
				},
			},
		},
	}
}

// NewAllowedAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given variables is one of the listed n-tuple assignments.
func NewAllowedAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	return newAssignmentsConstraintInternal(vars, assignments)
}

// NewForbiddenAssignmentsConstraint ensures that the values of the n-tuple
// formed by the given variables is not one of the listed n-tuple assignments.
func NewForbiddenAssignmentsConstraint(vars []IntVar, assignments [][]int64) Constraint {
	constraint := newAssignmentsConstraintInternal(vars, assignments)
	constraint.pb.GetTable().Negated = true
	return constraint
}

// NewLinearConstraint ensures that the linear expression lies in the given
// domain. It can be used to express linear equalities of the form:
//
// 		0 <= x + 2y <= 10
//
func NewLinearConstraint(e LinearExpr, d Domain) Constraint {
	var b strings.Builder
	b.WriteString("linear-constraint: ")
	b.WriteString(e.String())
	b.WriteString(" within ")
	b.WriteString(d.String())

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_Linear{
				Linear: &pb.LinearConstraintProto{
					Vars:   e.vars(),
					Coeffs: e.coeffs(),
					Domain: d.list(e.offset()),
				},
			},
		},
		str: b.String(),
	}
}

// NewLinearMaximumConstraint ensures that the target is equal to the maximum of
// all linear expressions.
func NewLinearMaximumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	var b strings.Builder
	b.WriteString("linear-max: ")
	b.WriteString(fmt.Sprintf("%s == max(", target.String()))
	for i, e := range exprs {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.String())
	}
	b.WriteString(")")

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_LinMax{
				LinMax: &pb.LinearArgumentProto{
					Target: target.proto(),
					Exprs:  linearExprList(exprs).protos(),
				},
			},
		},
		str: b.String(),
	}
}

// NewLinearMinimumConstraint ensures that the target is equal to the minimum of
// all linear expressions.
func NewLinearMinimumConstraint(target LinearExpr, exprs ...LinearExpr) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_LinMin{
				LinMin: &pb.LinearArgumentProto{
					Target: target.proto(),
					Exprs:  linearExprList(exprs).protos(),
				},
			},
		},
	}
}

// NewElementConstraint ensures that the target is equal to vars[index].
// Implicitly index takes on one of the values in [0, len(vars)).
func NewElementConstraint(target, index IntVar, vars ...IntVar) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_Element{
				Element: &pb.ElementConstraintProto{
					Target: target.index(),
					Index:  index.index(),
					Vars:   intVarList(vars).indexes(),
				},
			},
		},
	}
}

// NewNonOverlappingConstraint ensures that all the intervals are disjoint.
// More formally, there must exist a sequence such that for every pair of
// consecutive intervals, we have intervals[i].end <= intervals[i+1].start.
// Intervals of size zero matter for this constraint. This is also known as a
// disjunctive constraint in scheduling.
func NewNonOverlappingConstraint(intervals ...Interval) Constraint {
	var b strings.Builder
	b.WriteString("non-overlapping: ")
	for i, itv := range intervals {
		if i != 0 {
			b.WriteString(", ")
		}
		start, end, _ := itv.Parameters()
		b.WriteString(fmt.Sprintf("{%s, %s}", start.name(), end.name()))
	}

	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_NoOverlap{
				NoOverlap: &pb.NoOverlapConstraintProto{
					Intervals: intervalList(intervals).indexes(),
				},
			},
		},
		str: b.String(),
	}
}

// NewNonOverlapping2DConstraint ensures that the boxes defined by the following
// don't overlap:
//
// 		[xintervals[i].start, xintervals[i].end)
// 		[yintervals[i].start, yintervals[i].end)
//
// Intervals/boxes of size zero are considered for overlap if the last argument
// is true.
func NewNonOverlapping2DConstraint(
	xintervals []Interval,
	yintervals []Interval,
	boxesWithNoAreaCanOverlap bool,
) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_NoOverlap_2D{
				NoOverlap_2D: &pb.NoOverlap2DConstraintProto{
					XIntervals: intervalList(xintervals).indexes(),
					YIntervals: intervalList(yintervals).indexes(),

					BoxesWithNullAreaCanOverlap: boxesWithNoAreaCanOverlap,
				},
			},
		},
	}
}

// NewCumulativeConstraint ensures that the sum of the demands of the intervals
// (intervals[i]'s demand is specified in demands[i]) at each interval point
// cannot exceed a max capacity. The intervals are interpreted as [start, end).
// Intervals of size zero are ignored.
func NewCumulativeConstraint(capacity int32, intervals []Interval, demands []int32) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_Cumulative{
				Cumulative: &pb.CumulativeConstraintProto{
					Capacity:  capacity,
					Intervals: intervalList(intervals).indexes(),
					Demands:   demands,
				},
			},
		},
	}
}

// newAtMostOneConstraint is a special case of NewAtMostKConstraint that uses a
// more efficient internal encoding.
func newAtMostOneConstraint(literals ...Literal) Constraint {
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_AtMostOne{
				AtMostOne: &pb.BoolArgumentProto{
					Literals: asIntVars(literals).indexes(),
				},
			},
		},
	}
}

// OnlyEnforceIf is part of the Constraint interface.
func (c *constraint) OnlyEnforceIf(literals ...Literal) Constraint {
	c.pb.EnforcementLiteral = asIntVars(literals).indexes()
	c.enforcement = append(c.enforcement, literals...)
	return c
}

// protos is part of the Constraint interface.
func (c *constraint) protos() []*pb.ConstraintProto {
	return []*pb.ConstraintProto{c.pb}
}

type constraints struct {
	cs  []Constraint
	str string
}

func (c constraints) String() string {
	return c.str
}

var _ Constraint = &constraints{}

// OnlyEnforceIf is part of the Constraint interface.
func (c constraints) OnlyEnforceIf(literals ...Literal) Constraint {
	for _, c := range c.cs {
		c.OnlyEnforceIf(literals...)
	}
	return c
}

// protos is part of the Constraint interface.
func (c constraints) protos() []*pb.ConstraintProto {
	var res []*pb.ConstraintProto
	for _, c := range c.cs {
		res = append(res, c.protos()...)
	}
	return res
}

func newLiteralAssignmentsConstraintInternal(literals []Literal, assignments [][]bool) *constraint {
	var integerAssignments [][]int64
	for _, assignment := range assignments { // convert [][]bool to [][]int64
		var integerAssignment []int64
		for _, a := range assignment {
			i := 0
			if a {
				i = 1
			}
			integerAssignment = append(integerAssignment, int64(i))
		}
		integerAssignments = append(integerAssignments, integerAssignment)
	}

	return newAssignmentsConstraintInternal(asIntVars(literals), integerAssignments)
}

func newAssignmentsConstraintInternal(vars []IntVar, assignments [][]int64) *constraint {
	var values []int64
	for _, assignment := range assignments {
		if len(assignment) != len(vars) {
			panic("mismatched assignment and vars length")
		}
		values = append(values, assignment...)
	}
	return &constraint{
		pb: &pb.ConstraintProto{
			Constraint: &pb.ConstraintProto_Table{
				Table: &pb.TableConstraintProto{
					Vars:   intVarList(vars).indexes(),
					Values: values,
				},
			},
		},
	}
}