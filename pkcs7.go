package gopbook

// openssl pkcs12 -in signcertificate.p12 -out signcertificate.pem
// openssl smime -sign -in manifest.json -binary -out sig2 -certfile ~/Downloads/passbook/auth/wwdrcertificate.pem -inkey signcertificate.pem -signer signcertificate.pem -outform der

import (
	//"os"
	"os/exec"
)

func PKCS7(applecertfn, devcertfn, devcertpass, manifestfn, signaturefn string) error {
	cmd := exec.Command("openssl",
		"smime", "-sign", "-binary",
		"-in", manifestfn,
		"-out", signaturefn,
		"-certfile", applecertfn,
		"-inkey", devcertfn,
		"-signer", devcertfn,
		"-outform", "der",
		"-passin", "env:PKCS7PASS")
	cmd.Env = []string{"PKCS7PASS=" + devcertpass}
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	return cmd.Run()
}
