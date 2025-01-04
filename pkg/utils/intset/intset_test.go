package intset_test

import (
	"testing"

	"github.com/reidaa/ano/pkg/utils/intset"
)

func TestInsert(t *testing.T) {
	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)

	if got, want := data.Len(), 3; got != want {
		t.Errorf("length of data = %v, want %d", got, want)
	}
}

func TestInsert2(t *testing.T) {
	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)
	data.Insert(3)

	if got, want := data.Len(), 3; got != want {
		t.Errorf("length of data = %v, want %d", got, want)
	}
}

func TestDelete(t *testing.T) {
	data := intset.New()

	data.Insert(3)
	data.Insert(5)
	data.Insert(2)
	data.Insert(3)

	data.Delete(3)

	if got, want := data.Len(), 2; got != want {
		t.Errorf("length of data = %v, want %d", got, want)
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

	if got, want := i, 3; got != want {
		t.Errorf("number of iterations = %v, want %d", got, want)
	}
}
