package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"gitlab.com/steppelink/odin/odin-backend/database"
	"gitlab.com/steppelink/odin/odin-backend/database/models"
)

type BankAccountData struct {
	IBAN  string `json:"iban"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
	MsgId string `json:"msgId"`
}

func CheckBankAccountNumber(bankID, bankAccountNumber string) (BankAccountData, error) {
	db := database.Database

	paymentAPIURL := os.Getenv("PAYMENT_API_URL")
	if paymentAPIURL == "" {
		return BankAccountData{}, errors.New("payment api url not found")
	}

	token := os.Getenv("PAYMENT_API_KEY")
	if token == "" {
		return BankAccountData{}, errors.New("payment api key not found")
	}

	var bankCheck models.BankCode
	if result := db.GormDB.Where("id = ?", bankID).First(&bankCheck); result.Error != nil {
		return BankAccountData{}, errors.New("bank code not found")
	}

	url := paymentAPIURL + "/account/check?acct=" + bankAccountNumber + "&bank_code=" + bankCheck.BankCode
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BankAccountData{}, err
	}
	req.Header.Set("Authorization", "Token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return BankAccountData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BankAccountData{}, errors.New("failed to check bank account")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BankAccountData{}, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return BankAccountData{}, errors.New("invalid response structure")
	}

	bankAccountData := BankAccountData{}
	if iban, ok := data["iban"].(string); ok {
		bankAccountData.IBAN = iban
	}
	if alias, ok := data["alias"].(string); ok {
		bankAccountData.Alias = alias
	}
	if name, ok := data["name"].(string); ok {
		bankAccountData.Name = name
	} else {
		return BankAccountData{}, errors.New("Хүлээн авагчийн дансны дугаар ДНС бүртгэлгүй байна ")
	}
	if msgId, ok := data["msgId"].(string); ok {
		bankAccountData.MsgId = msgId
	}

	return bankAccountData, nil
}

func NameMatches(actual, returned string) bool {
	actual = strings.ToLower(strings.TrimSpace(actual))
	actual = strings.Replace(actual, "-", "", -1)

	returned = strings.ToLower(strings.TrimSpace(returned))
	returned = strings.Replace(returned, "-", "", -1)

	actualWords := strings.Fields(actual)
	returnedWords := strings.Fields(returned)

	actualWordCount := make(map[string]int)
	returnedWordCount := make(map[string]int)

	for _, word := range actualWords {
		actualWordCount[word]++
	}
	for _, word := range returnedWords {
		returnedWordCount[word]++
	}

	if len(actualWordCount) != len(returnedWordCount) {
		return false
	}
	for word, count := range actualWordCount {
		if returnedWordCount[word] != count {
			return false
		}
	}

	return true
}
func TransferToBankAccount(account, bankCode string, amount float64, description, currency, transferID string) error {
	paymentAPIURL := os.Getenv("PAYMENT_API_URL")
	if paymentAPIURL == "" {
		return errors.New("payment api url not found")
	}

	token := os.Getenv("PAYMENT_API_KEY")
	if token == "" {
		return errors.New("payment api key not found")
	}

	reqBody := map[string]interface{}{
		"account":     account,
		"bank_code":   bankCode,
		"amount":      amount,
		"description": description,
		"currency":    currency,
		"transferid":  transferID,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", paymentAPIURL+"/account/transfer", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to transfer")
	}

	return nil
}
