package reckontoken

var ENCODINGS map[string]*BPE

func init() {
	RegistryConstructors()
}
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
func GetEncoder(encoderName string) *BPE {
	if b, ok := ENCODINGS[encoderName]; ok {
		return b
	}
	return nil
}
