package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileRepository struct{}

func CreateFileRepository() FileRepository {
	return FileRepository{}
}

func (FileRepository) OnTransactionsAvailable(transactions []Transaction) {
	f, err := os.OpenFile("transactions", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening transactions file")
	}

	for _, transaction := range transactions {
		jsonData, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println("Error marshalling transaction to json")
		}
		f.Write(jsonData)
		f.WriteString("\n")
	}

	f.Close()
}
