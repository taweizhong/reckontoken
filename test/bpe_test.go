package test

import (
	"fmt"
	"github.com/taweizhong/reckontoken"
	"testing"
)

func TestBytePairSplit(t *testing.T) {
	b := reckontoken.NewBPE()
	text := "<|im_start|>system\n你是ChatGPT, 是来帮助人们解决问题的AI模型<|im_end|>\n<|im_start|>user\n<|im_end|>\n<|im_start|>assistant\n"
	encodedTokens := b.Encode(text)
	fmt.Println(len(encodedTokens))
}
