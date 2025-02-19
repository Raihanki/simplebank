package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/raihanki/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(10)
	expiredAt := time.Until(time.Now().Add(time.Hour * 24))
	token, err := pasetoMaker.CreateToken(username, expiredAt)
	require.NoError(t, err)

	payload, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(expiredAt), time.Second)
}

func TestExpiredToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(10)
	token, err := pasetoMaker.CreateToken(username, -time.Second)
	require.NoError(t, err)

	payload, err := pasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, fmt.Errorf("this token has expired").Error())
}

func TestInvalidSymentricKey(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(30))
	require.Error(t, err)
	require.Nil(t, pasetoMaker)
}
