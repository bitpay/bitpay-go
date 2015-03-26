package key_utils_test

import (
	. "github.com/bitpay/bitpay-go/key_utils"

	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"github.com/btcsuite/btcd/btcec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("generates pem", func() {

		pem := GeneratePem()
		match, _ := regexp.MatchString("-----BEGIN EC PRIVATE KEY-----\n.*\n.*\n.*\n--", pem)
		Expect(match).To(Equal(true))
	})
	It("generates sin from pem", func() {

		pem := "-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEICg7E4NN53YkaWuAwpoqjfAofjzKI7Jq1f532dX+0O6QoAcGBSuBBAAK\noUQDQgAEjZcNa6Kdz6GQwXcUD9iJ+t1tJZCx7hpqBuJV2/IrQBfue8jh8H7Q/4vX\nfAArmNMaGotTpjdnymWlMfszzXJhlw==\n-----END EC PRIVATE KEY-----\n"
		clientId := "TeyN4LPrXiG5t2yuSamKqP3ynVk3F52iHrX"
		result := GenerateSinFromPem(pem)
		if result != clientId {
			GinkgoT().Errorf("result: %s != %s", result, clientId)
		}
	})
	It("extracts private key from pem", func() {

		keya := GeneratePrivateKey()
		pema := GeneratePemFromKey(keya)
		keyb := ExtractKeyFromPem(pema)
		hexa := hex.EncodeToString(keya.Serialize())
		hexb := hex.EncodeToString(keyb.Serialize())
		if hexa != hexb {
			GinkgoT().Errorf("expected: %s\nreceived: %s", hexa, hexb)
		}
	})

	It("signs the sha256 with a pem", func() {
		// sign the message, then extract the signature from result
		pm := GeneratePem()
		message := "Hi Everybody!"
		signed := Sign(message, pm)
		byt, _ := hex.DecodeString(signed)
		signature, _ := btcec.ParseSignature(byt, btcec.S256())

		// create the expected message
		hash := sha256.New()
		hash.Write([]byte(message))
		expectedMessage := hash.Sum(nil)

		// get the public key from the PEM
		priv := ExtractKeyFromPem(pm)
		pub := priv.PubKey()

		Expect(signature.Verify(expectedMessage, pub)).To(Equal(true))
	})
})

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

func GeneratePrivateKey() *btcec.PrivateKey {
	priv, _ := btcec.NewPrivateKey(btcec.S256())
	return priv
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}
