package gopbook

import (
	"fmt"
	"os"
	"testing"
)

func TestPassbook0(t *testing.T) {
	pass, err := NewPassBookPass(
		"./appleSamples/static.pass",
		appleFN,
		devFN,
		passwd,
	)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	f, err := os.Create("./appleSamples/static2000.pkpass")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	err = pass.WritePassTo(f)
	if err != nil {
		t.Fatalf("Write to pass: %v\n", err)
	}
}

func TestPassbook(t *testing.T) {

	passes := make([]PassBookPass, 0, 1)

	for _, passdir := range PassSamples {
		pass, err := NewPassBookPass(
			passdir,
			appleFN,
			devFN,
			passwd,
		)
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		passes = append(passes, pass)
	}
	t.Logf("%v", len(passes))
	//t.Logf("%v", passes[0])
	for i, pass := range passes {
		fn := fmt.Sprintf("./appleSamples/pass_%d.pkpass", i)
		f, err := os.Create(fn)
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		err = pass.WritePassTo(f)
		if err != nil {
			t.Fatalf("Write to pass: %v\n", err)
		}
	}
}
