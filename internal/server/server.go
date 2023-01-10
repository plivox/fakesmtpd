package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fakesmtpd/internal/config"
	"io/ioutil"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// Error logger use zerolog
type ErrorLogger struct{}

func (l *ErrorLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}

func (l *ErrorLogger) Println(v ...interface{}) {
	log.Print(v)
}

func loadRootCAs(filepath string) (roots *x509.CertPool, err error) {
	rootPEM, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	roots = x509.NewCertPool()

	if ok := roots.AppendCertsFromPEM(rootPEM); !ok {
		err = errors.New("failed to parse root certificate")
	}
	return
}

func NewServer(config *config.Config) error {
	var (
		TLSConfig = &tls.Config{}
		err       error
	)

	if config.Server.TLS && config.Server.CACert != "" {
		if TLSConfig.RootCAs, err = loadRootCAs(config.Server.CACert); err != nil {
			return err
		}
	} else if config.Server.Insecure {
		TLSConfig.InsecureSkipVerify = true
	}

	s := smtp.NewServer(&Backend{})
	// s.AuthDisabled
	s.Addr = config.Server.Address
	s.Domain = config.Server.Domain
	s.TLSConfig = TLSConfig
	s.ErrorLog = &ErrorLogger{}
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Info().Msgf("Starting server at %s", config.Server.Address)

	if config.Server.TLS {
		err = s.ListenAndServeTLS()
	} else {
		err = s.ListenAndServe()
	}
	return err
}
