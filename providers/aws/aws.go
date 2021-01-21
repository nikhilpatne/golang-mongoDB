package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CheckConnection(access_key string, secret_access_key string) bool {

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(access_key, secret_access_key, ""),
	})
	svc := ec2.New(sess)
	_, err := svc.DescribeKeyPairs(nil)
	if err != nil {
		return false
	}

	return true
}

func LaunchInstance(access_key, secret_access_key, image, size, vm_name string) (bool, error) {

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(access_key, secret_access_key, ""),
	})

	// Create EC2 service client
	svc := ec2.New(sess)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(image),
		InstanceType: aws.String(size),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		return false, err
	}

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(vm_name),
			},
		},
	})
	if errtag != nil {
		return false, errtag
	}

	return true, nil
}
