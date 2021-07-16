package email

import (
	"ImaginatoGolangTestTask/shared/log"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"net/smtp"
	"os"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.

	//Sender = "sender@example.com"
	Sender = "sonu.jain.mi@gmail.com"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "recipient@example.com"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Amazon SES Test (AWS SDK for Go)"

	// The HTML body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Request struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func (r *Request) SendEmail(file string, credential map[string]string) error {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_CLIENT_ID"), os.Getenv("AWS_CLIENT_SECRET"), ""),
	})
	if err != nil {
		return err
	}

	// Create SES service client
	svc := ses.New(sess)

	//result, err := svc.ListIdentities(&ses.ListIdentitiesInput{IdentityType: aws.String("EmailAddress")})
	//
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	//for _, email := range result.Identities {
	//	var e = []*string{email}
	//
	//	verified, err := svc.GetIdentityVerificationAttributes(&ses.GetIdentityVerificationAttributesInput{Identities: e})
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	for _, va := range verified.VerificationAttributes {
	//		if *va.VerificationStatus == "Success" {
	//			fmt.Println(*email)
	//		}
	//	}
	//}

	// Parsing template
	if err := r.parseTemplate(file, credential); err != nil {
		fmt.Println("Error : parsing error")
	}

	sesEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(credential["email"])},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(r.Body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(r.Subject),
			},
		},
		Source: aws.String("subhransu.mohapatra.mi@gmail.com"),
		ReplyToAddresses: []*string{
			aws.String("subhransu.mohapatra.mi@gmail.com"),
		},
	}

	_, err = svc.SendEmail(sesEmailInput)
	if err != nil {
		return err
	}
	return nil
}

// NewMailRequest is
func NewMailRequest(to []string, subject string) *Request {
	return &Request{
		To:      to,
		Subject: subject,
	}
}

// parseTemplate is for parsing HTML template
func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.Body = buffer.String()
	return nil
}

func (r *Request) SendSmtpEmail(name string, templatePath string) {
	password := os.Getenv("Password")
	host := os.Getenv("Host")
	port := os.Getenv("EmailPort")

	smtpServer := smtpServer{host: host, port: port}

	// Authentication.
	auth := smtp.PlainAuth("", r.From, password, smtpServer.host)
	pwd, errDir := os.Getwd()
	if errDir != nil {
		log.GetLog().Info("ERROR(from repo) : ", errDir.Error())
	}
	t, err := template.ParseFiles(pwd + templatePath)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
	}
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: "+r.Subject+" \n%s\n\n", mimeHeaders)))
	bodyErr := t.Execute(&body, struct {
		Name    string
		Message string
		URL     string
	}{
		Name:    name,
		Message: r.Body,
		URL:     r.Body,
	})
	if bodyErr != nil {
		log.GetLog().Info("ERROR(from repo) : ", bodyErr.Error())
	}

	// Sending email
	smtpErr := smtp.SendMail(smtpServer.Address(), auth, r.From, r.To, body.Bytes())
	if smtpErr != nil {
		log.GetLog().Warn("ERROR(from repo) : ", smtpErr.Error())
	}
	log.GetLog().Info("Email sent!", "Email sent successfully!")
}
