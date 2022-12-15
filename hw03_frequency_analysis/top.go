package hw03frequencyanalysis

import (
	"strings"
)

type Pair struct {
	Word  string
	Count uint32
}

type dictionary struct {
	top   []*Pair
	words map[string]uint32
}

func (b *dictionary) topCheck(word *string, count uint32) {
	if count == 0 {
		return
	}

	if count <= b.top[9].Count {
		return
	}

	for idx := int32(9); idx >= 0; idx-- {
		if b.top[idx].Word == *word {
			b.top[idx].Count += 1
			b.upPair(idx)
			return
		}

		if count > b.top[idx].Count {
			continue
		}

		b.replaceTop(word, count, idx)
		return
	}

	b.top[0] = &Pair{
		Word:  *word,
		Count: count,
	}
}

func (b *dictionary) upPair(idx int32) {
	if idx <= 0 {
		return
	}

	if b.top[idx].Count >= b.top[idx-1].Count {
		return
	}

	oldTop := b.top[idx-1]
	b.top[idx-1] = b.top[idx]
	b.top[idx] = oldTop
}

func (b *dictionary) replaceTop(word *string, count uint32, idx int32) {
	if b.top[idx].Count == count && b.top[idx].Word != *word {
		idx += 1
	}

	for addIdx := int32(9); addIdx > idx; addIdx-- {
		b.top[addIdx] = b.top[addIdx-1]
	}

	b.top[idx] = &Pair{
		Word:  *word,
		Count: count,
	}
}

func (b *dictionary) SetValue(word string) {
	if val, ok := b.words[word]; ok {
		b.words[word] = val + 1
		b.topCheck(&word, val+1)
		return
	}

	b.words[word] = 1
	b.topCheck(&word, 1)
}

func (b *dictionary) GetValue() []string {
	words := make([]string, 10)
	for idx, word := range b.top {
		words[idx] = word.Word
	}
	return words[:]
}

func GetDictonary() *dictionary {
	return &dictionary{
		top: func() []*Pair {
			tops := make([]*Pair, 10)
			for idx := uint8(0); idx < 10; idx++ {
				tops[idx] = &Pair{
					Word:  "",
					Count: 0,
				}
			}
			return tops
		}(),
		words: make(map[string]uint32),
	}
}

func Top10(v string) []string {
	dict := GetDictonary()
	for _, word := range strings.Fields(v) {
		dict.SetValue(word)
	}
	return dict.GetValue()
}
