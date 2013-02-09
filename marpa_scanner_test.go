/* Copyright (C) 2013 Peter Stuifzand

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Lesser General Public License as published by the Free
Software Foundation, either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public License for more
details.

You should have received a copy of the GNU Lesser General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package marpa

import (
	//"fmt"
	"testing"
)

type Rule struct {
	Lhs string
	Rhs []string
}

func ActionRules(args []interface{}) interface{} {
	rules := []*Rule{}
	for _, n := range args {
		rule := n.(*Rule)
		rules = append(rules, rule)
	}
	return rules
}
func ActionRule(args []interface{}) interface{} {
	lhs := args[0].(string)
	rhs := args[2].([]string)
	return &Rule{Lhs: lhs, Rhs: rhs}
}
func ActionPlus(args []interface{}) interface{} {
	// not implemented
	return args
}
func ActionStar(args []interface{}) interface{} {
	// not implemented
	return args
}
func ActionLhs(args []interface{}) interface{} {
	return args[0]
}
func ActionRhs(args []interface{}) interface{} {
	return args[0]
}
func ActionNames(args []interface{}) interface{} {
	names := []string{}
	for _, n := range args {
		name := n.(string)
		names = append(names, name)
	}
	return names
}

func TestMarpaRules(t *testing.T) {
	g := NewGrammar()

	g.StartRule("rules")
	g.AddSequence("rules", "rule", Seq{Min: 1}, ActionRules)
	g.AddRule("rule", []string{"lhs", "bnfop", "rhs"}, ActionRule)
	g.AddRule("lhs", []string{"name"}, ActionLhs)
	g.AddRule("rhs", []string{"names"}, ActionRhs)
	g.AddRule("rhs", []string{"name", "plus"}, ActionPlus)
	g.AddRule("rhs", []string{"name", "star"}, ActionStar)
	g.AddSequence("names", "name", Seq{Min: 1}, ActionNames)

	g.Precompute()

	re, err := NewRecognizer(g)
	if err != nil {
		t.Errorf("Error should not be set: grammar is precomputed")
	}

	re.Read("name", "start")
	re.Read("bnfop", "::=")
	re.Read("name", "expression")

	re.Read("name", "expression")
	re.Read("bnfop", "::=")
	re.Read("name", "number")

	re.Read("name", "expression")
	re.Read("bnfop", "::=")
	re.Read("name", "expression")
	re.Read("name", "op")
	re.Read("name", "expression")

	val := re.Value()

	for val.Next() {
		rules := val.Value().([]*Rule)
		if rules[0].Lhs != "start" || rules[0].Rhs[0] != "expression" {
			t.Errorf("rules are not expected")
		}
		if rules[1].Lhs != "expression" || rules[1].Rhs[0] != "number" {
			t.Errorf("rules are not expected")
		}
		if rules[2].Lhs != "expression" {
			t.Errorf("rules are not expected")
		}
		if rules[2].Rhs[0] != "expression" || rules[2].Rhs[1] != "op" || rules[2].Rhs[0] != "expression" {
			t.Errorf("rules are not expected")
		}
	}
}
