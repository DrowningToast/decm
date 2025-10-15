package s3

import "fmt"

type S3Config struct {
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	BucketName      string `env:"BUCKET_NAME"`
	Endpoint        string `env:"ENDPOINT"`
}

func (c *S3Config) IsValid() bool {
	// Region is not required - will use default if not set
	return c.AccessKeyID != "" && c.SecretAccessKey != "" && c.BucketName != "" && c.Endpoint != ""
}

func (c *S3Config) String() string {
	return fmt.Sprintf("s3://%s@%s/%s", c.AccessKeyID, c.Endpoint, c.BucketName)
}
