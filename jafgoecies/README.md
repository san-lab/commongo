Just Another Go ECIES implementation
- based on btcec

ECIES is using Diffie-Hellman to derive a symmetric key

The main methods are `ECEncryptPub(...)` and `ECDecryptPriv(...)`

As the name suggests, the first one encrypts a message. 
An ephemeral key pair is generated, and a secret is derived (D-H) from the public key provided and the ephemeral key.
The boolean flag drives if a kdf function will be used on the derived secret (recommended, but may get in the way of interoperability).
The message is AES256-encrypted
The ephemeral public key is prepended to the encrypted message, so the secret can be re-generated if one has access to the private key.

```
func ECEncryptPub(pubkey *btcec.PublicKey, msg []byte, kdf bool) ([]byte, error) {...}
func ECDecryptPriv(privkey *btcec.PrivateKey, msg []byte, kdf bool) ([]byte, error) {...}
```
The output/input has the format: 
1) 65 bytes of btcec-encoded ephemeral public key
2) AES encrypted message with 32-bytes block and 12-bytes nonce (nonce comes first)
  
