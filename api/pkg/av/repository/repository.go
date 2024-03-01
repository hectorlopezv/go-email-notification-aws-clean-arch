package av

import (
	"av-send-email/api/pkg/entities"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Repository interface {
    ExecuteClamAvScan(paths []string) (*entities.ScanResult, error)
    ProcessScanResult(result *entities.ScanResult) (*entities.ProcessScanResult, error)
}

type repository struct {
    s3Client  *s3.Client
    sesClient *sesv2.Client
}

func (r *repository) ExecuteClamAvScan(paths []string) (*entities.ScanResult, error) {
    resultChan := make(chan *entities.ScanResult)
    go func() {
        stdout, stderr, err := r.runCommand(paths)
        result := &entities.ScanResult{}
        if err != nil {
            result.Error = err
            result.Stderr = stderr
        } else {
            result.Stdout = stdout
        }
        resultChan <- result
    }()

    result := <-resultChan
    if(result.Error != nil){
        return nil, result.Error
    }
    return result, nil
}

func (r *repository) ProcessScanResult(result *entities.ScanResult) (*entities.ProcessScanResult,error) {


    	 // Check for malware in the scan result
	 malwareFound := strings.Contains(result.Stdout, "Infected files:")
	 infectedFilesCount := 0
	 if malwareFound {
		 re := regexp.MustCompile(`Infected files:\s+(\d+)`)
		 matches := re.FindStringSubmatch(result.Stdout)
		 if len(matches) > 1 {
			 fmt.Sscanf(matches[1], "%d", &infectedFilesCount)
		 }
	 }
 
	 // Construct email response and S3 key
	 formatDate := time.Now().Format("2006-01-02")
	 var emailSubject, emailResponse, s3Key string
     enviroment := os.Getenv("ENVIROMENT")
	 if malwareFound && infectedFilesCount > 0 {
		 emailSubject = fmt.Sprintf("Malware Found on %s - %s", formatDate, enviroment)
		 emailResponse = "Malware found in the scan:\n\n" + result.Stdout
		 s3Key = fmt.Sprintf("malware_scan_results/%s/%s/malware_scan/%s.txt", enviroment,formatDate, strings.ReplaceAll(formatDate, "/", "_"))
	 } else {
		 emailSubject = fmt.Sprintf("No Malware Found on %s - %s", formatDate, enviroment)
		 emailResponse = "No malware found in the scan.\n\n" + result.Stdout
		 s3Key = fmt.Sprintf("malware_scan_results/%s/%s/no_malware_scan/%s.txt", enviroment,formatDate, strings.ReplaceAll(formatDate, "/", "_"))
	 }
        fmt.Printf("Result: %+v\n", emailResponse)
        fmt.Printf("Result: %+v\n", emailSubject)
        fmt.Printf("Result: %+v\n", s3Key)
    sesClient := r.sesClient
    s3Client := r.s3Client
    sesParams := &sesv2.SendEmailInput{
        FromEmailAddress: aws.String("securityalerts@sovereignrealestategroup.com"),
        Destination: &types.Destination{
            ToAddresses: []string{
                "alk@3atlantic.com",
                "hlopez@sovereignrealestategroup.com",
                "kbollepogu@sovereignrealestategroup.com",
            },
        },
        Content: &types.EmailContent{
            Simple: &types.Message{
                Subject: &types.Content{
                    Data: aws.String(emailSubject),
                },
                Body: &types.Body{
                    Text: &types.Content{
                        Data: aws.String(emailResponse),
                    },
                },
            },
        },
    }    
     _, errSendEmail := sesClient.SendEmail(context.TODO(), sesParams)
     if errSendEmail != nil {
         fmt.Println("Error sending email:", errSendEmail)
         return nil, errSendEmail
     }

         // Upload scan result to S3 bucket
    s3Params := &s3.PutObjectInput{
        Bucket: aws.String("avlogssvg"),
        Key:    aws.String(s3Key),
        Body:   strings.NewReader(emailResponse),
    }

    _, errS3 := s3Client.PutObject(context.TODO(), s3Params)
    if errS3 != nil {
        fmt.Println("Error uploading to S3:", errS3)
        return nil, errS3
    }

    // Send appropriate response based on scan results
    if malwareFound && infectedFilesCount > 0 {
        fmt.Println("Malware found. Email sent to recipients.")
        return &entities.ProcessScanResult{
            MalwareFound:       true,
            InfectedFilesCount: infectedFilesCount,
        },nil
    } else {
        fmt.Println("Scan completed. Scan result uploaded to S3.")
        return &entities.ProcessScanResult{
            MalwareFound:       false,
            InfectedFilesCount: 0,
        },nil
    }
}

func (r *repository) runCommand(paths []string) (stdout string, stderr string, err error) {
    var outb, errb bytes.Buffer

        // Starting with the base command and arguments
        args := []string{"-ri"}
        args = append(args, paths...)
        cmd := exec.Command("clamscan", args...)
    cmd.Stdout = &outb
    cmd.Stderr = &errb

    err = cmd.Run()
    return outb.String(), errb.String(), err
}

func NewRepository(s3Client *s3.Client, sesClient *sesv2.Client) Repository {
    return &repository{
        s3Client,
        sesClient,
    }
}
