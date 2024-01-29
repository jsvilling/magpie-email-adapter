package main

const currentUser = "me"
const accountMarker = "Konto: "
const balanceMarker = "Neuer Saldo: "
const amountMarker = "Betrag: "
const dateMarker = "Buchungsdatum: "
const valutaMarker = "Valutadatum: "
const currencyMarker = "CHF"
const sourceMarker = "Buchung:"

type Transaction struct {
	Account string `json:"Account"`
	Balance string `json:"Balance"`
	Amount  string `json:"Amount"`
	Date    string `json:"Date`
	Source  string `json:"Source`
}
