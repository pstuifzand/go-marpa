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
	"testing"
)

func ActionArg0(args []interface{}) interface{} {
	return args[0]
}

func ActionExprOp(args []interface{}) interface{} {
	l := args[0].(string)
	op := args[1].(string)
	r := args[2].(string)
	return "(" + l + " " + op + " " + r + ")"
}

func TestMarpa(t *testing.T) {
	g := NewGrammar()
	g.StartRule("start")
	g.AddRule("start", []string{"expression"}, ActionArg0)
	g.AddRule("expression", []string{"expression", "op", "expression"}, ActionExprOp)
	g.AddRule("expression", []string{"number"}, ActionArg0)
	g.Precompute()

	re, err := NewRecognizer(g)
	if err != nil {
		t.Errorf("Error should not be set: grammar is precomputed")
	}

	re.Read("number", "5")
	re.Read("op", "-")
	re.Read("number", "4")
	re.Read("op", "*")
	re.Read("number", "3")
	re.Read("op", "+")
	re.Read("number", "1")

	val := re.Value()

	for val.Next() {
		e := "(((5 - 4) * 3) + 1)"
		if val.Value() != e {
			t.Errorf("%s != %s", val.Value(), e)
		}
	}
}
