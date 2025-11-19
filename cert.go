// Copyright 2021 rpeat. All rights reserved.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

package rpeat

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func CreateX509(host string, path string, isCA bool) (TLS, error) {
	var pemfiles TLS

	validFor := 365 * 24 * time.Hour
	rsaBits := 2048
	privkey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Printf("Failed to generate serial number: %v", err)
		return pemfiles, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"rpeat local"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privkey.Public(), privkey)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v", err)
		return pemfiles, err
	}
	certFile := filepath.Join(path, "cert.pem")
	certOut, err := os.Create(certFile)
	if err != nil {
		fmt.Printf("Failed to open %s for writing: %v", certFile, err)
		return pemfiles, err
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		fmt.Printf("Failed to write data to %s: %v", certFile, err)
		return pemfiles, err
	}
	if err := certOut.Close(); err != nil {
		fmt.Printf("Error closing %s: %v", certFile, err)
		return pemfiles, err
	}

	keyFile := filepath.Join(path, "key.pem")
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Failed to open %s for writing: %v", keyFile, err)
		return pemfiles, err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(privkey)
	if err != nil {
		fmt.Printf("Unable to marshal private key: %v", err)
		return pemfiles, err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		fmt.Printf("Failed to write data to %s: %v", keyFile, err)
		return pemfiles, err
	}
	if err := keyOut.Close(); err != nil {
		fmt.Printf("Error closing %s: %v", keyFile, err)
		return pemfiles, err
	}
	return TLS{Cert: certFile, Key: keyFile}, nil
}
