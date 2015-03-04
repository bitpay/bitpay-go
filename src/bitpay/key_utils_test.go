package bitpay

import (
	"crypto/elliptic"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"regexp"
	"testing"
)

func TestGeneratePem(t *testing.T) {
	pem := GeneratePem()
	match, _ := regexp.MatchString("-----BEGIN EC PRIVATE KEY-----\n.*\n.*\n.*\n--", pem)
	if !match {
		t.Errorf(pem)
	}
}

func TestGenerateSinFromPublicKey(t *testing.T) {
	key := "031B36F2A119CBEDE86731403B1D5FCC3DCC48C220F0A17F903336587C6179527E"
	result := GenerateSinFromPublicKey(key)
	if result != "TfEn9UoPraR2iu746HreEXfqKBasFBm3dxw" {
		t.Errorf(result)
	}
}

func TestExtractPrivateKeyFromPem(t *testing.T) {
	keya := GeneratePrivateKey()
	pema := GeneratePemFromKey(keya)
	keyb := ExtractKeyFromPem(pema)
	if fmt.Sprintf("%x", keya.Serialize()) != fmt.Sprintf("%x", keyb.Serialize()) {
		t.Error(fmt.Sprintf("%x", keya.Serialize()), fmt.Sprintf("%x", keya.Serialize()))
	}
}

func GeneratePemFromKey(priv *btcec.PrivateKey) string {
	pub := priv.PubKey()
	ecd := pub.ToECDSA()
	oid := asn1.ObjectIdentifier{1, 3, 132, 0, 10}
	curve := btcec.S256()
	der, _ := asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    priv.D.Bytes(),
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(curve, ecd.X, ecd.Y)},
	})
	blck := pem.Block{Type: "EC PRIVATE KEY", Bytes: der}
	pm := pem.EncodeToMemory(&blck)
	return string(pm)
}
