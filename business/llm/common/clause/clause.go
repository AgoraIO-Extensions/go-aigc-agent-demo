package clause

var CharMap map[rune]bool

func init() {
	CharMap = make(map[rune]bool)
	charSet := ",，.。!！?？:：;；)）*~、"
	for _, c := range charSet {
		CharMap[c] = true
	}
}
