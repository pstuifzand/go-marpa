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
	"testing"
)

func ActionArg0(args []string) string {
	return args[0]
}
func ActionExprOp(args []string) string {
	return "(" + args[0] + " " + args[1] + " " + args[2] + ")"
}

func TestMarpa(t *testing.T) {
	g := NewGrammar()
	g.StartRule("start")
	g.AddRule("start", []string{"expression"}, ActionArg0)
	g.AddRule("expression", []string{"expression", "op", "expression"}, ActionExprOp)
	g.AddRule("expression", []string{"number"}, ActionArg0)
	g.Precompute()

	re := NewRecognizer(g)

	re.Read("number", "5")
	re.Read("op", "-")
	re.Read("number", "4")
	re.Read("op", "*")
	re.Read("number", "3")
	re.Read("op", "+")
	re.Read("number", "1")

	val := re.Value()

	for val.Next() {
		fmt.Printf("%s\n", val.Value())
	}
}
