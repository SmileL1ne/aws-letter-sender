package sender

import (
	"fmt"
)

type IAWSSender interface {
	SendLetter(email, subject, body string) (string, error)
	SendLetterTemplate(email, subject, tmplpath string, data interface{}) (string, error)
}

func NewSender(email, region, format string) (IAWSSender, error) {
	switch {
	case format == ".txt" || format == ".text" || format == ".html":
		return &Sender{
			email:  email,
			region: region,
			format: format[1:]}, nil
	default:
		return nil, fmt.Errorf("%w: '%v'", ErrUndefinedSender, format)
	}
}
