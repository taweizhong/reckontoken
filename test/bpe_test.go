package test

import (
	"encoding/base64"
	"fmt"
	"github.com/taweizhong/reckontoken"
	"strings"
	"testing"
)

func TestEncodeOrdinary(t *testing.T) {
	encoder := reckontoken.GetEncoder("cl100k_base")
	text := "By studying DeepSeek R1’s architecture, we managed to selectively quantize certain layers to higher bits (like 4bit) & leave most MoE layers (like those used in GPT-4) to 1.5bit. Naively quantizing all layers breaks the model entirely, causing endless loops & gibberish outputs. Our dynamic quants solve this."
	encodedTokens := encoder.EncodeOrdinary(text)
	fmt.Println(len(encodedTokens))
}
func TestSplit(t *testing.T) {
	encoder := reckontoken.GetEncoder("cl100k_base")
	text := "DeepSeek"
	encodedTokens := encoder.Split(text)
	fmt.Println(encodedTokens)
}
func TestEncode(t *testing.T) {
	encoder := reckontoken.GetEncoder("cl100k_base")
	text := "DeepSeek-R1 has been making waves recently by rivaling OpenAI's O1 reasoning model while being fully open-source. We explored how to enable more local users to run it & managed to quantize DeepSeek’s R1 671B parameter model to 131GB in size, a 80% reduction in size from the original 720GB, whilst being very functional."
	allowedSpecial := map[string]bool{"<|endoftext|>": true, "<|special|>": true}
	encodedTokens := encoder.Encode(text, allowedSpecial)
	fmt.Println(len(encodedTokens))
}
func TestBase64Encode(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace("IGJvbWJlcg=="))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decoded))
}
