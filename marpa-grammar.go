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

type Grammar struct {
	thin     *mt.Grammar
	symbols  map[string]mt.SymbolID
	rule_ids []mt.RuleID
	actions  map[mt.RuleID]func([]interface{}) interface{}
}

type Seq struct {
	Min       int
	Proper    bool
	Separator string
}

func NewGrammar() *Grammar {
	var config mt.Config
	mt.ConfigInit(&config)
	var grammar Grammar
	grammar.thin = mt.NewGrammar(&config)
	grammar.symbols = make(map[string]mt.SymbolID)
	grammar.actions = make(map[mt.RuleID]func([]interface{}) interface{})
	return &grammar
}

func (grammar *Grammar) Precompute() {
	if grammar.thin.Precompute() == -2 {
		cnt := grammar.thin.EventCount()
		fmt.Printf("ERRORS: %d events\n", cnt)
		for i := 0; i < cnt; i += 1 {
			var evt mt.Event
			evttype := grammar.thin.Event(&evt, i)
			fmt.Printf("event type %d\n", evttype)
		}
	}
	return
}

func (grammar *Grammar) IsPrecomputed() bool {
	ret := grammar.thin.IsPrecomputed()
	if ret == -2 {
		panic("ERROR: IsPrecomputed")
	}
	return ret == 1
}

func (grammar *Grammar) StartRule(lhs string) {
	lhs_id := grammar.symbol(lhs)
	grammar.thin.StartSymbolSet(lhs_id)
	return
}

func (grammar *Grammar) AddSequence(lhs string, rhs string, seq Seq, action func(args []interface{}) interface{}) {
	lhs_id := grammar.symbol(lhs)
	rhs_id := grammar.symbol(rhs)
	var sep_id mt.SymbolID

	if seq.Separator != "" {
		sep_id = grammar.symbol(seq.Separator)
	} else {
		sep_id = mt.SymbolID(-1)
	}

	flags := 0
	if seq.Proper {
		flags = 0x2
	}

	rule_id := grammar.thin.NewSequence(lhs_id, rhs_id, sep_id, seq.Min, flags)
	grammar.actions[rule_id] = action
	grammar.rule_ids = append(grammar.rule_ids, rule_id)
	return
}

func (grammar *Grammar) AddRule(lhs string, rhs []string, action func(args []interface{}) interface{}) {
	lhs_id := grammar.symbol(lhs)
	rhs_ids := make([]mt.SymbolID, len(rhs))
	for i, id := range rhs {
		rhs_ids[i] = grammar.symbol(id)
	}
	rule_id := grammar.thin.NewRule(lhs_id, rhs_ids)
	grammar.actions[rule_id] = action
	grammar.rule_ids = append(grammar.rule_ids, rule_id)
	return
}

func (grammar *Grammar) symbol(sym string) mt.SymbolID {
	id, ok := grammar.symbols[sym]
	if !ok {
		id = grammar.thin.NewSymbol()
		grammar.symbols[sym] = id
	}
	return id
}
