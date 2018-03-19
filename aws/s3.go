package aws

import (
	_ "github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Bucket    string
	Key       string
	AccessKey string
	SecretKey string
}

func New(b string, k string) S3 {
	return S3{
		Bucket: b,
		Key:    k,
	}
}
func (s *S3) GetS3BucketContent() {}
