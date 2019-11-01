package genericLexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Item struct {
	Type  ItemType
	Value string
	pos   int
	lname *string
}

func (i Item) String() string {
	s := fmt.Sprint(*i.lname, " ", i.pos, " ")
	switch i.Type {
	case ItemError:
		return s + "Error"
	case ItemComma:
		return s + "Comma"
	case ItemParan:
		return s + "Parentheses" + i.Value
	case ItemLetter:
		return s + fmt.Sprintf("Letter \"%s\"", i.Value)
	case ItemWord:
		return s + fmt.Sprintf("Word \"%s\"", i.Value)
	case ItemNumber:
		return s + fmt.Sprint("Number ", i.Value)
	case ItemWSP:
		return s + "WSP"
	default:
		return fmt.Sprintf("%s \"%s\"", s, i.Value)
	}
}

type Lexer struct {
	name      string
	input     string
	start     int
	pos       int
	width     int
	Items     chan Item
	buffer    [3]Item
	peekcount int
}

type ItemType int

const (
	ItemError ItemType = iota
	ItemDot
	ItemEOS
	ItemLetter
	ItemWord
	ItemNumber
	ItemComma
	ItemFlag
	ItemWSP
	ItemParan
)

type stateFn func(*Lexer) stateFn

func Lex(name, input string) (*Lexer, chan Item) {
	l := &Lexer{
		name:  name,
		input: input,
		Items: make(chan Item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.Items
}

const eof = -1

func (l *Lexer) run() {
	for state := lexD; state != nil; {
		state = state(l)
	}
	l.Items <- Item{Type: ItemEOS}
	close(l.Items) // No more tokens will be delivered.
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) NextItem() Item {
	if l.peekcount > 0 {
		l.peekcount--
		//	fmt.Println("NextItem got peeked Item", l.buffer[0].String())
	} else {
		l.buffer[0] = <-l.Items
		//	fmt.Println("NextItem got new Item", l.buffer[0].String())
	}
	return l.buffer[0]
}

func (l *Lexer) PeekItem() Item {
	if l.peekcount > 0 {
		//	fmt.Println("peekItem got already peeked Item", l.buffer[0].String())
		return l.buffer[0]
	}
	//	fmt.Println("peekItem needs new Item")
	l.buffer[0] = l.NextItem()
	l.peekcount = 1
	//	fmt.Println("peekItem got new Item", l.buffer[0].String())
	return l.buffer[0]
}

func lexNumber(l *Lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("e") {
		l.accept("+-")
		l.acceptRun(digits)
	}
	l.emit(ItemNumber)
	return lexD
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) emit(t ItemType) {

	i := Item{t, l.input[l.start:l.pos], l.start, &l.name}
	l.Items <- i
	l.start = l.pos
}

func lexWord(l *Lexer) stateFn {
	l.acceptRun("abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ")
	l.emit(ItemWord)
	return lexD
}

func lexLetter(l *Lexer) stateFn {
	l.accept("abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ")
	if unicode.IsLetter(l.peek()) {
		return lexWord
	}
	l.emit(ItemLetter)
	return lexD
}

func lexComma(l *Lexer) stateFn {
	l.accept(",")
	l.emit(ItemComma)
	return lexD
}

func isWSP(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func lexWSP(l *Lexer) stateFn {
	l.accept(" \t\r\n\f")
	l.emit(ItemWSP)
	return lexD
}

func lexD(l *Lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			l.emit(ItemEOS)
			return nil
		case isWSP(r):
			return lexWSP
		case unicode.IsLetter(r):
			return lexLetter
		case r == '-' || r == '+':
			return lexNumber
		case unicode.IsNumber(r):
			return lexNumber
		case r == ',':
			return lexComma
		case r == '(' || r == ')':
			return lexParan
		default:
			return nil
		}
	}
	return nil
}

func lexParan(l *Lexer) stateFn {
	l.accept("()")
	l.emit(ItemParan)
	return lexD
}

func (l *Lexer) ConsumeWhiteSpace() error {
	for l.PeekItem().Type == ItemWSP {
		l.NextItem()
	}
	return nil
}

func (l *Lexer) ConsumeComma() error {
	for l.PeekItem().Type == ItemComma {
		l.NextItem()
	}
	return nil
}
