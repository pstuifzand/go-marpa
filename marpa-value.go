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
	mt "github.com/pstuifzand/go-marpa-thin"
)

type Value struct {
	tree  *mt.Tree
	val   *mt.Value
	re    *Recognizer
	value interface{}
}

func (val *Value) Next() bool {
	b := val.tree.Next() >= 0
	val.val = mt.NewValue(val.tree)
	for _, rule_id := range val.re.grammar.rule_ids {
		val.val.RuleIsValuedSet(rule_id, 1)
	}
	stack := make([]interface{}, 100)
VALUE:
	for {
		step_type := val.val.Step()

		switch step_type {
		case mt.STEP_TOKEN:
			argn := val.val.ArgN()
			stack[argn] = val.re.tokens[val.val.TokenValue()]
		case mt.STEP_RULE:
			arg0 := val.val.Arg0()
			argn := val.val.ArgN()
			args := stack[arg0 : argn+1]
			action, e := val.re.grammar.actions[val.val.Rule()]
			if e {
				stack[arg0] = action(args)
			}
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
