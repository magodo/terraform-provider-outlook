package clients

import (
	"net/http"

	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

type Client struct {
	UserFeature
	Batch        *msgraph.BatchRequestBuilder
	MailFolders  *msgraph.UserMailFoldersCollectionRequestBuilder
	MessageRules *msgraph.MailFolderMessageRulesCollectionRequestBuilder
}

func NewClient(client *http.Client, feature UserFeature) *Client {
	return &Client{
		UserFeature:  feature,
		Batch:        msgraph.NewBatch(client),
		MailFolders:  msgraph.NewClient(client).Me().MailFolders(),
		MessageRules: msgraph.NewClient(client).Me().MailFolders().ID("inbox").MessageRules(),
	}
}
