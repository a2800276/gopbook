package gopbook

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Manifest struct {
	Prefix string
	files  map[string]string
}

func NewManifest(prefix string) Manifest {
	return Manifest{prefix, make(map[string]string)}
}

func (m *Manifest) AddFile(fn string) error {
	if strings.Index(fn, m.Prefix) != 0 {
		return fmt.Errorf("unexpected prefix: %s should start with %s", fn, m.Prefix)
	}
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	sha := sha1.New()
	io.Copy(sha, file)
	shaHex := fmt.Sprintf("%x", sha.Sum(nil))
	name := fn[len(m.Prefix):]
	m.files[name] = shaHex
	return nil
}

func (m *Manifest) ToJSON() (string, error) {
	str, err := json.Marshal(m.files)
	return string(str), err
}
