package intset

type IntSet struct {
	ds map[int]bool
}

func New() *IntSet {
	s := map[int]bool{}

	return &IntSet{
		ds: s,
	}
}

func (set *IntSet) Delete(v int) {
	delete(set.ds, v)
}

func (set *IntSet) Insert(v int) {
	set.ds[v] = true
}

func (set *IntSet) Len() int {
	return len(set.ds)
}

func (set *IntSet) Values(yield func(int) bool) {
	for v := range set.ds {
		if !yield(v) {
			return
		}
	}
}

func (set *IntSet) Slice() []int {
	r := []int{}

	for v := range set.Values {
		r = append(r, v)
	}

	return r
}
