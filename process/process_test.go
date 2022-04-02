package process

import (
	"mrsydar/tkl/k360/client"
	"testing"
)

func TestProcess(t *testing.T) {
	client := client.New(
		"816ff0fc-2ecb-4fe2-ae5f-9d60cfce338a",
		"w03NSjlUVgsJ/8Dy7lpmNnSYOc2RhieBa0uw2rAYRq0=",
	)

	err := ProcessInvoices(*client, "/home/mrsydar/Desktop/TKL EXPORT.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Fail()
}
