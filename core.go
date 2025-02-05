package reckontoken

import (
	"log"
	"regexp"
)

type BPE struct {
	encoder                map[string]int
	special_tokens_encoder map[string]int
	decoder                map[int]string
	special_tokens_decoder map[int]string
	regex_tls              string
	special_regex_tls      string
	sorted_token_bytes     string
}

func NewBPE() *BPE {
	path := TokenFilePath["o200k_base"]
	tokens, _ := LoadTokens(path)
	return &BPE{
		encoder: tokens,
	}
}

func (b *BPE) getRegex() string {
	return b.regex_tls
}
func (b *BPE) getSpecialTokens() string {
	return b.special_regex_tls
}
func (b *BPE) decodeNative(ir []indexRank) string {
	result := ""
	for _, token := range ir {
		result += b.decoder[token.rank]
	}
	return result
}
func (b *BPE) encodeOrdinaryNative(text string) []int {
	regex := b.getRegex()
	ret := make([]int, 0)
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}
	result := re.FindAllString(text, -1)
	for _, mat := range result {
		if piece, ok := b.encoder[mat]; ok {
			ret = append(ret, piece)
		} else {
			irList := bytePairMerge(b.encoder, mat)
			ranks := irList.getRank()
			ret = append(ret, ranks...)
		}
	}
	return ret
}
func (b *BPE) Encode(text string) []string {
	irList := BytePairEncode(b.encoder, text)
	return irList
}
