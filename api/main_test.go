package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/thanaphat2005/simple-bank/db/sqlc"
	"github.com/thanaphat2005/simple-bank/token"
	"github.com/thanaphat2005/simple-bank/utils"
)

func newTestServer(t *testing.T, store db.Store) *Server {

	config := utils.Configuration{
		SecretKey: utils.RandomString(32),
	}

	tokenMaker, _ := token.NewPasetoMaker(config.SecretKey)

	server, err := NewServer(store, config, tokenMaker)

	require.NoError(t, err)

	return server

}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
