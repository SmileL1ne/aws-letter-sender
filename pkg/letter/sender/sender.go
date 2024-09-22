package sender

import (
	"bytes"
	"fmt"

	tmpl "text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	text = "txt"
	html = "html"
)

type Sender struct {
	email  string // ex: example@example.com
	region string // ex: eu-west-1
	format string // ex: txt, html
}

func (s *Sender) SendLetterTemplate(email, subject, tmplpath string, data interface{}) (string, error) {
	t, err := tmpl.ParseFiles(tmplpath)
	if err != nil {
		return "", fmt.Errorf("ParseFiles: %v", err)
	}

	var body bytes.Buffer
	err = t.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("(SendLetterTemplate) template Execute: %v", err)
	}

	msgID, err := s.SendLetter(email, subject, body.String())
	if err != nil {
		return "", fmt.Errorf("(SendLetterTemplate) SendLetter: %v", err)
	}
	return msgID, err
}

func (s *Sender) SendLetter(email, subject, body string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region)},
	)
	if err != nil {
		return "", fmt.Errorf("session.NewSession: %v", err)
	}

	var charset string = "UTF-8"

	// Create an SES session.
	svc := ses.New(sess)
	// Assemble the email.
	input := s.assembleEmail(email, subject, charset, body)
	// Attempt to send the email.
	sent, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return "", fmt.Errorf("svc.SendEmail: %s, error: %s", ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return "", fmt.Errorf("svc.SendEmail: %s; error: %s", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return "", fmt.Errorf("svc.SendEmail: %s; error: %s", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				return "", fmt.Errorf("svc.SendEmail error: %s" + aerr.Error())
			}
		} else {
			return "", fmt.Errorf("svc.SendEmail error: %s", err.Error())
		}
	}

	return *sent.MessageId, nil
}

func (s *Sender) assembleEmail(email, subject, charset, body string) *ses.SendEmailInput {
	letterBody := &ses.Body{}
	switch s.format {
	case text:
		letterBody = &ses.Body{
			Text: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(body),
			},
		}
	case html:
		letterBody = &ses.Body{
			Html: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(body),
			},
		}
	}

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: letterBody,
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.email),
	}
}
