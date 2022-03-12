package taxpayer

import "testing"

func TestGetTaxpayerDataSuccess(t *testing.T) {
	taxpayer, err := GetTaxpayerData("7792465289")
	if err != nil {
		t.Fatalf("error was not expected: %v", err)
	}

	expected := "RIWO SYSTEMS SPÓŁKA Z OGRANICZONĄ ODPOWIEDZIALNOŚCIĄ"
	if taxpayer.name != expected {
		t.Fatalf("expected name %q, but got %q", expected, taxpayer.name)
	}

	expected = "7792465289"
	if taxpayer.nip != expected {
		t.Fatalf("expected nip %q, but got %q", expected, taxpayer.nip)
	}

	expected = "367435452"
	if taxpayer.regon != expected {
		t.Fatalf("expected regon %q, but got %q", expected, taxpayer.regon)
	}

	expectedAddress := Address{"SZAMOTULSKA 40/1A", "60-366", "POZNAŃ"}
	if *taxpayer.address != expectedAddress {
		t.Fatalf("expected address %q, but got %q", expectedAddress, taxpayer.address)
	}

	expected = "PL"
	if taxpayer.country != expected {
		t.Fatalf("expected country %q, but got %q", expected, taxpayer.country)
	}
}

func TestGetTaxpayerDataBadNip(t *testing.T) {
	_, err := GetTaxpayerData("0000000000")
	if err == nil {
		t.Fatalf("error was expected")
	}
}
