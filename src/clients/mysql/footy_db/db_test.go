package footy_db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDatabase(t *testing.T) {
	assert.Panics(t, func() {
		ConnectToDatabase("", "", "", "", "")
	})
}
