package util

import (
	"regexp"
	"testing"
)

var wordGroupRegex = regexp.MustCompile("([a-z]+)")

func TestZeroSubgroup(t *testing.T) {
	_, err := GetFirstSubgroupMatch("", wordGroupRegex)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOneSubgroup(t *testing.T) {
	expected := "test"
	actual, err := GetFirstSubgroupMatch("#$%^test1#4", wordGroupRegex)

	if err != nil {
		t.Fatalf("error was not expected: %v", err)
	}

	if actual != expected {
		t.Fatalf("got %q, but expected %q", actual, expected)
	}
}

func TestMultipleSubgroups(t *testing.T) {
	expected := "testa"
	actual, err := GetFirstSubgroupMatch("testa1#4testb#$%^", wordGroupRegex)

	if err != nil {
		t.Fatalf("error was not expected: %v", err)
	}

	if actual != expected {
		t.Fatalf("got %q, but expected %q", actual, expected)
	}
}
