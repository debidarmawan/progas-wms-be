package config

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
)

func registerMySQLTLS() {
	// Dev Aiven (open access): DSN pakai tls=skip-verify — tidak perlu ca.pem
	if err := mysqldriver.RegisterTLSConfig("skip-verify", &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}); err != nil && !strings.Contains(err.Error(), "already registered") {
		log.Fatalf("MySQL TLS skip-verify: %v", err)
	}

	// Production: DSN pakai tls=aiven + DB_CA_CERT=/path/to/ca.pem
	caPath := os.Getenv("DB_CA_CERT")
	if caPath == "" {
		return
	}

	pem, err := os.ReadFile(caPath)
	if err != nil {
		log.Fatalf("DB_CA_CERT: cannot read %s: %v", caPath, err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(pem) {
		log.Fatal("DB_CA_CERT: failed to parse CA certificate")
	}

	err = mysqldriver.RegisterTLSConfig("aiven", &tls.Config{
		RootCAs:    rootCertPool,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil && !strings.Contains(err.Error(), "already registered") {
		log.Fatalf("DB_CA_CERT: register TLS: %v", err)
	}
}
