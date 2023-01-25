package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"sync"
)

type Pair struct {
	Word  string
	Count uint32
	Cost  uint64
}

func (p *Pair) CheckCost() {
	if p.Cost > 0 {
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

type Top struct {
	mx       sync.Mutex
	list     [10]*Pair
	minValue uint32
}

func (t *Top) SetInTop(p *Pair) {
	t.mx.Lock()
	defer func() {
		t.minValue = t.list[0].Count
		t.mx.Unlock()
	}()
	if p.Count < t.minValue {
		return
	}
	p.CheckCost()
	for currentPair, pair := range t.list {
		if pair.Count < p.Count {
			continue
		}
		if currentPair > 0 {
			for rewrite := 0; rewrite < currentPair; rewrite++ {
				if t.list[rewrite+1].Word == pair.Word {
					continue
				}
				if rewrite < 9 {
					t.list[rewrite] = t.list[rewrite+1]
				}
			}
			t.list[currentPair] = pair
		}
		break
	}
}

type dictionary struct {
	mx    sync.Mutex
	top   *Top
	words map[string]*uint32
}

func (b *dictionary) SetValue(v string, wg *sync.WaitGroup) {
	b.mx.Lock()
	defer func() {
		wg.Done()
		b.mx.Unlock()
	}()
	count := uint32(0)
	if val, ok := b.words[v]; ok {
		*val += 1
		count = *val
		return
	}
	b.words[v] = &count
	pair := &Pair{
		Word:  v,
		Count: count,
	}
	b.top.SetInTop(pair)
}

func (b *dictionary) GetValue() []string {
	result := make([]string, 0, 10)
	sort.Slice(b.top.list[:], func(i, j int) bool {
		if b.top.list[i].Count > b.top.list[j].Count {
			if b.top.list[i].Cost > b.top.list[j].Cost {
				return true
			}
		}
		return false
	})
	for _, pair := range b.top.list {
		result = append(result, pair.Word)
	}
	return result
}

func GetDictonary() *dictionary {
	return &dictionary{
		top: &Top{
			list: [10]*Pair{},
		},
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
		word := strings.TrimSpace(word)
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
