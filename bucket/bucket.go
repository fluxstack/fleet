package bucket

import "gocloud.dev/blob"

type Bucket struct {
	*blob.Bucket
}
