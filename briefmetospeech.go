package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"net/http"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
	"io"
	"io/ioutil"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("You must supply a file name and a bucket name")
		os.Exit(1)
	}

	// The name of the text file to convert to MP3
	fileName := os.Args[1]
	bucketName := os.Args[2]

	// Open text file and get it's contents as a string
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Got error opening file " + fileName)
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Convert bytes to string
	s := string(contents[:])

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	// Create Polly client
	svc := polly.New(sess)
	fmt.Println(s)
	// Output to MP3 using voice Joanna
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Mathieu")}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("Got error calling SynthesizeSpeech:")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Save as MP3
	names := strings.Split(fileName, ".")
	name := names[0]
	name = "test"
	mp3File := name + ".mp3"

	outFile, err := os.Create(mp3File)
	if err != nil {
		fmt.Println("Got error creating " + mp3File + ":")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		fmt.Println("Got error saving MP3:")
		fmt.Print(err.Error())
		os.Exit(1)
	}else{
		AddFileToS3(sess, mp3File, bucketName)
	}
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileDir string, bucketName string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}