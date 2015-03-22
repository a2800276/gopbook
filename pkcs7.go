package gopbook

// openssl smime -sign -in manifest.json -binary -out sig2 -certfile ~/Downloads/passbook/auth/wwdrcertificate.pem -inkey signcertificate.pem -signer signcertificate.pem -outform der

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

var pkcs7_signed_data_oid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
var pkcs7_data_oid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
var eContenttype = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 4}
var sha1oid = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
var contentTypeOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
var signingTimeOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
var messageDigestOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
var rsaOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
var sha1withRSAOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 5}

type PKCS7 struct {
	ContentType asn1.ObjectIdentifier
	Content     []interface{}
}

// convert pkcs12 using:
// openssl pkcs12 -in signcertificate.p12 -out signcert.pem
func MakePKCS7(appleCertFn, devCertFn, manifest, pass string) ([]byte, error) {
	appleCertFile, err := os.Open(appleCertFn)
	if err != nil {
		return nil, err
	}
	defer appleCertFile.Close()
	data, err := ioutil.ReadAll(appleCertFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)

	appleCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	devCertFile, err := os.Open(devCertFn)
	if err != nil {
		return nil, err
	}
	defer devCertFile.Close()
	data, err = ioutil.ReadAll(devCertFile)

	var privKey *rsa.PrivateKey
	var devCert *x509.Certificate

	for block, data = pem.Decode(data); block != nil; block, data = pem.Decode(data) {
		if block.Type == "CERTIFICATE" {
			devCert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
		if block.Type == "ENCRYPTED PRIVATE KEY" {
			bytes, err := x509.DecryptPEMBlock(block, []byte(pass))
			if err != nil {
				return nil, err
			}
			privKey, err = x509.ParsePKCS1PrivateKey(bytes)
			if err != nil {
				return nil, err
			}
		}
		if block.Type == "PRIVATE KEY" {
			pKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			privKey = pKey.(*rsa.PrivateKey)
		}
	}
	// check everything loaded
	return makePKCS7(appleCert, devCert, privKey, manifest)

}

func makePKCS7(appleCert, devCert *x509.Certificate, priv *rsa.PrivateKey, manifest string) ([]byte, error) {
	var pkcs7 PKCS7
	pkcs7.ContentType = pkcs7_signed_data_oid

	var signedData SignedData
	signedData.Version = 1
	signedData.DigestAlgorithms = make([]AlgorithmIdentifier, 1, 1)
	signedData.DigestAlgorithms[0] = newAlgorithmIdentifierWithNull(sha1oid)
	signedData.ContentInfo = ContentInfo{pkcs7_data_oid, asn1.RawValue{Tag: 5}}
	signedData.Certificates = make([]Certificate, 2, 2)
	signedData.Certificates[0] = copyCert(appleCert)
	signedData.Certificates[1] = copyCert(devCert)

	var signerInfo SignerInfo
	signerInfo.Version = 1
	signerInfo.IssuerAndSerialNumber = newIssuerAndSerialNumber(devCert)
	signerInfo.DigestAlgorithm = newAlgorithmIdentifierWithNull(sha1oid)
	// digest = calc sha1 over manifest
	digest := sha1.Sum([]byte(manifest))
	signerInfo.AuthenticatedAttributes = newAuthenticatedAttributes(digest[0:])
	signerInfo.DigestEncryptionAlgorithm = newAlgorithmIdentifierWithNull(rsaOID)
	log.Printf("here")
	authAttrDER, err := asn1.Marshal(signerInfo.AuthenticatedAttributes)
	if err != nil {
		log.Printf("here")
		return nil, err
	}
	authAttrDigest := sha1.Sum(authAttrDER)
	encDigest, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, authAttrDigest[0:])

	// encDigest = authenticated_attributes? in der
	signerInfo.EncryptedDigest = encDigest
	signedData.SignerInfo = make([]SignerInfo, 1, 1)
	signedData.SignerInfo[0] = signerInfo
	pkcs7.Content = make([]interface{}, 1, 1)
	pkcs7.Content[0] = signedData
	return asn1.Marshal(pkcs7)
}

type SignedData struct {
	Version          int
	DigestAlgorithms []AlgorithmIdentifier `asn1:"set"`
	ContentInfo      ContentInfo
	Certificates     []Certificate
	CRLS             []pkix.CertificateList
	SignerInfo       []SignerInfo `asn1:"set"`
}

type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue
}

func newAlgorithmIdentifierWithNull(alg asn1.ObjectIdentifier) AlgorithmIdentifier {
	var ai AlgorithmIdentifier
	ai.Algorithm = alg
	ai.Parameters = asn1.RawValue{Tag: 5}
	return ai
}

type ContentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     interface{}
}

type SignerInfo struct {
	Version                   int
	IssuerAndSerialNumber     IssuerAndSerialNumber
	DigestAlgorithm           AlgorithmIdentifier
	AuthenticatedAttributes   AuthenticatedAttributes
	DigestEncryptionAlgorithm AlgorithmIdentifier
	EncryptedDigest           []byte
	UnauthenticatedAttributes []byte `asn1:"optional,omitempty`
}

type IssuerAndSerialNumber struct {
	Issuer       pkix.Name
	SerialNumber *big.Int
}

func newIssuerAndSerialNumber(cert *x509.Certificate) IssuerAndSerialNumber {
	var ias IssuerAndSerialNumber
	ias.Issuer = cert.Issuer
	ias.SerialNumber = cert.SerialNumber
	return ias
}

type AuthenticatedAttributes struct {
	ContentType   ContentType
	SigningTime   SigningTime
	MessageDigest MessageDigest
}

func newAuthenticatedAttributes(digest []byte) AuthenticatedAttributes {
	var aa AuthenticatedAttributes
	aa.ContentType = newContentType()
	aa.SigningTime = newSigningTime()
	aa.MessageDigest = newMessageDigest(digest)
	return aa
}

type ContentType struct {
	OID     asn1.ObjectIdentifier
	TypeSet []asn1.ObjectIdentifier
}

func newContentType() ContentType {
	var ct ContentType
	ct.OID = contentTypeOID
	ct.TypeSet = make([]asn1.ObjectIdentifier, 1, 1)
	ct.TypeSet[0] = pkcs7_data_oid
	return ct
}

type SigningTime struct {
	OID     asn1.ObjectIdentifier
	TimeSet []time.Time
}

func newSigningTime() SigningTime {
	var st SigningTime
	st.OID = signingTimeOID
	st.TimeSet = make([]time.Time, 1, 1)
	st.TimeSet[0] = time.Now()
	return st
}

type MessageDigest struct {
	OID   asn1.ObjectIdentifier
	MDSet [][]byte
}

func newMessageDigest(digest []byte) MessageDigest {
	var md MessageDigest
	md.OID = messageDigestOID
	md.MDSet = make([][]byte, 1, 1)
	md.MDSet[0] = digest
	return md
}

type publicKeyInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type validity struct {
	NotBefore, NotAfter time.Time
}

type TBSCertificate struct {
	Version            int `asn1:"optional,explicit,default:1,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueId           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueId    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
}

type Certificate struct {
	TBSCertificate     TBSCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

func copyCert(c *x509.Certificate) Certificate {
	var cert Certificate
	_, err := asn1.Unmarshal(c.Raw, &cert)
	if err != nil {
		panic(err.Error())
	}
	return cert
}
