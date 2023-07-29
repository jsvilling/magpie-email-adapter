package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type TransactionAvailableListener interface {
	OnTransactionsAvailable(transactions []Transaction)
}

type GmailSrv struct {
	gmail.Service
	transactionAvailableListeners []TransactionAvailableListener
}

func SetupGmailSrv(ctx *context.Context, listeners *[]TransactionAvailableListener) *GmailSrv {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		fmt.Print("Error reading file")
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		fmt.Print("Error creating config")
	}

	client := getClient(config)

	googleService, err := gmail.NewService(*ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("Error creating gmailSrv")
	}

	return &GmailSrv{*googleService, *listeners}
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (gmailSrv *GmailSrv) LoadTransactions() []Transaction {
	call := gmailSrv.Users.Messages.List("me")
	call.LabelIds("Label_2644345171270105163")
	msgRespnse, err := call.Do()
	if err != nil {
		fmt.Print("Error loading message list")
	}

	transactions := make([]Transaction, len(msgRespnse.Messages))

	for i, l := range msgRespnse.Messages {
		transaction := gmailSrv.GetTransactionFromEmail(l.Id)
		transactions[i] = transaction
		fmt.Println(transaction)
	}

	for _, listener := range gmailSrv.transactionAvailableListeners {
		listener.OnTransactionsAvailable(transactions)
	}

	return transactions
}

func (gmailSrv *GmailSrv) GetTransactionFromEmail(messageId string) Transaction {

	msg, err := gmailSrv.Users.Messages.Get(currentUser, messageId).Format("full").Do()
	if err != nil {
		fmt.Println("error loading message content: ", err)
	}

	data, _ := base64.URLEncoding.DecodeString(msg.Payload.Parts[0].Body.Data)
	msgContent := string(data)

	trimCurrency := mkSplitGetAtN(currencyMarker, 0)
	trimAccountName := mkSplitGetAtN("\"", 0)

	iAccount := strings.Index(msgContent, accountMarker)
	iBalance := strings.Index(msgContent, balanceMarker)
	iAmount := strings.Index(msgContent, amountMarker)
	iDate := strings.Index(msgContent, dateMarker)
	iValuta := strings.Index(msgContent, valutaMarker)

	account := trimAccountName(strings.Trim(msgContent[iAccount+len(accountMarker):iBalance], "\n"))
	balance := trimCurrency(strings.Trim(msgContent[iBalance+len(balanceMarker):iAmount], "\n"))
	amount := trimCurrency(strings.Trim(msgContent[iAmount+len(amountMarker):iDate], "\n"))
	date := strings.Trim(msgContent[iDate+len(dateMarker):iValuta], "\n")

	return Transaction{account, balance, amount, date}
}

func mkSplitGetAtN(delim string, n int) func(string) string {
	return func(str string) string {
		return strings.Split(str, delim)[n]
	}
}
