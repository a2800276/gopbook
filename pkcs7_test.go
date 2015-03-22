package gopbook

import (
	"testing"
)

const appleFN = "./test_pkcs7/wwdrcertificate.pem"
const devFN = "./test_pkcs7/signcertificate.pem"
const manifestFN = "./test_pkcs7/manifest.json"
const outFN = "./test_pkcs7/signature"
const passwd = "1234567890"

func TestPKCS7(t *testing.T) {

	// TODO check test dir exists and warn.

	err := PKCS7(appleFN, devFN, passwd, manifestFN, outFN)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

}
