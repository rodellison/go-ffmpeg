package clients

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"

	//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"
)

var (
	ECSIfaceClient ecsiface.ECSAPI
	ECSSvcClient   *ecs.ECS
)

func init() {

	//During testing, we'll override the endpoint to ensure testing against local DynamoDB Docker image
	cfg := aws.Config{
		//		Endpoint: aws.String("http://localhost:8000"),
		Region:     aws.String("us-east-1"),
		MaxRetries: aws.Int(3),
	}

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the DynamoDB service client, to be used for inserting DB entries
	ECSSvcClient = ecs.New(sess, &cfg)
	ECSIfaceClient = ECSSvcClient

	//Note: XRay is unnecessary - but using it to try out tracing services..
	xray.AWS(ECSSvcClient.Client)

}

// func InsertDBEvent converts Eventdata into appropriate DynamoDB table attributes, and puts the item into the DB.
func RunECSTask(ctx aws.Context, s3_video_url, thumbnail_file, frame_pos string) (err error) {

	var myNWC = &ecs.NetworkConfiguration{
		AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
			AssignPublicIp: aws.String("ENABLED"),
			//SecurityGroups: nil,
			Subnets: []*string{aws.String(os.Getenv("ECS_TASK_VPC_SUBNET_1")),
				aws.String(os.Getenv("ECS_TASK_VPC_SUBNET_2")),
			},
		},
	}
	var myContainerOverride = &ecs.ContainerOverride{
	//	Command: nil,
	//	Cpu:     nil,
		Environment: []*ecs.KeyValuePair{
			{
				Name:  aws.String("INPUT_VIDEO_FILE_URL"),
				Value: aws.String(s3_video_url),
			},
			{
				Name:  aws.String("OUTPUT_THUMBS_FILE_NAME"),
				Value: aws.String(thumbnail_file),
			},
			{
				Name:  aws.String("POSITION_TIME_DURATION"),
				Value: aws.String(frame_pos),
			},
			{
				Name:  aws.String("OUTPUT_S3_PATH"),
				Value: aws.String(os.Getenv("OUTPUT_S3_PATH")),
			},
			{
				Name:  aws.String("AWS_REGION"),
				Value: aws.String(os.Getenv("OUTPUT_S3_AWS_REGION")),
			},
		},
	//	EnvironmentFiles:     nil,
	//	Memory:               nil,
	//	MemoryReservation:    nil,
		Name:                 aws.String(os.Getenv("ECS_CONTAINER_OVERRIDE_NAME")),
	//	ResourceRequirements: nil,
	}

	var myTaskOverrides = &ecs.TaskOverride{
		ContainerOverrides: []*ecs.ContainerOverride{
			myContainerOverride,
		},
	//	Cpu:                           nil,
	//	ExecutionRoleArn:              nil,
	//	InferenceAcceleratorOverrides: nil,
	//	Memory:                        nil,
	//	TaskRoleArn:                   nil,
	}

	var myRunTaskInput = ecs.RunTaskInput{
		//CapacityProviderStrategy: nil,
		Cluster: aws.String(os.Getenv("ECS_CLUSTER_NAME")),
		Count:   aws.Int64(1),
		//EnableECSManagedTags:     nil,
		//Group:                    nil,
		LaunchType:           aws.String("FARGATE"),
		NetworkConfiguration: myNWC,
		Overrides:            myTaskOverrides,
		//	PlacementConstraints:     nil,
		//	PlacementStrategy:        nil,
		PlatformVersion:          aws.String("LATEST"),
		//	PropagateTags:            nil,
		//	ReferenceId:              nil,
		//	StartedBy:                nil,
		//	Tags:                     nil,
		TaskDefinition: aws.String(os.Getenv("ECS_TASK_DEFINITION")),
	}

	_, err = ECSIfaceClient.RunTaskWithContext(ctx, &myRunTaskInput)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}

}
