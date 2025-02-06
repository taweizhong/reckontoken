package reckontoken

var ENCODINGS map[string]*BPE

func init() {
	RegistryConstructors()
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
