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
	mockdb "github.com/imrishuroy/read-cache/db/mock"
	db "github.com/imrishuroy/read-cache/db/sqlc"
	"github.com/imrishuroy/read-cache/util"
	"github.com/stretchr/testify/require"
)

func TestGetCache(t *testing.T) {
	cache := randomCache(t)

	// anonymous struct
	testCases := []struct {
		name          string
		ID            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   cache.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCache(gomock.Any(), gomock.Eq(cache.ID)).
					Times(1).
					Return(cache, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
				requireBodyMatchCache(t, recorder.Body, cache)
			},
		},
		{
			name: "NotFound",
			ID:   cache.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCache(gomock.Any(), gomock.Eq(cache.ID)).
					Times(1).
					Return(db.Cache{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name: "InternalError",
			ID:   cache.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCache(gomock.Any(), gomock.Eq(cache.ID)).
					Times(1).
					Return(db.Cache{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name: "InvalidID",
			ID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCache(gomock.Any(), gomock.Any()).
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
			// create mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create mock store
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// create server
			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			// create request
			url := fmt.Sprintf("/caches/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// execute request
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}

}

func randomCache(t *testing.T) db.Cache {
	return db.Cache{
		ID:    util.RandomInt(1, 1000),
		Title: util.RandomString(6),
		Link:  util.RandomString(6),
	}
}

func requireBodyMatchCache(t *testing.T, body *bytes.Buffer, cache db.Cache) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getCache db.Cache
	err = json.Unmarshal(data, &getCache)
	require.NoError(t, err)
	require.Equal(t, cache, getCache)
}
