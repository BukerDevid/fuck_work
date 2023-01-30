package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size int32

	first *ListItem
	last  *ListItem
}

func (ls *list) Len() int {
	return (int)(ls.size)
}

func (ls *list) Front() *ListItem { return ls.first }

func (ls *list) Back() *ListItem { return ls.last }

func (ls *list) checkFirst(new *ListItem) bool {
	if ls.first == nil {
		ls.first = new
		ls.last = new
		return false
	}
	return true
}

func (ls *list) checkRelation() {
	if ls.size == 2 {
		if ls.first == nil && ls.last != nil {
			ls.first = ls.last
		} else if ls.last == nil && ls.first != nil {
			ls.last = ls.first
		}
	}
}

func (ls *list) PushFront(v interface{}) *ListItem {
	defer func() {
		ls.size++
	}()
	newItem := &ListItem{
		Value: v,
	}
	if !ls.checkFirst(newItem) {
		return ls.first
	}
	buf := ls.first
	buf.Prev = newItem
	ls.first = newItem
	ls.first.Next = buf
	return ls.first
}

func (ls *list) PushBack(v interface{}) *ListItem {
	defer func() {
		ls.size++
	}()
	newItem := &ListItem{
		Value: v,
	}
	if !ls.checkFirst(newItem) {
		return ls.first
	}
	buf := ls.last
	buf.Next = newItem
	ls.last = newItem
	ls.last.Prev = buf
	return ls.last
}

func (ls *list) Remove(i *ListItem) {
	defer func(i *ListItem) {
		ls.size--
		ls.checkRelation()
	}(i)
	if ls.size == 0 {
		return
	}
	if i.Next != nil && i.Prev != nil {
		i.Next.Prev = i.Prev
		if ls.last == i {
			ls.last = i.Next
		}
		i.Prev.Next = i.Next
		if ls.first == i {
			ls.first = i.Prev
		}
	}
}

func (ls *list) MoveToFront(i *ListItem) {
	if ls.size < 2 || ls.first == i {
		return
	}
	buf := ls.first
	buf.Prev = i
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	ls.first = i
	ls.first.Next = buf
	ls.first.Prev = nil
}

func NewList() List {
	return new(list)
}
