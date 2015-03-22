package gopbook

import (
	"os"
	"testing"
)

func TestPKCS7(t *testing.T) {
	appleFN := "./test_pkcs7/wwdrcertificate.pem"
	devFN := "./test_pkcs7/signcert.pem"

	bytes, err := MakePKCS7(appleFN, devFN, "hallo", "1234567890")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	pkcs7, err := os.Create("out.pkcs7")
	if err != nil {

		t.Fatalf("%v\n", err)
	}
	defer pkcs7.Close()
	pkcs7.Write(bytes)
}
