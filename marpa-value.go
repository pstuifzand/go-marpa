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
	"fmt"
	mt "github.com/pstuifzand/go-marpa-thin"
)

const DEBUG = false

type Value struct {
	tree  *mt.Tree
	val   *mt.Value
	re    *Recognizer
	value interface{}
}

func resize(stack []interface{}, size int) []interface{} {
	size += 1
	if DEBUG {
		for i, s := range stack {
			fmt.Printf("stack[%d]=%#v\n", i, s)
		}
	}
	l := len(stack)
	if l >= size {
		return stack
	}
	n := make([]interface{}, size)
	copy(n, stack)
	return n
}

func (val *Value) Next() bool {
	b := val.tree.Next() >= 0
	val.val = mt.NewValue(val.tree)
	for _, rule_id := range val.re.grammar.rule_ids {
		val.val.RuleIsValuedSet(rule_id, 1)
	}

	var stack []interface{}

VALUE:
	for {
		step_type := val.val.Step()

		switch step_type {
		case mt.STEP_INITIAL:
			stack = make([]interface{}, 1)
		case mt.STEP_TOKEN:
			res := val.val.Result()
			stack = resize(stack, res)
			if DEBUG {
				fmt.Printf("put stack[%d] <= %s\n", res, val.re.tokens[val.val.TokenValue()])
			}
			stack[res] = val.re.tokens[val.val.TokenValue()]
		case mt.STEP_RULE:
			arg0 := val.val.Arg0()
			argn := val.val.ArgN()
			action, e := val.re.grammar.actions[val.val.Rule()]
			if e {
				res := val.val.Result()
				stack = resize(stack, res)
				if DEBUG {
					fmt.Printf("reduce %d <= [%d -- %d] %d\n", res, arg0, argn, val.val.Rule())
				}
				args := stack[arg0 : argn+1]
				stack[res] = action(args)
			}
		case mt.STEP_NULLING_SYMBOL:
			res := val.val.Result()
			stack = resize(stack, res+1)
			stack[res] = val.re.tokens[val.val.TokenValue()]
		case mt.STEP_INACTIVE:
			break VALUE
		}
	}
	val.value = stack[0]
	return b
}

func (val *Value) Value() interface{} {
	return val.value
}
