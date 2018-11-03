package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("You must supply a bucket name and a label id")
		os.Exit(1)
	}

	// The name of the text file to convert to MP3
	bucketName := os.Args[1]
	labelId := os.Args[2]

	emails := retrieveEmail(labelId)
	fmt.Print(emails)

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	// Create Polly client
	svc := polly.New(sess)

	for _, s := range emails {
		fmt.Println(len(s))

		input := &polly.StartSpeechSynthesisTaskInput{OutputS3BucketName: aws.String(bucketName), OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Mathieu")}

		_, err := svc.StartSpeechSynthesisTask(input)
		if err != nil {
			fmt.Println("Got error calling SynthesizeSpeech:")
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}
}