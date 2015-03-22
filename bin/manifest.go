package main

import "gopbook"
import "fmt"

var PassJsonSamples = []string{
	"./appleSamples/BoardingPass.pass/pass.json",
	"./appleSamples/Coupon.pass/pass.json",
	"./appleSamples/Event.pass/pass.json",
	"./appleSamples/Generic.pass/pass.json",
	"./appleSamples/StoreCard.pass/pass.json",
}

const MANIFEST_RESULT = ``

func main() {
	m := gopbook.Manifest{Prefix: "./appleSamples/"}
	for _, fn := range PassJsonSamples {
		err := m.AddFile(fn)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	str, err := m.ToJSON()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	if str != MANIFEST_RESULT {
		fmt.Printf("expected: %v, got: %v", str)
	}

}
