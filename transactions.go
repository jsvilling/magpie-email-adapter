package main

const currentUser = "me"
const accountMarker = "Konto: "
const balanceMarker = "Neuer Saldo: "
const amountMarker = "Betrag: "
const dateMarker = "Buchungsdatum: "
const valutaMarker = "Valutadatum: "
const currencyMarker = "CHF"

type Transaction struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
	Amount  string `json:"amount"`
	Date    string `json:"date`
}
