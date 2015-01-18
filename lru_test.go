package lru

import "testing"

func TestEvict(t *testing.T) {
	cache := New(2)

	cache.Add("foo1", "bar1")
	cache.Add("foo2", "bar2")
	cache.Add("foo2", "bar2")
	cache.Add("foo2", "bar2")
	cache.Add("foo3", "bar3")
	cache.Add("foo1", "bar1")

	if _, ok := cache.Get("foo2"); ok {
		t.Fatal("should evict item foo2")
	}
}
