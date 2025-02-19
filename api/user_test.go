package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/raihanki/simplebank/db/mock"
	db "github.com/raihanki/simplebank/db/sqlc"
	"github.com/raihanki/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user := randomUser(t)

	testCases := []struct {
		name      string
		body      map[string]interface{}
		buildStub func(store *mockdb.MockStore)
		checkResp func(reccorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  util.RandomString(8),
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(user, nil)
			},
			checkResp: func(reccorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, reccorder.Code)
				requireBodyMatchUser(t, reccorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStub(mockStore)

			recorder := httptest.NewRecorder()
			server := newTestServer(t, mockStore)

			url := "/users"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			tc.checkResp(recorder)
		})
	}
}

func randomUser(t *testing.T) db.User {
	password, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	var ts pgtype.Timestamptz
	ts.Time = time.Now()

	user := db.User{
		Username:  util.RandomString(10),
		Password:  password,
		FullName:  util.RandomString(6) + " " + util.RandomString(5),
		Email:     util.RandomString(6) + "@simplebank.com",
		CreatedAt: ts,
	}

	return user
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
}
