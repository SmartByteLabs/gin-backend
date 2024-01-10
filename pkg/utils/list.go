package utils

type ComparableList[T comparable] struct {
	items []T
}

func NewComparableList[T comparable]() *ComparableList[T] {
	return &ComparableList[T]{
		items: []T{},
	}
}

func (l *ComparableList[T]) Add(item T) {
	l.items = append(l.items, item)
}

func (l *ComparableList[T]) Remove(item T) {
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
}

func (l *ComparableList[T]) Contains(item T) bool {
	for _, v := range l.items {
		if v == item {
			return true
		}
	}
	return false
}

func (l *ComparableList[T]) Size() int {
	return len(l.items)
}

func (l *ComparableList[T]) FindIndex(item T) int {
	for i, v := range l.items {
		if v == item {
			return i
		}
	}
	return -1
}
