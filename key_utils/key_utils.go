package key_utils

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"

	"golang.org/x/crypto/ripemd160"
)

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

func GeneratePem() string {
	priv, _ := btcec.NewPrivateKey(btcec.S256())
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

func GenerateSinFromPem(pm string) string {
	key := ExtractKeyFromPem(pm)
	sin := generateSinFromKey(key)
	return sin
}

func ExtractCompressedPublicKey(pm string) string {
	key := ExtractKeyFromPem(pm)
	pub := key.PubKey()
	comp := pub.SerializeCompressed()
	hexb := hex.EncodeToString(comp)
	return hexb
}

func ExtractKeyFromPem(pm string) *btcec.PrivateKey {
	byta := []byte(pm)
	blck, _ := pem.Decode(byta)
	var ecp ecPrivateKey
	asn1.Unmarshal(blck.Bytes, &ecp)
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), ecp.PrivateKey)
	return priv
}

func ExtractSerializedKeyFromPem(pm string) string {
	priv := ExtractKeyFromPem(pm)
	ser := priv.Serialize()
	hexa := hex.EncodeToString(ser)
	return hexa
}

func GeneratePrivateKey() *btcec.PrivateKey {
	priv, _ := btcec.NewPrivateKey(btcec.S256())
	return priv
}

func Sign(message string, pm string) string {
	key := ExtractKeyFromPem(pm)
	byta := []byte(message)
	hash := sha256.New()
	hash.Write(byta)
	bytb := hash.Sum(nil)
	sig, _ := key.Sign(bytb)
	ser := sig.Serialize()
	hexa := hex.EncodeToString(ser)
	return hexa
}

func generateSinFromKey(key *btcec.PrivateKey) string {
	pub := key.PubKey()
	comp := pub.SerializeCompressed()
	hexb := hex.EncodeToString(comp)
	stx := generateSinFromPublicKey(hexb)
	return stx
}

func generateSinFromPublicKey(key string) string {
	hexa := sha256ofHexString(key)
	hexa = ripemd160ofHexString(hexa)
	versionSinTypeEtc := "0F02" + hexa
	hexa = sha256ofHexString(versionSinTypeEtc)
	hexa = sha256ofHexString(hexa)
	checksum := hexa[0:8]
	hexa = versionSinTypeEtc + checksum
	byta, _ := hex.DecodeString(hexa)
	sin := base58.Encode(byta)
	return sin
}

func sha256ofHexString(hexa string) string {
	byta, _ := hex.DecodeString(hexa)
	hash := sha256.New()
	hash.Write(byta)
	hashsum := hash.Sum(nil)
	hexb := hex.EncodeToString(hashsum)
	return hexb
}

func ripemd160ofHexString(hexa string) string {
	byta, _ := hex.DecodeString(hexa)
	hash := ripemd160.New()
	hash.Write(byta)
	hashsum := hash.Sum(nil)
	hexb := hex.EncodeToString(hashsum)
	return hexb
}
