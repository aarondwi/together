package engine

import "testing"

func TestLinkedListBatch(t *testing.T) {
	ll := llBatch{}

	ok := ll.hasItem()
	if ok {
		t.Fatalf("It should be false, cause still no item, but it is true")
	}

	for i := 0; i < 3; i++ {
		ll.add(&Batch{})
	}

	ok = ll.hasItem()
	if !ok {
		t.Fatalf("It should be true, cause already has item, but it is false")
	}

	for i := 0; i < 3; i++ {
		b := ll.get()
		if b == nil {
			t.Fatalf("it should not be nil, cause exactly has same number of items as get, but fail at iter %d", i)
		}
	}

	b := ll.get()
	if b != nil {
		t.Fatal("it should be nil, cause already all taken, but instead we got something")
	}

	ll.add(&Batch{})
	b = ll.get()
	if b == nil {
		t.Fatal("it should not be nil, cause already put 1 more, but instead we got nil")
	}
}
