package reckontoken

import (
	"fmt"
	"log"
	"strings"
)

// 编码器
var ENCODINGCONSTRUCTORS = map[string]factory{
	"cl100k_base": cl100kBase(),
	"o200k_base":  o200kBase(),
}

var (
	ENDOFTEXT   = "<|endoftext|>"
	FIM_PREFIX  = "<|fim_prefix|>"
	FIM_MIDDLE  = "<|fim_middle|>"
	FIM_SUFFIX  = "<|fim_suffix|>"
	ENDOFPROMPT = "<|endofprompt|>"
)

// 词典下载地址
var TokenFilePath = map[string]string{
	"o200k_base":  "https://gitee.com/taweizhong/encodings/raw/master/o200k_base.tiktoken",
	"cl100k_base": "https://gitee.com/taweizhong/encodings/raw/master/cl100k_base.tiktoken",
}

// 工厂函数
type factory func() *Base

// 基础数据
type Base struct {
	name           string
	mergeAbleRanks map[string]int
	patStr         string
	specialTokens  map[string]int
}

// cl100kBase
func cl100kBase() factory {
	return func() *Base {
		mergeable_ranks, err := LoadTokens(TokenFilePath["cl100k_base"])
		if err != nil {
			log.Fatalln(fmt.Errorf("cl100k_base load error: %v", err))
		}
		special_tokens := map[string]int{
			ENDOFTEXT:   100257,
			FIM_PREFIX:  100258,
			FIM_MIDDLE:  100259,
			FIM_SUFFIX:  100260,
			ENDOFPROMPT: 100276,
		}
		return &Base{
			name:           "cl100kBase",
			mergeAbleRanks: mergeable_ranks,
			specialTokens:  special_tokens,
			patStr:         `(?i)'([sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s+$|\s*[\r\n]|\s+|\s`,
		}
	}

}

// o200kBase
func o200kBase() factory {
	return func() *Base {
		mergeable_ranks, err := LoadTokens(TokenFilePath["o200k_base"])
		if err != nil {
			log.Fatalln(fmt.Errorf("o200kBase load error: %v", err))
		}
		special_tokens := map[string]int{ENDOFTEXT: 199999, ENDOFPROMPT: 200018}
		patterns := []string{
			`[^\r\n\pL\pN]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]*[\p{Ll}\p{Lm}\p{Lo}\p{M}]+(?i:'s|'t|'re|'ve|'m|'ll|'d)?`,
			`[^\r\n\pL\pN]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]+[\p{Ll}\p{Lm}\p{Lo}\p{M}]*((?i:'s|'t|'re|'ve|'m|'ll|'d)?)`,
			`\p{N}{1,3}`,
			` ?[^\s\pL\pN]+[\r\n/]*`,
			`\s*[\r\n]+`,
			`\s+(?!\S)`,
			`\s+`,
		}
		// 使用 | 连接所有模式
		patStr := strings.Join(patterns, "|")
		return &Base{
			name:           "o200k_base",
			mergeAbleRanks: mergeable_ranks,
			patStr:         patStr,
			specialTokens:  special_tokens,
		}
	}
}
