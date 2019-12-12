package hiper

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func Login(username, password string) (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %w", err)
	}

	s := Session{
		Client: http.Client{
			Jar: jar,
		},
	}

	values := url.Values{}
	values.Add("email", username)
	values.Add("password", password)

	resp, err := s.PostForm("https://www.hiper.dk/mit-hiper/login/registrer-dynamisk", values)
	if err != nil {
		return nil, fmt.Errorf("unable to post login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected login http statuscode: %d", resp.StatusCode)
	}
	return &s, nil
}

// Session represents a login session with hiper
type Session struct {
	http.Client
}

func (s *Session) Invoices() (Invoices, error) {
	resp, err := s.Get("https://www.hiper.dk/mit-hiper/regninger-dynamisk")
	if err != nil {
		return nil, fmt.Errorf("unable to get invoices: %s")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected get invoices http status code: %d", resp.StatusCode)
	}

	invoices := make(Invoices, 0)
	err = invoices.FromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot blarh: %w", err)
	}

	return invoices, nil
}

func (s *Session) Logout() error {
	resp, err := s.Get("https://www.hiper.dk/mit-hiper/logud")
	if err != nil {
		return fmt.Errorf("unable to logout: %w", err)
	}

	// we dont need the body
	resp.Body.Close()

	// we also dont really care about the http status code
	return nil
}
