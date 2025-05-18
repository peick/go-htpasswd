# htpasswd for Go

![](https://github.com/peick/go-htpasswd/workflows/Go/badge.svg)
[![Go Doc](https://godoc.org/github.com/peick/go-htpasswd?status.svg)](https://godoc.org/github.com/peick/go-htpasswd)
[![Go Report Card](https://goreportcard.com/badge/github.com/peick/go-htpasswd)](https://goreportcard.com/report/github.com/peick/go-htpasswd)

This is a go libary to validate user credentials against an HTPasswd file as used in apache and nginx.

## Currently, this supports:

* SSHA
* MD5Crypt
* APR1Crypt
* SHA
* Bcrypt
* Plain text
* Crypt with SHA-256 and SHA-512

## Thanks to

This library was forked from <https://github.com/jimstudt/http-authentication/tree/master/basic>
with modifications support SSHA, Md5Crypt, Bcrypt, Crypt with SHA-256 and SHA-512 support.
