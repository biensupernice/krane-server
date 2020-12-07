package controllers

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/uuid"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
)

// RequestLoginPhrase : request a preliminary login request for authentication with the krane server.
// This will return a request id and phrase. The phrase should be encrypted using the clients private auth.
// This route does not return a token. You must use /auth and provide the signed phrase.
func RequestLoginPhrase(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.Generate().String()
	phrase := []byte(fmt.Sprintf("Authenticating with krane %s", reqID))

	err := store.Client().Put(constants.AuthenticationCollectionName, reqID, phrase)
	if err != nil {
		logger.Error(err)

		err = store.Client().Remove(constants.AuthenticationCollectionName, reqID)
		if err != nil {
			logger.Error(err)
			response.HTTPBad(w, err)
			return
		}
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, struct {
		RequestID string `json:"request_id"`
		Phrase    string `json:"phrase"`
	}{
		RequestID: reqID,
		Phrase:    string(phrase),
	})
}
