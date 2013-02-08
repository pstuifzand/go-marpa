package marpa

import (
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
	//fmt.Printf("%s = %d\n", token, id)
	return id
}

func NewRecognizer(grammar *Grammar) *Recognizer {
	thin_re := mt.NewRecognizer(grammar.thin)
	recognizer := &Recognizer{thin_re, grammar, nil, nil, 1}
	recognizer.tokens = make(map[int]string)
	recognizer.tokensi = make(map[string]int)
	thin_re.StartInput()
	return recognizer
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
