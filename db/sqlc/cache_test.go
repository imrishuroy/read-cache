package db

import (
	"context"
	"testing"
	"time"

	"github.com/imrishuroy/read-cache-api/util"

	"github.com/stretchr/testify/require"
)

func createRandomCache(t *testing.T) Cache {
	arg := CreateCacheParams{
		Title: util.RandomTitle(),
		Link:  util.RandomLink(),
	}

	cache, err := testQueries.CreateCache(context.Background(), arg)
	require.NoError(t, err) // check that error must be nill
	require.NotEmpty(t, cache)

	require.Equal(t, arg.Title, cache.Title)
	require.Equal(t, arg.Link, cache.Link)

	require.NotZero(t, cache.ID)
	require.NotZero(t, cache.CreatedAt)

	return cache
}

func TestCreateCache(t *testing.T) {
	createRandomCache(t)
}

func TestGetCache(t *testing.T) {
	cache1 := createRandomCache(t)
	cache2, err := testQueries.GetCache(context.Background(), cache1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cache2)

	require.Equal(t, cache1.ID, cache2.ID)
	require.Equal(t, cache1.Title, cache2.Title)
	require.Equal(t, cache1.Link, cache2.Link)
	require.WithinDuration(t, cache1.CreatedAt, cache2.CreatedAt, time.Second)
}

func TestUpdateCache(t *testing.T) {
	cache1 := createRandomCache(t)

	arg := UpdateCacheParams{
		ID:    cache1.ID,
		Title: cache1.Title,
		Link:  util.RandomLink(),
	}

	cache2, err := testQueries.UpdateCache(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cache2)

	require.Equal(t, cache1.ID, cache2.ID)
	require.Equal(t, cache1.Title, cache2.Title)
	require.Equal(t, arg.Link, cache2.Link)
	require.WithinDuration(t, cache1.CreatedAt, cache2.CreatedAt, time.Second)
}

func TestDeleteCache(t *testing.T) {
	cache1 := createRandomCache(t)
	err := testQueries.DeleteCache(context.Background(), cache1.ID)
	require.NoError(t, err)

	cache2, err := testQueries.GetCache(context.Background(), cache1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, cache2)
}

func TestListCache(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomCache(t)
	}

	arg := ListCachesParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListCaches(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
