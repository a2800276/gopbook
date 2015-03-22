package gopbook

import (
	"testing"
)

const MANIFEST_RESULT = `{"BoardingPass.pass/pass.json":"6b52f6026911d42af539fe5eaab20669fed9c97d","Coupon.pass/pass.json":"ac5eccd991c295c58d7daf0675c05f67973a6321","Event.pass/pass.json":"0cb9cfac51e8ad7716476dbb69857ea88482cdf7","Generic.pass/pass.json":"4d32e183999a336a10233f87f0eb3eb907de136a","StoreCard.pass/pass.json":"9639271ddb26c9cb5baba145eb1d1ebb36a031a1"}`

func TestManifest(t *testing.T) {
	m := NewManifest("./appleSamples/")
	for _, fn := range PassJsonSamples {
		err := m.AddFile(fn)
		if err != nil {
			t.Fatalf("%v\n", err)
		}
	}
	str, err := m.ToJSON()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	if str != MANIFEST_RESULT {
		t.Fatalf("expected: %v, got: %v", MANIFEST_RESULT, str)
	}

}
