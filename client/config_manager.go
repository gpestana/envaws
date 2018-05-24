package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	awss "github.com/gpestana/envaws/aws"
	"github.com/gpestana/envaws/config"
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

type S3Config struct {
	Env             map[string]string
	Interface       *awss.S3
	lastCheckedETag string
	config          managerConfig
}

func New(c config.C, f func()) (ConfigManager, error) {
	var cm ConfigManager
	var err error

	switch c.Service {
	case "s3":
		cm = NewS3ConfigManager(c, f)
	case "ssm":
		cm = NewSSMConfigManager(c, f)
	default:
		err = errors.New(fmt.Sprintf("Service %v is not valid.", c.Service))
	}
	return cm, err
}

func NewS3ConfigManager(c config.C, f func()) S3Config {
	intf := awss.NewS3(c.S3.Bucket, c.S3.Key, c.Region)

	return S3Config{
		Interface:       &intf,
		lastCheckedETag: "",
		config: managerConfig{
			changedFn:  f,
			accessKey:  c.AccessKey,
			secretKey:  c.SecretKey,
			pollingInt: c.PollingInterval,
		},
	}
}

func (c S3Config) GetConfigurations() ([]byte, error) {
	var res []byte
	r, err := c.Interface.GetContent()
	if err != nil {
		return res, err
	}
	res, err = json.Marshal(r)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (c S3Config) StartPolling() {
	for _ = range time.Tick(time.Duration(time.Duration(c.config.pollingInt) * time.Second)) {
		go func() {
			etag, err := c.Interface.GeCurrentETag()
			if err != nil {
				log.Println(err)
				return
			}

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

type SSMConfig struct {
	Parameters []string
	Interface  *ssm.SSM
	config     managerConfig
}

func NewSSMConfigManager(c config.C, f func()) SSMConfig {
	awsConf := &aws.Config{Region: aws.String(c.Region)}
	sess := session.Must(session.NewSession(awsConf))
	intf := ssm.New(sess)

	return SSMConfig{
		Interface:  intf,
		Parameters: c.Ssm.Parameters,
		config: managerConfig{
			changedFn:  f,
			accessKey:  c.AccessKey,
			secretKey:  c.SecretKey,
			pollingInt: c.PollingInterval,
		},
	}
}

func (c SSMConfig) GetConfigurations() ([]byte, error) {
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

	confs, err = json.Marshal(out)
	if err != nil {
		return confs, err
	}

	// TODO: zip configs into { parameter: value }
	return confs, nil
}

func (c SSMConfig) StartPolling() {}
