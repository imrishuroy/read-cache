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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/imrishuroy/read-cache/db/mock"
	db "github.com/imrishuroy/read-cache/db/sqlc"
	"github.com/imrishuroy/read-cache/util"
	"github.com/stretchr/testify/require"
)

func TestCreateCacheAPI(t *testing.T) {
	cache := randomCache(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title": cache.Title,
				"link":  cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCacheParams{
					Title: cache.Title,
					Link:  cache.Link,
				}
				store.EXPECT().
					CreateCache(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(cache, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCache(t, recorder.Body, cache)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title": cache.Title,
				"link":  cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCache(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Cache{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: gin.H{

				"link": cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCache(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// create request
			url := "/caches"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

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
			server := newTestServer(t, store)
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

func TestListCachesAPI(t *testing.T) {
	n := 5
	caches := make([]db.Cache, n)
	for i := 0; i < n; i++ {
		caches[i] = randomCache(t)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListCachesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListCaches(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(caches, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCaches(t, recorder.Body, caches)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCaches(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Cache{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCaches(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/caches"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCacheAPI(t *testing.T) {
	cache := randomCache(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"Id":    cache.ID,
				"title": "This is test title",
				"link":  cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCacheParams{
					ID:    cache.ID,
					Title: "This is test title",
					Link:  cache.Link,
				}

				store.EXPECT().
					UpdateCache(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(cache, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCache(t, recorder.Body, cache)

			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"Id":    cache.ID,
				"title": "This is test title",
				"link":  cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateCache(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Cache{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name: "InvalidTitle",
			body: gin.H{
				"Id": cache.ID,
				// "title": "This is test title",
				"link": cache.Link,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateCache(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// create request
			url := "/caches"
			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestDeleteCacheAPI(t *testing.T) {
	cache := randomCache(t)

	testCases := []struct {
		name          string
		ID            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   cache.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCache(gomock.Any(), gomock.Eq(cache.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name: "InternalError",
			ID:   cache.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCache(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCache(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// create request
			url := fmt.Sprintf("/caches/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			// execute request
			server.router.ServeHTTP(recorder, request)
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

func requireBodyMatchCaches(t *testing.T, body *bytes.Buffer, caches []db.Cache) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getCaches []db.Cache
	err = json.Unmarshal(data, &getCaches)
	require.NoError(t, err)
	require.Equal(t, caches, getCaches)
}
