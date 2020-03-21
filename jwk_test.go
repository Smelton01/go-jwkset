package jwkset

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
)

func TestFetcher(t *testing.T) {
	assert := assert.New(t)
	fetcher := &HTTPFetcher{
		Client: &http.Client{},
	}
	jwksresp, err := fetcher.FetchJWKs("https://www.googleapis.com/oauth2/v3/certs")
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 3)
}

func TestJWKsCacher(t *testing.T) {
	defaultExpiration := 10 * time.Minute
	cleanupInterval := time.Minute

	assert := assert.New(t)
	cacher := NewCacher(defaultExpiration, cleanupInterval, &HTTPFetcher{
		Client: &http.Client{},
	})

	cacheKey := "https://www.googleapis.com/oauth2/v3/certs"
	jwksresp, err := cacher.FetchJWKs(cacheKey)
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 3)

	cachedResp, found := cacher.cache.Get(cacheKey)
	assert.True(found)

	resp, ok := cachedResp.([]jose.JSONWebKey)
	if assert.True(ok, "cached response should be []jose.JSONWebKey but %#v", cachedResp) {
		assert.Equal(jwksresp.Keys, resp)
	}

	jwksresp, err = cacher.FetchJWKs(cacheKey)
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 3)
}
