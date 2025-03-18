package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/singhJasvinder101/go_bank/db/mock"
	db "github.com/singhJasvinder101/go_bank/db/sqlc"
	"github.com/singhJasvinder101/go_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestGETAccountAPI(t *testing.T) {
	user := utils.RandomString()

	account := randomAccount(user)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start test server
	server := NewServer(store)
	// response recorder
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request := httptest.NewRequest("GET", url, nil)

	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, 200, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func TestGETAccountWithtestCasessAPI(t *testing.T) {
	user := utils.RandomString()

	account := randomAccount(user)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
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
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			testCase.buildStubs(store)

			// Start test server and server request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       int64(utils.RandomInt(1, 10000)),
		Owner:    owner,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, boddy *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(boddy)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
