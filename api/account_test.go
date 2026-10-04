package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	mockdb "github.com/thanaphat2005/simple-bank/db/mock"
	db "github.com/thanaphat2005/simple-bank/db/sqlc"
	"github.com/thanaphat2005/simple-bank/token"
	"github.com/thanaphat2005/simple-bank/utils"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct {
		name          string
		accountID     int64
		setup         func(t *testing.T, req *http.Request, token token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:      "Ok",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(account, nil).Times(1)
		},
		setup: func(t *testing.T, req *http.Request, token token.Maker) {
			addAuthorization(t, req, token, "Bearer", user.Username, time.Minute)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, account)
		},
	},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(db.Account{}, pgx.ErrNoRows).Times(1)
			},
			setup: func(t *testing.T, req *http.Request, token token.Maker) {
				addAuthorization(t, req, token, "Bearer", user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(db.Account{}, sql.ErrConnDone).Times(1)
			},
			setup: func(t *testing.T, req *http.Request, token token.Maker) {
				addAuthorization(t, req, token, "Bearer", user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			url := fmt.Sprintf("/accounts/%d", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			tc.setup(t, request, server.maker)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)

	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
