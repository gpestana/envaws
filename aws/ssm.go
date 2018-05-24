package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
)

type SSM struct {
	Client *ssm.SSM
}

func NewSSM(r string) SSM {
	awsConf := &aws.Config{Region: aws.String(r)}
	sess := session.Must(session.NewSession(awsConf))
	cli := ssm.New(sess)

	return SSM{
		Client: cli,
	}
}

func (s *SSM) GetContent() error {
	in := ssm.DescribeParametersInput{}
	res, err := s.Client.DescribeParameters(&in)
	if err != nil {
		return err
	}
	// TODO: check for next results if response is paginated

	log.Println(res.String())
	return nil
}
