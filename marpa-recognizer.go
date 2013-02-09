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
	"errors"
	"fmt"
	mt "github.com/pstuifzand/go-marpa-thin"
)

type Recognizer struct {
	thin        *mt.Recognizer
	grammar     *Grammar
	tokens      map[int]string
	tokensi     map[string]int
	token_count int
}

func (re *Recognizer) token(token string) int {
	id, e := re.tokensi[token]
	if !e {
		id = re.token_count
		re.tokens[id] = token
		re.tokensi[token] = id
		re.token_count++
	}
	return id
}

func NewRecognizer(grammar *Grammar) (*Recognizer, error) {
	if !grammar.IsPrecomputed() {
		return nil, errors.New("Grammar is not precomputed. Call grammar.Precompute()")
	}
	thin_re := mt.NewRecognizer(grammar.thin)
	recognizer := &Recognizer{thin_re, grammar, nil, nil, 1}
	recognizer.tokens = make(map[int]string)
	recognizer.tokensi = make(map[string]int)
	thin_re.StartInput()
	return recognizer, nil
}

func (re *Recognizer) Read(terminal, value string) {
	re.thin.Alternative(re.grammar.symbol(terminal), re.token(value), 1)
	re.thin.EarlemeComplete()
}

func (re *Recognizer) Value() *Value {
	latest_earley_set_id := re.thin.LatestEarleySet()
	bocage, err := mt.NewBocage(re.thin, latest_earley_set_id)
	if err != nil {
		fmt.Printf("ERROR: %s %s\n", err, mt.ErrStr[re.grammar.thin.Error()])
		panic("ERROR")
	}
	order := mt.NewOrder(bocage)
	tree := mt.NewTree(order)
	value := &Value{tree, nil, re, ""}
	return value
}
