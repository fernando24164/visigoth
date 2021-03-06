package tokenizer

import (
	"bytes"
	"golang.org/x/text/runes"
	"log"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Tokenizer interface {
	TokenizeText([]byte) [][]byte
	TokenizeWord([]byte) ([]byte, bool)
}

type Transformer interface {
	Transform(payload []byte) ([]byte, error)
}

type lowercaseTransformer struct {
	Tf transform.Transformer
}

type IsMarkNonSpacingChecker struct{}

func (c *IsMarkNonSpacingChecker) Contains(r rune) bool {
	return isMarkNonSpacing(r)
}

var runeRemover = &IsMarkNonSpacingChecker{}

func newLowercaseTransformer() *lowercaseTransformer {
	tf := &lowercaseTransformer{}
	tf.Tf = transform.Chain(norm.NFD, runes.Remove(runeRemover), norm.NFC)
	return tf
}

func (t *lowercaseTransformer) Transform(payload []byte) ([]byte, error) {
	final := make([]byte, len(payload))
	nDst, _, err := t.Tf.Transform(final, bytes.ToLowerSpecial(unicode.CaseRanges, payload), true)
	if err != nil {
		return nil, err
	}
	final = final[:nDst]
	return final, nil
}

type SimpleTokenizer struct {
	Tf Transformer
}

func NewSimpleTokenizer() Tokenizer {
	st := new(SimpleTokenizer)
	st.Tf = newLowercaseTransformer()
	return st
}

func (s *SimpleTokenizer) TokenizeText(payload []byte) [][]byte {
	var res [][]byte
	for _, bword := range bytes.Fields(payload) {
		bword, ok := s.TokenizeWord(bword)
		if ok {
			res = append(res, bword)
		}
	}
	return res
}

func (s *SimpleTokenizer) TokenizeWord(payload []byte) ([]byte, bool) {
	bword := bytes.Trim(payload, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
	if len(bword) > 0 {
		word, err := s.Tf.Transform(bword)
		if err != nil {
			log.Println("SimpleTokenizer, fatal error: ", err)
			return nil, false
		}
		return word, true
	}
	return nil, false
}
