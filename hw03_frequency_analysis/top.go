package hw03frequencyanalysis

import (
	"log"
	"sort"
	"strings"
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
		p.Cost <<= 8
	}
}

type Top struct {
	list     [21]*Pair
	minValue uint32
}

func (t *Top) SetInTop(p *Pair) {
	defer func() {
		t.minValue = t.list[0].Count
	}()
	if p.Count < t.minValue {
		return
	}
	p.CheckCost()
	var exist bool
	counter := 0
	for currentPair, pair := range t.list {
		if p.Count > pair.Count && currentPair < 20 {
			counter++
			continue
		}
		if counter > 0 {
			counter--
		}
		for _, val := range t.list {
			if val.Word == p.Word {
				exist = true
				break
			}
		}
		for rewrite := 0; rewrite < currentPair; rewrite++ {
			if !exist {
				t.list[rewrite] = t.list[rewrite+1]
				continue
			}
			if p.Word == t.list[rewrite].Word {
				t.list[rewrite] = t.list[rewrite+1]
				exist = false
			}
		}
		t.list[counter] = p
		break
	}
}

type dictionary struct {
	top   *Top
	words map[string]*uint32
}

func (b *dictionary) SetValue(v string) {
	count := uint32(1)
	if val, ok := b.words[v]; ok {
		*val += 1
		count = *val
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
		if b.top.list[i].Count < b.top.list[j].Count {
			return true
		}
		if b.top.list[i].Count == b.top.list[j].Count {
			if b.top.list[i].Cost > b.top.list[j].Cost {
				return true
			}
		}
		return false
	})
	for _, val := range b.top.list {
		log.Println(val)
	}
	top := b.top.list[11:]
	for idx := 9; idx >= 0; idx-- {
		if top[idx].Cost > 0 {
			result = append(result, top[idx].Word)
		}
	}
	return result
}

func GetDictonary() *dictionary {
	return &dictionary{
		top: &Top{
			list: func() [21]*Pair {
				var list [21]*Pair
				for idx := 0; idx < len(list); idx++ {
					list[idx] = &Pair{}
				}
				return list
			}(),
			minValue: 0,
		},
		words: make(map[string]*uint32),
	}
}

func Top10(v string) []string {
	if v == "" {
		return make([]string, 0)
	}
	dict := GetDictonary()
	for _, word := range strings.Fields(v) {
		word := strings.TrimSpace(word)
		word = strings.Trim(word, ",")
		word = strings.Trim(word, ".")
		word = strings.Trim(word, ":")
		word = strings.Trim(word, ";")
		word = strings.Trim(word, "!")
		word = strings.Trim(word, "?")
		word = strings.Trim(word, "\"")
		word = strings.Trim(word, "'")
		if word == "-" {
			continue
		}
		word = strings.ToLower(word)
		dict.SetValue(word)
	}
	return dict.GetValue()
}
