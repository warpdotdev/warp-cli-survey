package store

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type webStore struct {
	serverRoot string
}

// NewWebStore makes a webstore pointing at the given serverRoot
func NewWebStore(serverRoot string) Storer {
	return &webStore{serverRoot: serverRoot}
}

func (ws *webStore) Write(response Response) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(response)

	resp, err := http.Post(ws.serverRoot+"/answer", "application/json", b)
	if err != nil {
		log.Println("Unable to save response. Mind connecting to the internet and trying again?")
		return
	}

	defer resp.Body.Close()
}
