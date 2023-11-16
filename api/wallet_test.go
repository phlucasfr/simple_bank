package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "picpay_simplificado/db/mock"
	db "picpay_simplificado/db/sqlc"
	"picpay_simplificado/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetWalletAPI(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name          string
		walletID      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), wallet.ID).
					Times(1).
					Return(wallet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWallet(t, recorder.Body, wallet)
			},
		},
		{
			name:     "NotFound",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), wallet.ID).
					Times(1).
					Return(db.Wallet{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), wallet.ID).
					Times(1).
					Return(db.Wallet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidID",
			walletID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/wallets/%d", tc.walletID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteWalletAPI(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name          string
		Owner         int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			Owner: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteWallet(gomock.Any(), wallet.ID).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			Owner: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteWallet(gomock.Any(), wallet.ID).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "InternalError",
			Owner: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteWallet(gomock.Any(), wallet.ID).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "InvalidID",
			Owner: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/wallets/%d", tc.Owner)
			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomWallet() db.Wallet {
	return db.Wallet{
		ID:       util.RandomInt(1, 100),
		Owner:    util.RandomString(10),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchWallet(t *testing.T, body *bytes.Buffer, wallet db.Wallet) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotWallet db.Wallet
	err = json.Unmarshal(data, &gotWallet)
	require.NoError(t, err)
	require.Equal(t, wallet, gotWallet)
}
