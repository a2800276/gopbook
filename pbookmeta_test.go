package gopbook

import (
	"os"
	"testing"
	"time"
)

var PassSamples = []string{
	"./appleSamples/BoardingPass.pass",
	"./appleSamples/Coupon.pass",
	"./appleSamples/Event.pass",
	"./appleSamples/Generic.pass",
	"./appleSamples/StoreCard.pass",
}

const InstallSamplesMes = `
Could not find Apple sample Passbook file in 'appleSamples', these
are not distributed along with this package for license reasons. Please
download them from here: https://developer.apple.com/downloads/index.action?name=Passbook
and copy the 'Sample Passes' directory to 'appleSamples'. Could not find %s`

func TestLoad(t *testing.T) {
	for _, fn := range PassSamples {
		file, err := os.Open(fn + "/pass.json")
		if err != nil {
			t.Fatalf(InstallSamplesMes, fn)
			return
		}
		defer file.Close()
		pb, err := LoadPassMetaData(file)
		if err != nil {
			t.Errorf("Could not load: %s, %v\n", fn, err)
			continue
		}
		//t.Logf("%v", pb)
		void(pb)
	}

}

func void(o interface{}) {}

func TestTime(t *testing.T) {
	format := "2006-01-02T15:04-07:00"
	t0 := "2012-07-22T14:25-08:00"
	if ti, err := time.Parse(format, t0); err != nil {
		t.Errorf("%v, %v, %v\n", format, t0, err)
	} else {
		t.Logf("%v, %v, %v\n", format, t0, ti)

	}
}
