package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNonceGenerator(t *testing.T) {
	nonce := defaultNonce()
	assert.Len(t, nonce, 36)
	assert.Regexp(t, "^[[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}$", nonce)
}

func TestDefaultTimestamp(t *testing.T) {
	ts := timestamp()
	assert.Regexp(t, `^\d{13}$`, ts)
}
