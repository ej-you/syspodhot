package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	password := []byte("Qwerty123!")

	hash, err := Encode(password)
	require.NoError(t, err, "encode password")

	t.Logf("Hash: %s", hash)
}
