package emailVerifier

import (
	"errors"
	"net"
	"net/textproto"
	"strings"
)

// FindEmail takes a first, last and company name, finds MX and mail server addresses,
// creates possibilities for email addresses and tries them against the server.
func FindEmail(firstName string, lastName string, companyName string) (string, error) {
	test, err := tryEmails(firstName, lastName, companyName)
	if err != nil {
		return "", err
	}
	return test[0], nil
}

// VerifyEmail takes an email and checks whether the related MX server for the host agrees that it exists.
func VerifyEmail(email string) error {
	host := strings.Split(email, "@")[1]

	res, err := net.LookupMX(host)
	if err != nil {
		return errors.New("Incorrect Host Address")
	}
	mxServer := strings.TrimRight(res[0].Host, ".")
	conn, err := textproto.Dial("tcp", mxServer+":25")
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := setupMX(conn, email); err != nil {
		return err
	}

	if err := checkResponse(conn, "rcpt to: <"+email+">", 250); err != nil {
		return errors.New("Recipient " + email + " invalid: " + err.Error())
	}

	return nil
}
