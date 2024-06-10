package main

import (
	"testing"
)

func TestServer(t *testing.T) {
	kv := NewKeyValeStore()

	kv.Set("name", "guido")
	expected := "guido"
	received := kv.Get("name")

	if expected != received {
		t.Errorf("Expected %q, got %q", expected, received)
	}

	received = kv.Get("Arsenal")
	expected = "Not Found"
	if expected != received {
		t.Errorf("Expected %q, got %q", expected, received)
	}

	kv.Set("x", "123")

	kv.Snapshot("kv.snapshot")
	kv2 := NewKeyValeStore()
	kv2.RetrieveSnapshot("kv.snapshot")

	received = kv2.Get("name")
	expected = "guido"
	if expected != received {
		t.Errorf("Expected %q, got %q", expected, received)
	}

	received = kv2.Get("x")
	expected = "123"
	if expected != received {
		t.Errorf("Expected %q, got %q", expected, received)
	}

}
