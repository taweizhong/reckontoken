package reckontoken

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// option函数
type Option func(bpe *BPE)

// BPE结构体
type BPE struct {
	encoder                map[string]int // 编码词典
	special_tokens_encoder map[string]int // 特殊token的编码
	decoder                map[int]string // 解码词典
	special_tokens_decoder map[int]string // 特殊token的解码
	regex                  *regexp.Regexp // 匹配普通token的正则，用于划分
	special_regex          *regexp.Regexp // 匹配特殊token的正则
	sorted_token_bytes     string         // 排序
}

// 创建一个BPE
func NewBPE(opt ...Option) *BPE {
	bpe := &BPE{}
	for _, opt := range opt {
		opt(bpe)
	}
	return bpe
}

func WithEncoder(encoder map[string]int) Option {
	return func(bpe *BPE) {
		bpe.encoder = encoder
	}
}

func WithDecoder(encoder map[string]int) Option {
	decoder := make(map[int]string)
	for k, v := range encoder {
		decoder[v] = k
	}
	return func(bpe *BPE) {
		bpe.decoder = decoder
	}
}

func WithRegex(regex_tls string) Option {
	regex, err := regexp.Compile(regex_tls)
	if err != nil {
		log.Fatalln(fmt.Errorf("regex error: %v", err))
	}
	return func(bpe *BPE) {
		bpe.regex = regex
	}
}

func WithSpecialRegex(specialTokens map[string]int) Option {
	escapedTokens := make([]string, len(specialTokens))
	i := 0
	for s, _ := range specialTokens {
		escapedTokens[i] = regexp.QuoteMeta(s)
		i++
	}
	pattern := strings.Join(escapedTokens, "|") // 拼接成正则
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalln(fmt.Errorf("regex error: %v", err))
	}
	return func(bpe *BPE) {
		bpe.special_regex = regex
	}
}

func WithSpecialTokensEncoder(special_tokens_encoder map[string]int) Option {
	return func(bpe *BPE) {
		bpe.special_tokens_encoder = special_tokens_encoder
	}
}

func (b *BPE) getRegex() *regexp.Regexp        { return b.regex }
func (b *BPE) getSpecialRegex() *regexp.Regexp { return b.special_regex }

// 解码
func (b *BPE) DecodeNative(ir []indexRank) string {
	result := ""
	for _, token := range ir {
		result += b.decoder[token.rank]
	}
	return result
}

// 分词
func (b *BPE) Split(text string) []string {
	regex := b.getRegex()
	ret := make([]string, 0)
	result := regex.FindAllString(text, -1)
	for _, mat := range result {
		if _, ok := b.encoder[mat]; ok {
			ret = append(ret, mat)
		} else {
			tokens := bytePairSplit(b.encoder, mat)
			ret = append(ret, tokens...)
		}
	}
	return ret
}

// 原始编码，不排除特殊token
func (b *BPE) EncodeOrdinary(text string) []int {
	regex := b.getRegex()
	ret := make([]int, 0)
	result := regex.FindAllString(text, -1)
	for _, mat := range result {
		if piece, ok := b.encoder[mat]; ok {
			ret = append(ret, piece)
		} else {
			ranks := bytePairEncode(b.encoder, mat)
			ret = append(ret, ranks...)
		}
	}
	return ret
}

// 编码，排除特殊的token
// 该方法匹配所有的特殊token,如果匹配到的token在allowedSpecial中，会将该token的rank值加入到ret作为一个token。
func (b *BPE) Encode(text string, allowedSpecial map[string]bool) []int {
	regex := b.getRegex()
	specialRegex := b.getSpecialRegex()
	var ret []int
	start := 0
	for {
		startFind := start
		var specialIndex []int
		// 匹配特殊的token
		for {
			// 从前往后匹配第一个特殊的token
			specialIndex = specialRegex.FindStringIndex(text[startFind:])
			// 没有匹配到直接返回
			if specialIndex == nil {
				break
			} else {
				// 匹配到了且在允许的特殊token列表中直接返回
				if v, ok := allowedSpecial[text[specialIndex[0]:specialIndex[1]]]; ok {
					if v {
						break
					}
				}
				// 匹配到了但是不在允许的特殊token列表中，进行下一个特殊token的匹配
				startFind = specialIndex[0] + 1
			}
		}
		end := len(text)
		// 匹配到了且在允许的特殊token列表中时，specialIndex才不为nil
		if specialIndex != nil {
			end = specialIndex[0]
		}
		// 处理普通文本
		matches := regex.FindAllString(text[start:end], -1)
		for _, piece := range matches {
			if token, exists := b.encoder[piece]; exists {
				ret = append(ret, token)
				continue
			}
			// BPE 编码
			tokens := bytePairEncode(b.encoder, piece)
			ret = append(ret, tokens...)
		}

		// 处理特殊 token
		if specialIndex != nil {
			token := b.special_tokens_encoder[text[specialIndex[0]:specialIndex[1]]]
			ret = append(ret, token)
			start = specialIndex[1]
		} else {
			break
		}
	}
	return ret
}
