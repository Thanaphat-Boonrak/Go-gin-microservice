package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	mockdb "github.com/thanaphat2005/simple-bank/db/mock"
	db "github.com/thanaphat2005/simple-bank/db/sqlc"
	"github.com/thanaphat2005/simple-bank/utils"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUsersParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUsersParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("match arg %v and password", e.arg)
}

func EqCreateUsersParams(arg db.CreateUsersParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUsersParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().
					CreateUsers(gomock.Any(), EqCreateUsersParams(arg, password)).
					Return(user, nil).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUsers(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DulpicatedUsername",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUsers(gomock.Any(), gomock.Any()).
					Return(db.User{}, &pgconn.PgError{Code: "23505"}).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
				"username":  "invalid#Username",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email":     "Thanaphat",
				"password":  password,
				"full_name": user.FullName,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"email":     user.Email,
				"password":  "123",
				"full_name": user.FullName,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandonEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)

	require.NoError(t, err)
	var gotAccount db.User
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, user, gotAccount)
}
