// Package htpasswd provides HTTP Basic Authentication using Apache-style htpasswd files
// for the user and password data.
//
// It supports most common hashing systems used over the decades and can be easily extended
// by the programmer to support others. (See the sha.go source file as a guide.)
//
// You will want to use something like...
//
//	myauth := htpasswd.New("./my-htpasswd-file", htpasswd.DefaultSystems, nil)
//	ok := myauth.Match(user, password)
//
// ...to use in your handler code.
// You should read about that nil, as well as Reread() too.
package htpasswd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

// An EncodedPasswd is created from the encoded password in a password file by a PasswdParser.
//
// The password files consist of lines like "user:passwd-encoding". The user part is stripped off and
// the passwd-encoding part is captured in an EncodedPasswd.
type EncodedPasswd interface {
	// MatchesPassword returns true if the string matches the password.
	// This may cache the result in the case of expensive comparison functions.
	MatchesPassword(pw string) bool
}

// PasswdParser examines an encoded password, and if it is formatted correctly and sane, return an
// EncodedPasswd which will recognize it.
//
// If the format is not understood, then return nil
// so that another parser may have a chance. If the format is understood but not sane,
// return an error to prevent other formats from possibly claiming it
//
// You may write and supply one of these functions to support a format (e.g. bcrypt) not
// already included in this package. Use sha.c as a template, it is simple but not too simple.
type PasswdParser func(pw string) (EncodedPasswd, error)

type passwdTable map[string]EncodedPasswd

// A Htpasswd encompasses an Apache-style htpasswd file for HTTP Basic authentication
type Htpasswd struct {
	filePath string
	passwds  atomic.Pointer[passwdTable]
	parsers  []PasswdParser
}

// DefaultSystems is an array of PasswdParser including all builtin parsers. Notice that Plain is last, since it accepts anything
var DefaultSystems = []PasswdParser{
	Md5,
	Sha,
	Bcrypt,
	Ssha,
	CryptSha,
	Plain,
}

type parameters struct {
	parsers []PasswdParser
}

type Option func(*parameters)

func WithParsers(parsers ...PasswdParser) Option {
	return func(p *parameters) {
		p.parsers = parsers
	}
}

// New creates an Htpasswd from an Apache-style htpasswd file for HTTP Basic Authentication.
//
// The realm is presented to the user in the login dialog.
//
// The filename must exist and be accessible to the process, as well as being a valid htpasswd file.
//
// parsers is a list of functions to handle various hashing systems. In practice you will probably
// just pass htpasswd.DefaultSystems, but you could make your own to explicitly reject some formats or
// implement your own.
//
// bad is a function, which if not nil will be called for each malformed or rejected entry in
// the password file.
func New(filename string, opts ...Option) (*Htpasswd, error) {
	params := &parameters{parsers: DefaultSystems}
	for _, opt := range opts {
		opt(params)
	}

	bf := Htpasswd{
		filePath: filename,
		parsers:  params.parsers,
	}

	if err := bf.Reload(); err != nil {
		return nil, err
	}

	return &bf, nil
}

// NewFromReader is like new but reads from r instead of a named file. Calling
// Reload on the returned Htpasswd will result in an error; use
// ReloadFromReader instead.
func NewFromReader(r io.Reader, opts ...Option) (*Htpasswd, error) {
	params := &parameters{parsers: DefaultSystems}
	for _, opt := range opts {
		opt(params)
	}

	bf := Htpasswd{
		parsers: params.parsers,
	}

	if err := bf.ReloadFromReader(r); err != nil {
		return nil, err
	}

	return &bf, nil
}

// Match checks the username and password combination to see if it represents
// a valid account from the htpassword file.
func (bf *Htpasswd) Match(username, password string) bool {
	matcher, ok := (*bf.passwds.Load())[username]

	if ok && matcher.MatchesPassword(password) {
		// we are good
		return true
	}

	return false
}

// Reload rereads the htpasswd file.
// You will need to call this to notice any changes to the password file.
// This function is thread safe. Someone versed in fsnotify might make it
// happen automatically. However, you might also connect a SIGHUP handler to
// this function.
func (bf *Htpasswd) Reload() error {
	f, err := os.Open(bf.filePath)
	if err != nil {
		return fmt.Errorf("failed to open htpasswd file %s: %w", bf.filePath, err)
	}
	defer f.Close()

	return bf.ReloadFromReader(f)
}

// ReloadFromReader is like Reload but reads credentials from r instead of a named
// file. If Htpasswd was created by New, it is okay to call Reload and
// ReloadFromReader as desired.
func (bf *Htpasswd) ReloadFromReader(r io.Reader) error {
	newPasswdMap := &passwdTable{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if perr := bf.addHtpasswdUser(newPasswdMap, line); perr != nil {
			return perr
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning htpasswd file failed: %w", err)
	}

	bf.passwds.Store(newPasswdMap)

	return nil
}

// addHtpasswdUser processes a line from an htpasswd file and add it to the user/password map. We may
// encounter some malformed lines, this will not be an error, but we will log them if
// the caller has given us a logger.
func (bf *Htpasswd) addHtpasswdUser(pwmap *passwdTable, rawLine string) error {
	// ignore empty line
	line := strings.TrimSpace(rawLine)
	if line == "" {
		return nil
	}

	// ignore comment line. Inline comments are not allowed
	if strings.HasPrefix(line, "#") {
		return nil
	}

	// split "user:encoding" at colon
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("malformed line, no colon: %s", line)
	}

	user := parts[0]
	encoding := parts[1]

	// give each parser a shot. The first one to produce a matcher wins.
	// If one produces an error then stop (to prevent Plain from catching it)
	for _, p := range bf.parsers {
		matcher, err := p(encoding)
		if err != nil {
			return err
		}
		if matcher != nil {
			(*pwmap)[user] = matcher
			return nil // we are done, we took to first match
		}
	}

	return fmt.Errorf("unable to recognize password for %s in %s", user, encoding)
}
