package kvs

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	err := Put("key", "value")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestGet(t *testing.T) {
	err := Put("key", "value")
	if err != nil {
		t.Fatal(err.Error())
	}

	value, err := Get("key")
	if err != nil {
		t.Fatal(err.Error())
	}

	if value != "value" {
		t.Fatalf(`Unexpected value %s, expecting "value"`, value)
	}
}

func TestDelete(t *testing.T) {
	err := Put("key", "value")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = Delete("key")
	if err != nil {
		t.Fatal(err.Error())
	}

	value, err := Get("key")
	if err == nil || !errors.Is(err, ErrorNoSuchKey) {
		t.Fatalf(`Deleted KV ("key" => "value") left on the store: %s`, value)
	}
}
