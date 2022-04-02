package taxpayer

import "testing"

func TestGetTaxpayerDataSuccess(t *testing.T) {
	taxpayer, err := GetTaxpayerData("1070023596")
	if err != nil {
		t.Fatalf("error was not expected: %v", err)
	}

	expected := "RIWO SYSTEMS SPÓŁKA Z OGRANICZONĄ ODPOWIEDZIALNOŚCIĄ"
	if taxpayer.Name != expected {
		t.Fatalf("expected name %q, but got %q", expected, taxpayer.Name)
	}

	expected = "7792465289"
	if taxpayer.Nip != expected {
		t.Fatalf("expected nip %q, but got %q", expected, taxpayer.Nip)
	}

	expected = "367435452"
	if taxpayer.Regon != expected {
		t.Fatalf("expected regon %q, but got %q", expected, taxpayer.Regon)
	}

	expectedAddress := Address{"SZAMOTULSKA 40/1A", "60-366", "POZNAŃ", "POLSKA", "PL"}
	if *taxpayer.Address != expectedAddress {
		t.Fatalf("expected address %q, but got %q", expectedAddress, taxpayer.Address)
	}
}

func TestGetTaxpayerDataBadNip(t *testing.T) {
	_, err := GetTaxpayerData("0000000000")
	if err == nil {
		t.Fatalf("error was expected")
	}
}
