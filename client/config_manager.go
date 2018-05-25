package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gpestana/envaws/config"
	"io/ioutil"
	"log"
	"time"
)

type ConfigManager interface {
	GetConfigurations() ([]byte, error)
	StartPolling()
}

type managerConfig struct {
	changedFn  func()
	pollingInt int
	accessKey  string
	secretKey  string
}

// Response from managerConfig
type Response struct {
	Configs map[string]string
	Errors  []string
}

func NewResponse() *Response {
	c := map[string]string{}
	e := []string{}
	return &Response{Configs: c, Errors: e}
}

func (r *Response) parseSsm(params *ssm.GetParametersOutput) {
	for _, p := range params.Parameters {
		r.Configs[*p.Name] = *p.Value
	}
	for _, e := range params.InvalidParameters {
		err := fmt.Sprintf("'%v' is an Invalid Parameter", *e)
		r.Errors = append(r.Errors, err)
	}
}

type S3 struct {
	Env             map[string]string
	Interface       *s3.S3
	lastCheckedETag string
	bucket          string
	key             string
	config          managerConfig
}

func New(c config.C, f func()) (ConfigManager, error) {
	var cm ConfigManager
	var err error

	switch c.Service {
	case "s3":
		cm = NewS3Manager(c, f)
	case "ssm":
		cm = NewSSMManager(c, f)
	default:
		err = errors.New(fmt.Sprintf("Service %v is not valid.", c.Service))
	}
	return cm, err
}

func NewS3Manager(c config.C, f func()) S3 {
	awsConf := &aws.Config{
		Region: aws.String(c.Region),
	}
	sess := session.Must(session.NewSession(awsConf))
	intf := s3.New(sess)

	return S3{
		Interface:       intf,
		lastCheckedETag: "",
		bucket:          c.S3.Bucket,
		key:             c.S3.Key,
		config: managerConfig{
			changedFn:  f,
			accessKey:  c.AccessKey,
			secretKey:  c.SecretKey,
			pollingInt: c.PollingInterval,
		},
	}
}

func (c S3) getObjectS3() (*s3.GetObjectOutput, error) {
	objIn := &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &c.key,
	}
	obj, err := c.Interface.GetObject(objIn)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

func (c S3) GetConfigurations() ([]byte, error) {
	out, err := c.getObjectS3()
	if err != nil {
		return nil, err
	}
	confs, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return confs, err
	}
	return confs, err
}

func (c S3) StartPolling() {
	for _ = range time.Tick(time.Duration(time.Duration(c.config.pollingInt) * time.Second)) {
		go func() {
			o, err := c.getObjectS3()
			if err != nil {
				log.Println(err)
				return
			}
			etag := *o.ETag

			if c.lastCheckedETag == "" {
				c.lastCheckedETag = etag
				return
			}

			if etag != c.lastCheckedETag {
				c.lastCheckedETag = etag
				c.config.changedFn()
			}
		}()
	}
}

type SSM struct {
	Parameters            []string
	Interface             *ssm.SSM
	lastCheckedConfigHash string
	config                managerConfig
}

func NewSSMManager(c config.C, f func()) SSM {
	awsConf := &aws.Config{Region: aws.String(c.Region)}
	sess := session.Must(session.NewSession(awsConf))
	intf := ssm.New(sess)

	return SSM{
		Interface:             intf,
		Parameters:            c.Ssm.Parameters,
		lastCheckedConfigHash: "",
		config: managerConfig{
			changedFn:  f,
			accessKey:  c.AccessKey,
			secretKey:  c.SecretKey,
			pollingInt: c.PollingInterval,
		},
	}
}

func (c SSM) GetConfigurations() ([]byte, error) {
	var confs []byte

	// builds []*string of parameter names as expected by ssm.GetParametersInput
	params := make([]*string, len(c.Parameters))
	for i, p := range c.Parameters {
		ptr := p
		params[i] = &ptr
	}

	in := ssm.GetParametersInput{
		Names: params,
	}
	out, err := c.Interface.GetParameters(&in)
	if err != nil {
		return confs, err
	}

	// parses configurationsOut into res structure expected by client
	resRaw := NewResponse()
	resRaw.parseSsm(out)
	confs, err = json.Marshal(resRaw)
	if err != nil {
		return confs, err
	}

	return confs, nil
}

func (c SSM) StartPolling() {
	for _ = range time.Tick(time.Duration(time.Duration(c.config.pollingInt) * time.Second)) {
		go func() {
			obj, err := c.GetConfigurations()
			if err != nil {
				log.Println(err)
				return
			}

			hash := md5Hash(obj)
			if c.lastCheckedConfigHash == "" {
				c.lastCheckedConfigHash = hash
				return
			}

			if hash != c.lastCheckedConfigHash {
				c.lastCheckedConfigHash = hash
				c.config.changedFn()
			}
		}()
	}
}

func md5Hash(b []byte) string {
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
