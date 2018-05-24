package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Region    string
	Bucket    string
	Key       string
	AccessKey string
	SecretKey string
	Client    *s3.S3
}

func NewS3(b string, k string, r string) S3 {
	awsConf := &aws.Config{
		Region: aws.String(r),
	}
	sess := session.Must(session.NewSession(awsConf))
	cli := s3.New(sess)

	return S3{
		Bucket: b,
		Key:    k,
		Client: cli,
	}
}

func (s *S3) getObjectS3() (*s3.GetObjectOutput, error) {
	objIn := &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &s.Key,
	}
	obj, err := s.Client.GetObject(objIn)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *S3) GetContent() (bytes.Buffer, error) {
	b, err := s.getObjectS3()
	if err != nil {
		return bytes.Buffer{}, err
	}

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(b.Body)

	return *bodyBuf, nil
}

func (s *S3) GeCurrentETag() (string, error) {
	b, err := s.getObjectS3()
	if err != nil {
		return "", err
	}
	return *b.ETag, nil
}
