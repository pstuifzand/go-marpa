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

func actionArg0(args []interface{}) interface{} {
	return args[0]
}

func TestMarpaNoPrecompute(t *testing.T) {
	g := NewGrammar()

	g.StartRule("start")
	g.AddRule("start", []string{"expression"}, actionArg0)
	// g.Precompute() // Don't precompute

	_, err := NewRecognizer(g)
	if err == nil {
		t.Errorf("Error should be set of grammar is not precomputed")
	}
}

func TestMarpaNoStartRule(t *testing.T) {
	g := NewGrammar()
	g.StartRule("start")
	g.AddRule("start", []string{"expression"}, actionArg0)
	err := g.Precompute()
	if err == nil {
		t.Fatalf("Err == nil, but should be a marpa error")
	}
	if _, ok := err.(MarpaError); ok {
		t.Logf("MarpaError")
	}
}
