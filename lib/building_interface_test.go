package cookie_clicker

import (
	"testing"
)

func TestMakeBuilding(t *testing.T) {
	name := "some-building"
	cps := 1.0

	b := newStandardBuilding(
		name,
		"",
		nil,
		cps,
	)

	if b.GetName() != name {
		t.Errorf("Unexpected name: %s != %s", b.GetName(), name)
	}

	if b.GetCPS() != cps {
		t.Errorf("Unexpected CPS: %e != %e", b.GetCPS(), cps)
	}
}
