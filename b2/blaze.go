package b2;

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Link struct {
  Text string
  Href string
}


func (l Link) String() string {
  return l.Text + ": " + l.Href
}

func GetBucketFiles(b2KeyID, b2Key, bucketName string) ([]Link, error) {
  bucket := aws.String(bucketName)

  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-central-003"),
    Endpoint: aws.String("https://s3.eu-central-003.backblazeb2.com"),
    Credentials: credentials.NewStaticCredentials(b2KeyID, b2Key, ""),
    S3ForcePathStyle: aws.Bool(true),
  });

  if err != nil {
    log.Printf("Failed to create session: %v", err)
    return nil, err
  }

  svc := s3.New(sess)
  response, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket})
  if err != nil {
    log.Printf("Unable to list items in bucket: %q, %v", *bucket, err)
    return nil, err
  }

  var links []Link

  for _, item := range response.Contents {
    objectKey := *item.Key 
    friendlyURL := fmt.Sprintf("https://f003.backblazeb2.com/file/%s/%s", *bucket, objectKey)
    links = append(links, Link{Text: objectKey, Href: friendlyURL})
  }

  return links, nil
}

