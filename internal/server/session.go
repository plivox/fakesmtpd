package server

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Info().Str("from", from).Msg("")
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Info().Str("rcpt-to", to).Msg("")
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Info().Str("data", string(b)).Msg("")
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
