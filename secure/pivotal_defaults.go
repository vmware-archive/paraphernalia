package secure

import "crypto/tls"

// PivotalTLSConfig creates a *tls.Config that is suitable for use in internal
// communication links between Pivotal services. It is not guaranteed to be
// suitable for communication to other external services as it contains a
// strict definition of acceptable standards. The standards were taken from the
// "Consolidated Remarks" internal document.
//
// This has yet to be audited and approved.
func PivotalTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.CurveP384,
		},
	}
}
