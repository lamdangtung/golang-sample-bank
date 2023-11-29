package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/lamdangtung/golang-sample-bank/db/mock"
	db "github.com/lamdangtung/golang-sample-bank/db/sqlc"
	"github.com/lamdangtung/golang-sample-bank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransfer(t *testing.T) {

	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()
	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.EUR
	amount := int64(10)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore, t *testing.T)
		checkResponse func(recorder *httptest.ResponseRecorder, t *testing.T)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
		{
			name: "FromAcountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "CAD",
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "InvalidAccount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.VND,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "InternalErrorWhenGetAccount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.VND,
			},
			buildStubs: func(store *mockdb.MockStore, t *testing.T) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, t *testing.T) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store, t)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/transfers")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, t)
		})
	}

}
