package htpasswd

import (
	"testing"
)

func Test_Ssha(t *testing.T) {
	testParserGood(t, "ssha", Ssha, nil, "{SSHA}/lLSOXpMWipWr3ifiighLCpqBiFoMzBM", "password")
	testParserBad(t, "ssha", nil, RejectSsha, "{SSHA}0")
	testParserNot(t, "ssha", nil, RejectSsha, "plaintext")
}
