package bucket

import (
	"context"
	"errors"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"net/http"
	"net/url"
)

type Bucket struct {
	*blob.Bucket
	name       string
	exposedURL string
}

func (b *Bucket) ExposeURL(file string) string {
	u, _ := url.JoinPath(b.exposedURL, b.name, file)
	return u
}

func (b *Bucket) UploadFromURL(ctx context.Context, path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := b.Bucket.Upload(ctx, path, resp.Body, &blob.WriterOptions{
		ContentType: resp.Header.Get("Content-Type"),
	}); err != nil {
		return err
	}
	return nil
}

type Option struct {
	Provider   string       `json:"provider"`
	Local      *LocalOption `json:"local"`
	ExposedURL string       `json:"exposed_url"`
}

type LocalOption struct {
	Dir string `json:"dir"`
}

type Manager struct {
	buckets map[string]*Bucket
}

func NewManager(opts map[string]Option) (*Manager, error) {
	buckets := make(map[string]*Bucket)
	for k, opt := range opts {
		switch opt.Provider {
		case "local":
			// Create a file-based bucket.
			bucket, err := fileblob.OpenBucket(opt.Local.Dir, &fileblob.Options{
				CreateDir:   true,
				DirFileMode: 775,
				NoTempDir:   true,
			})
			if err != nil {
				return nil, err
			}
			buckets[k] = &Bucket{Bucket: bucket, exposedURL: opt.ExposedURL, name: k}
		default:

		}
	}
	return &Manager{
		buckets: buckets,
	}, nil
}

func (man *Manager) ExposedURL(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	if u.Scheme == "bucket" {
		buc, err := man.Get(u.Host)
		if err != nil {
			return rawUrl
		}
		return buc.ExposeURL(u.Path)
	}
	return rawUrl
}

func (man *Manager) Get(bucket string) (*Bucket, error) {
	b, ok := man.buckets[bucket]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}
