package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	err := ConnectDatabase()
	assert.NoError(t, err, "ConnectDatabase() should not return an error")
	assert.NotNil(t, DB, "DB should not be nil after ConnectDatabase()")

	hasUsers := DB.Migrator().HasTable("users")
	hasMessages := DB.Migrator().HasTable("messages")

	assert.True(t, hasUsers, "Table 'users' should exist after migration")
	assert.True(t, hasMessages, "Table 'messages' should exist after migration")
}
