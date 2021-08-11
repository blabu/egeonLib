package golang

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
)

// encodePemBlock - encode pem bloc that crypted with password and retrun encrypted data
func encodePemBlock(data []byte, password string) ([]byte, error) {
	keyBlock, _ := pem.Decode(data)
	var keyDER []byte
	var err error
	if x509.IsEncryptedPEMBlock(keyBlock) {
		keyDER, err = x509.DecryptPEMBlock(keyBlock, []byte(password))
		if err != nil {
			return nil, err
		}
	} else {
		return pem.EncodeToMemory(keyBlock), nil
	}
	// Update keyBlock with the plaintext bytes and clear the now obsolete
	// headers.
	keyBlock.Bytes = keyDER
	keyBlock.Headers = nil

	// Turn the key back into PEM format so we can leverage tls.X509KeyPair,
	// which will deal with the intricacies of error handling, different key
	// types, certificate chains, etc.
	keyPEM := pem.EncodeToMemory(keyBlock)
	return keyPEM, nil
}

//GetEncryptedCert return tls.Certificate encrypted with password
func GetEncryptedCert(privateKeyPath, certPath, password string) (tls.Certificate, error) {
	if len(password) == 0 {
		return tls.LoadX509KeyPair(certPath, privateKeyPath)
	}
	key, err := ReadFile(privateKeyPath)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPEM, err := encodePemBlock(key, password)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM, err := ReadFile(certPath)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.X509KeyPair(certPEM, keyPEM)
}
