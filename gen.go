package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/minio/cli"
)

// Generate lambda.
var genCmd = cli.Command{
	Name:   "gen",
	Usage:  "Generates lambda",
	Action: mainGen,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bucket",
			Usage: "Bucket to run lambda on.",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "Prefix to run lambda on.",
		},
		cli.StringFlag{
			Name:  "suffix",
			Usage: "Suffix to run lambda on.",
		},
		cli.StringFlag{
			Name:  "events",
			Usage: "Events to run lambda on.",
			Value: "s3:ObjectCreated:*,s3:ObjectRemoved:*",
		},
	},
	CustomHelpTemplate: `NAME:
   minl {{.Name}} - {{.Usage}}

USAGE:
   minl {{.Name}} [FLAGS]

FLAGS:
  {{range .Flags}}{{.}}
  {{end}}
`,
}

// LambdaMetadata struct.
type LambdaMetadata struct {
	PackageName string
	Endpoint    string
	AccessKey   string
	SecretKey   string
	Secure      bool
	Region      string
	Bucket      string
	Events      []string
	Prefix      string
	Suffix      string
}

var notificationFile = `
package main

import (
        "fmt"
        "github.com/minio/minio-go/v6"
)

func enableBucketNotification(s3Client *minio.Client, lambdaArn minio.Arn) error {
       // ARN represents a notification channel that needs to be created in your S3 provider
       //  (e.g. http://docs.aws.amazon.com/sns/latest/dg/CreateTopic.html)

       // An example of an lambda ARN:
       //             arn:minio:lambda:us-east-1:myfunc:lambda
       //                  ^      ^      ^         ^       ^
       //       Provider __|      |      |         |       |
       //                         |    Region  Account ID  |_ Notification Name
       //                Service _|
       //
       // You should replace YOUR-PROVIDER, YOUR-SERVICE, YOUR-REGION, YOUR-ACCOUNT-ID and YOUR-RESOURCE
       // with actual values that you receive from the S3 provider

       // Here you create a new lambda notification
       lambdaConfig := minio.NewNotificationConfig(lambdaArn)
       {{range $event := .Events}}
       lambdaConfig.AddEvents(minio.NotificationEventType("{{ $event }}"))
       {{end}}
       {{if .Prefix}} lambdaConfig.AddFilterPrefix("{{ .Prefix }}") {{end}}
       {{if .Suffix}} lambdaConfig.AddFilterSuffix("{{ .Suffix }}") {{end}}

       // Now, set all previously created notification configs
       bucketNotification := minio.BucketNotification{}
       bucketNotification.AddLambda(lambdaConfig)
       err := s3Client.SetBucketNotification("{{ .Bucket }}", bucketNotification)
       if err != nil {
            return err
       }
       return nil
}

func listenBucketNotification(s3Client *minio.Client, bucket string, lambdaArn minio.Arn, lambdaFn LambdaFunc) error {
       // Create a done channel to control 'ListenBucketNotification' go routine.
       doneCh := make(chan struct{})

       // Indicate a background go-routine to exit cleanly upon return.
       defer close(doneCh)

       // Notification account ARN, variable fields to note here are:
       //  - region
       //  - account-id
       // Account id here is the same which was chosen during SetBucketNotification.
       for notificationInfo := range s3Client.ListenBucketNotification(bucket, lambdaArn, doneCh) {
               if err := lambdaFn(notificationInfo.Records, notificationInfo.Err); err != nil {
                       return err
               }
       }
       return nil
}

func main() {
       // Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
       // This boolean value is the last argument for New().

       // New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
       // determined based on the Endpoint value.
       s3Client, err := minio.New("{{ .Endpoint }}", "{{ .AccessKey }}", "{{ .SecretKey }}", {{ .Secure }})
       if err != nil {
             fmt.Println("Unable to initialize minio client", err)
             return
       }

       // Create a new lambda ARN.
       lambdaArn := minio.NewArn("minio", "lambda", "{{ .Region }}", "{{ .PackageName }}", "lambda")

       // Enable bucket notifications.
       if err = enableBucketNotification(s3Client, lambdaArn); err != nil {
              fmt.Println("Unable to enable bucket notification.", err)
              return
       }     

       // Listen bucket notification.
       if err = listenBucketNotification(s3Client, "{{ .Bucket }}", lambdaArn, YourFunc); err != nil {
              fmt.Println("Unable to listen bucket notification.", err)
              return
       }
}

type LambdaFunc func(events []minio.NotificationEvent, err error) error

func YourFunc(events []minio.NotificationEvent, err error) error {
       if err != nil {
               return err
       }
       /// Your code here.
       return nil
}
`

var supportedEventTypes = []string{
	"s3:ObjectCreated:*",
	"s3:ObjectCreated:Put",
	"s3:ObjectCreated:Post",
	"s3:ObjectCreated:Copy",
	"sh:ObjectCreated:CompleteMultipartUpload",
	"s3:ObjectRemoved:*",
	"s3:ObjectRemoved:Delete",
}

func isValidEventType(eventType string) bool {
	for _, supportedEventType := range supportedEventTypes {
		if supportedEventType == eventType {
			return true
		}
	}
	return false
}

func parseEvents(events []string) (parsedEvents []string) {
	for _, eventType := range events {
		if isValidEventType(eventType) {
			parsedEvents = append(parsedEvents, eventType)
		}
	}
	return parsedEvents
}

// checkGenSyntax - validate all the passed arguments
func checkGenSyntax(ctx *cli.Context) {
	if !ctx.Args().Present() {
		cli.ShowCommandHelpAndExit(ctx, ctx.Args().First(), 1)
	}
}

func initLambdaDir(lambda string) error {
	return os.MkdirAll(lambda, 0755)
}

func newLambdaMeta(ctx *cli.Context) LambdaMetadata {
	lmeta := LambdaMetadata{
		PackageName: ctx.Args().First(),
		Endpoint:    os.Getenv("S3_ENDPOINT"),
		AccessKey:   os.Getenv("ACCESS_KEY"),
		SecretKey:   os.Getenv("SECRET_KEY"),
		Secure:      os.Getenv("S3_SECURE") == "1",
		Region:      os.Getenv("S3_REGION"),
		Bucket:      ctx.String("bucket"),
		Events:      parseEvents(strings.Split(ctx.String("events"), ",")),
		Prefix:      ctx.String("prefix"),
		Suffix:      ctx.String("suffix"),
	}
	return lmeta
}

func mainGen(ctx *cli.Context) {
	checkGenSyntax(ctx)

	lambda := ctx.Args().First()
	tmpl := template.Must(template.New(lambda).Parse(notificationFile))

	initLambdaDir(lambda)
	lmeta := newLambdaMeta(ctx)

	templateFile := path.Join(lambda, lambda+".go")
	w, err := os.Create(templateFile)
	if err != nil {
		fmt.Println("Unable to write", templateFile, err)
		return
	}
	defer w.Close()
	err = tmpl.Execute(w, lmeta)
	if err != nil {
		fmt.Println("Unable to write", templateFile, err)
		return
	}
}
