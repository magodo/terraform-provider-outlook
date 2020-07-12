package clients

import (
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

type Client struct {
	UserFeature
	MailFolders  *msgraph.UserMailFoldersCollectionRequestBuilder
	MessageRules *msgraph.MailFolderMessageRulesCollectionRequestBuilder
}

func NewClient(b msgraph.BaseRequestBuilder, feature UserFeature) *Client {
	b.SetURL(b.URL() + "/me")
	userClient := msgraph.UserRequestBuilder{BaseRequestBuilder: b}
	return &Client{
		UserFeature:  feature,
		MailFolders:  userClient.MailFolders(),
		MessageRules: userClient.MailFolders().ID("inbox").MessageRules(),
	}
}
