package htpasswd

import (
	"testing"
)

func Test_Plain(t *testing.T) {
	testParserGood(t, "plain", Plain, RejectPlain, "bar", "bar")
	// testParserBad() plain takes anything
	// testParserNot() plain takes anything
}
