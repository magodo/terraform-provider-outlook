package clients

import (
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

type Client struct {
	MailFolders  *msgraph.UserMailFoldersCollectionRequestBuilder
	MessageRules *msgraph.MailFolderMessageRulesCollectionRequestBuilder
	Categories	*msgraph.OutlookUserMasterCategoriesCollectionRequestBuilder
}

func NewClient(b msgraph.BaseRequestBuilder) *Client {
	b.SetURL(b.URL() + "/me")
	userClient := msgraph.UserRequestBuilder{BaseRequestBuilder: b}
	outlookClient := msgraph.OutlookUserRequestBuilder{BaseRequestBuilder: b}
	return &Client{
		MailFolders:  userClient.MailFolders(),
		MessageRules: userClient.MailFolders().ID("inbox").MessageRules(),
		Categories:   outlookClient.MasterCategories(),
	}
}
