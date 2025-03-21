package bucket

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBucket_UploadFromURL(t *testing.T) {
	url := "https://static.ghost.org/v4.0.0/images/feature-image.jpg"
	ctx := context.Background()
	man, err := NewManager(map[string]Option{
		"default": {
			Provider:   "local",
			ExposedURL: "https://fs.weflux.cn",
			Local:      &LocalOption{Dir: "/tmp/uploads"},
		},
	})
	require.NoError(t, err)
	bucket, err := man.Get("default")
	require.NoError(t, err)
	err := bucket.UploadFromURL(ctx, "hello/test.jpg", url)
	require.NoError(t, err)
}
