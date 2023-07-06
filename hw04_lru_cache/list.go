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

type list struct {
	count int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	item.addBetween(l.back, nil)
	l.back = item
	if l.front == nil {
		l.front = item
	}
	l.count++

	return item
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	item.addBetween(nil, l.front)
	l.front = item
	if l.back == nil {
		l.back = item
	}
	l.count++

	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i == l.front {
		l.front = i.Next
	}

	if i == l.back {
		l.back = i.Prev
	}

	i.removeSelf()
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if i == l.front {
		return
	}

	if i == l.back {
		l.back = i.Prev
	}

	i.removeSelf()
	i.addBetween(nil, l.front)
	l.front = i
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func (i *ListItem) removeSelf() {
	if i.Prev != nil {
		i.Prev.Next = nil
	}

	if i.Next != nil {
		i.Next.Prev = nil
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
}

func (item *ListItem) addBetween(prev, next *ListItem) {
	if prev != nil {
		prev.Next = item
		item.Prev = prev
	}

	if next != nil {
		next.Prev = item
		item.Next = next
	}
}

func NewList() List {
	return new(list)
}
