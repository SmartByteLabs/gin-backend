package utils

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: map[T]struct{}{},
	}
}

func NewSetFromSlice[T comparable](slice []T) *Set[T] {
	set := NewSet[T]()
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

func (s *Set[T]) Add(items ...T) *Set[T] {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *Set[T]) Remove(item T) *Set[T] {
	delete(s.m, item)
	return s
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.m)
}

func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.m))
	for item := range s.m {
		slice = append(slice, item)
	}
	return slice
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	union := NewSet[T]()
	for item := range s.m {
		union.Add(item)
	}
	for item := range other.m {
		union.Add(item)
	}
	return union
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	intersection := NewSet[T]()
	for item := range s.m {
		if other.Contains(item) {
			intersection.Add(item)
		}
	}
	return intersection
}

func (s *Set[T]) GetCommonElements(ar []T) []T {
	out := make([]T, 0)
	for _, item := range ar {
		if s.Contains(item) {
			out = append(out, item)
		}
	}

	return out
}
