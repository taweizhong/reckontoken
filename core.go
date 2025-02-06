package reckontoken

import (
	"math"
)

var rankMax = math.MaxInt

type IRList []indexRank

// 核心结构体，记录索引位置和rank
type indexRank struct {
	index int
	rank  int
}

// 分割成每一个token的编码
func bytePairEncode(ranks map[string]int, piece string) []int {
	ir := bytePairMerge(ranks, piece)
	rankList := make([]int, 0)
	for i := 0; i < len(ir)-1; i++ {
		token := piece[ir[i].index:ir[i+1].index]
		rank := ranks[token]
		rankList = append(rankList, rank)
	}
	return rankList
}

// 将字符串使用BPE方法分割
func bytePairSplit(ranks map[string]int, piece string) []string {
	ir := bytePairMerge(ranks, piece)
	tokens := make([]string, 0)
	for i := 0; i < len(ir)-1; i++ {
		token := piece[ir[i].index:ir[i+1].index]
		tokens = append(tokens, token)
	}
	return tokens
}

// 字节对合并
func bytePairMerge(ranks map[string]int, piece string) IRList {
	parts := make([]indexRank, 0, len(piece)+1)

	minRank := indexRank{
		index: rankMax,
		rank:  rankMax,
	}
	for i := 0; i < len(piece)-1; i++ {
		rank := rankMax
		if r, exist := ranks[piece[i:i+2]]; exist {
			rank = r
			if rank < minRank.rank {
				minRank = indexRank{index: i, rank: rank}
			}
		}
		parts = append(parts, indexRank{i, rank})
	}
	// 确保最后一个字符被处理
	parts = append(parts, indexRank{len(piece) - 1, rankMax})
	parts = append(parts, indexRank{len(piece), rankMax})

	getRank := func(parts []indexRank, index int) int {
		rank := rankMax
		if index+3 < len(parts) {
			if r, exist := ranks[piece[parts[index].index:parts[index+3].index]]; exist {
				rank = r
				return rank
			}
		}
		return rank
	}
	for minRank.rank < rankMax {
		index := minRank.index
		if index > 0 {
			parts[index-1].rank = getRank(parts, index-1)
		}
		parts[index].rank = getRank(parts, index)
		l := parts[:index+1]
		r := parts[index+2:]
		parts = append(l, r...)
		minRank = indexRank{
			index: rankMax,
			rank:  rankMax,
		}
		for i := 0; i < len(parts)-1; i++ {
			if parts[i].rank < minRank.rank {
				minRank.rank = parts[i].rank
				minRank.index = i
			}
		}
	}
	return parts
}
