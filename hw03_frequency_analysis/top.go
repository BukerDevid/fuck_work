package hw03frequencyanalysis

import (
	"strings"
	"sync"
)

type Pair struct {
	mx sync.Mutex

	Word  string
	Count uint32
	Cost  uint64
}

type dictionary struct {
	mx sync.Mutex

	top   *Top
	words map[string]*uint32
}

type Top struct {
	list     []*Pair
	minCount uint32
}

func (b *dictionary) SetValue(v string, wg *sync.WaitGroup) {
	defer wg.Done()

	if val, ok := b.words[v]; ok {
		*val += 1
		return
	}

	count := uint32(0)
	b.words[v] = &count
}

func (b *dictionary) GetValue() []string {

	return nil
}

func CheckSum(p *Pair) {
	if p.Cost != 0 {
		return
	}

	maxLength := 0
	if len(p.Word) > 4 {
		maxLength = 4
	} else {
		maxLength = len(p.Word)
	}

	for _, sbm := range p.Word[:maxLength] {
		p.Cost += uint64(sbm)
		p.Cost = p.Cost << 8
	}
}

func GetDictonary() *dictionary {
	return &dictionary{
		words: make(map[string]*uint32),
	}
}

func Top10(v string) []string {
	if v == "" {
		return make([]string, 0)
	}

	dict := GetDictonary()
	wg := sync.WaitGroup{}

	for _, word := range strings.Fields(v) {
		word = strings.TrimSpace(word)
		if word == "-" {
			continue
		}

		word = strings.ToLower(word)
		wg.Add(1)
		go dict.SetValue(word, &wg)
	}

	wg.Wait()
	return dict.GetValue()
}
