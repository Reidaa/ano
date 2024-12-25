package intset_test

import (
	"testing"

	"github.com/reidaa/ano/pkg/utils/intset"
)

func TestInsert(t *testing.T) {
	set := intset.New()

	set.Insert(3)
	set.Insert(5)
	set.Insert(2)

	if set.Len() != 3 {
		t.Errorf("length of data is wrong")
	}
}

func TestInsert2(t *testing.T) {
	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)
	data.Insert(3)

	if data.Len() != 3 {
		t.Errorf("length of data is wrong")
	}
}

func TestDelete(t *testing.T) {
	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)
	data.Insert(3)

	data.Delete(3)

	if data.Len() != 2 {
		t.Errorf("length of data is wrong")
	}
}

func TestIterate(t *testing.T) {
	var i int

	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)
	data.Insert(3)

	for range data.Values {
		i++
	}

	if i != 3 {
		t.Errorf("iteration failed")
	}
}
