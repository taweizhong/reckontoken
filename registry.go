package reckontoken

import (
	"fmt"
	"log"
)

var ENCODINGS map[string]*BPE

func init() {
	RegistryConstructors()
}

func NewBase(url string, patStr string, special_tokens map[string]int) *BPE {
	mergeable_ranks, err := LoadTokens(url)
	if err != nil {
		log.Fatalln(fmt.Errorf("file error: %v", err))
	}
	if mergeable_ranks == nil || len(mergeable_ranks) == 0 {
		log.Fatalln(fmt.Errorf("format error: mergeable_ranks error"))
	}
	if patStr != "" {
		log.Fatalln(fmt.Errorf("patStr error"))
	}
	if special_tokens != nil {
		log.Fatalln(fmt.Errorf("special_tokens error"))
	}
	bpe := NewBPE(
		WithEncoder(mergeable_ranks),
		WithDecoder(mergeable_ranks),
		WithRegex(patStr),
		WithSpecialRegex(special_tokens),
		WithSpecialTokensEncoder(special_tokens),
	)
	return bpe
}

// 注册工厂函数
func RegistryConstructors() {
	ENCODINGS = make(map[string]*BPE, len(ENCODINGCONSTRUCTORS))
	for name, f := range ENCODINGCONSTRUCTORS {
		base := f()
		bpe := NewBPE(
			WithEncoder(base.mergeAbleRanks),
			WithDecoder(base.mergeAbleRanks),
			WithRegex(base.patStr),
			WithSpecialRegex(base.specialTokens),
			WithSpecialTokensEncoder(base.specialTokens),
		)
		ENCODINGS[name] = bpe
	}
}

// 获取具体的编码器
func GetEncoder(encoderName string) *BPE {
	if b, ok := ENCODINGS[encoderName]; ok {
		return b
	}
	return nil
}
