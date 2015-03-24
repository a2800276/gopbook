package gopbook

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type PassBookPass struct {
	//PassMetaData PassMetaData
	PassMetaData map[string]interface{} // pass.json
	Manifest     Manifest               // manifest.json
	AppleCertFn  string                 // apple cert filename
	DevCertFn    string                 // developer cert (pem, see pkcs7.go for conversion)
	DevCertPass  string                 // dev cert password
	tempDir      string                 // where temporary passes are stored for assembly
}

func NewPassBookPass(
	templateDir string, // directory containing base assets
	appleCertFn string, // apple cert
	devCertFn string, // developer cert
	devCertPass string, // password
) (pass PassBookPass, err error) {
	// copy to tmp dir
	pass.tempDir, err = ioutil.TempDir("", "gopbook")
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if fi.Name() == ".DS_Store" {
			continue
		}
		infile, err := os.Open(fmt.Sprintf("%s/%s", templateDir, fi.Name()))
		if err != nil {
			return pass, err
		}
		defer infile.Close()
		outfile, err := os.Create(fmt.Sprintf("%s/%s", pass.tempDir, fi.Name()))
		if err != nil {
			return pass, err
		}
		defer outfile.Close()
		if _, err = io.Copy(outfile, infile); err != nil {
			return pass, err
		}
	}
	// load pass.json
	passFile, err := os.Open(pass.tempDir + "/pass.json")
	if err != nil {
		return pass, fmt.Errorf("could not read pass.json: %v", err)
	}
	defer passFile.Close()

	decoder := json.NewDecoder(passFile)
	//pass.PassMetaData, err = LoadPassMetaData(passFile)
	err = decoder.Decode(&pass.PassMetaData)
	if err != nil {
		println("here")
		return
	}
	// create manifest
	pass.Manifest = NewManifest(pass.tempDir)
	pass.AppleCertFn = appleCertFn
	pass.DevCertFn = devCertFn
	pass.DevCertPass = devCertPass
	return
}

func (p *PassBookPass) finalizePass() error {
	// export the pass.json
	passJsonFile, err := os.Create(fmt.Sprintf("%s/pass.json", p.tempDir))
	if err != nil {
		return err
	}
	defer passJsonFile.Close()
	//p.PassMetaData.SavePassMeta(passJsonFile)
	encoder := json.NewEncoder(passJsonFile)
	err = encoder.Encode(p.PassMetaData)
	if err != nil {
		return err
	}
	// update manifest
	files, err := ioutil.ReadDir(p.tempDir)
	if err != nil {
		return fmt.Errorf("readdir: %v", err)
	}
	for _, fi := range files {
		if fi.Name() == "manifest.json" || fi.Name() == "signature" {
			continue
		}
		if err := p.Manifest.AddFile(fmt.Sprintf("%s/%s", p.tempDir, fi.Name())); err != nil {
			return fmt.Errorf("manifest: %v", err)
		}
	}
	manifestFn := fmt.Sprintf("%s/%s", p.tempDir, "manifest.json")
	signatureFn := fmt.Sprintf("%s/%s", p.tempDir, "signature")

	p.Manifest.WriteTo(manifestFn)
	// create signature
	return PKCS7(p.AppleCertFn, p.DevCertFn, p.DevCertPass, manifestFn, signatureFn)

}

func (p *PassBookPass) WritePassTo(writer io.Writer) error {
	if err := p.finalizePass(); err != nil {
		return fmt.Errorf("finalize: %v", err)
	}
	zipper := zip.NewWriter(writer)
	defer zipper.Close()
	files, err := ioutil.ReadDir(p.tempDir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		infile, err := os.Open(fmt.Sprintf("%s/%s", p.tempDir, fi.Name()))
		if err != nil {
			return err
		}
		defer infile.Close()
		f, err := zipper.Create(fi.Name())
		if err != nil {
			return err
		}
		_, err = io.Copy(f, infile)
		if err != nil {
			return err
		}
	}
	return nil
}
