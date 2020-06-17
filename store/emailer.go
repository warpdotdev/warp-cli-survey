package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Emailer sends summary emails
type Emailer struct {
	serverRoot string
}

// NewEmailer makes an emailer pointing at the given serverRoot
func NewEmailer(serverRoot string) *Emailer {
	return &Emailer{serverRoot: serverRoot}
}

// SendSummaryEmail sends the given summary to the given email address
func (e *Emailer) SendSummaryEmail(email string, summary string) {
	fmt.Println("Sending summary email to", email)

	summaryMap := make(map[string]interface{})
	summaryMap["email"] = email
	summaryMap["summary"] = summary
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(summaryMap)

	resp, err := http.Post(e.serverRoot+"/summary", "application/json", b)
	if err != nil {
		log.Println("Unable to save response...are you online?", err)
	}

	defer resp.Body.Close()
}
