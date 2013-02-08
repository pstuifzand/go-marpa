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
