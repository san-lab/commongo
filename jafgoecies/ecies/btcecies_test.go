package ecies

import (
	"bytes"
	"fmt"
	"testing"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v2"
)

var privKey *secp256k1.PrivateKey
var ciphertext []byte

func ECEncryptDecrypt(t *testing.T, kdf bool) {
	privKey, err := secp256k1.GeneratePrivateKey()
	fmt.Printf("Priv key %x\n", privKey.Serialize())
	if err != nil {
		t.Fatal(err)
	}
	teststring := []byte("Samo glavnoye eto dusha!")
	fmt.Println("plaintext:", teststring)
	ciphertext, err = ECEncryptPub(privKey.PubKey(), (teststring), kdf)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := ECDecryptPriv(privKey, ciphertext, kdf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("deciphered:", string(plaintext))

	if !(bytes.Compare(teststring, plaintext) == 0) {
		t.Errorf("Decoded not equal to source: %s, %s", string(teststring), string(plaintext))
	}

}

func TestECEncryptDecrypt(t *testing.T) {
	ECEncryptDecrypt(t, true)
	ECEncryptDecrypt(t, false)
}
