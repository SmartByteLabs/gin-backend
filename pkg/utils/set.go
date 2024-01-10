package utils

type set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *set[T] {
	return &set[T]{
		m: map[T]struct{}{},
	}
}

func NewSetFromSlice[T comparable](slice []T) *set[T] {
	set := NewSet[T]()
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

func (s *set[T]) Add(item T) *set[T] {
	s.m[item] = struct{}{}
	return s
}

func (s *set[T]) Remove(item T) *set[T] {
	delete(s.m, item)
	return s
}

func (s *set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *set[T]) Size() int {
	return len(s.m)
}

func (s *set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.m))
	for item := range s.m {
		slice = append(slice, item)
	}
	return slice
}

func (s *set[T]) Union(other *set[T]) *set[T] {
	union := NewSet[T]()
	for item := range s.m {
		union.Add(item)
	}
	for item := range other.m {
		union.Add(item)
	}
	return union
}

func (s *set[T]) Intersection(other *set[T]) *set[T] {
	intersection := NewSet[T]()
	for item := range s.m {
		if other.Contains(item) {
			intersection.Add(item)
		}
	}
	return intersection
}

func (s *set[T]) GetCommonElements(ar []T) []T {
	out := make([]T, 0)
	for _, item := range ar {
		if s.Contains(item) {
			out = append(out, item)
		}
	}

	return out
}
