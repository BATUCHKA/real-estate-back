package util

import (
	"fmt"
	"net/http"
	"os"
)

func SendOTPToMobile(otp string, phoneNumber string) error {

	smsAPIURL := os.Getenv("SMS_API_URL")

	if smsAPIURL == "" {
		return fmt.Errorf("SMS_API_URL not set in environment variables")
	}

	apiKey := os.Getenv("SMS_API_KEY")

	if apiKey == "" {
		return fmt.Errorf("SMS_API_KEY not set in environment variables")
	}

	smsNumber := os.Getenv("SMS_NUMBER")

	if smsNumber == "" {
		return fmt.Errorf("SMS_NUMBER not set in environment variables")
	}

	req, err := http.NewRequest("GET", smsAPIURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)

	q := req.URL.Query()

	q.Add("from", smsNumber)
	q.Add("to", phoneNumber)
	q.Add("text", fmt.Sprintf("Batalgaajuulah code: %s", otp))

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)

	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMS API error: %s", resp.Status)
	}

	fmt.Printf("OTP sent to mobile: %s\n", phoneNumber)

	return nil
}

func SendLinkToMobile(text string, phoneNumber string) error {
	smsAPIURL := os.Getenv("SMS_API_URL")

	if smsAPIURL == "" {
		return fmt.Errorf("SMS_API_URL not set in environment variables")
	}

	apiKey := os.Getenv("SMS_API_KEY")

	if apiKey == "" {
		return fmt.Errorf("SMS_API_KEY not set in environment variables")
	}

	smsNumber := os.Getenv("SMS_NUMBER")

	if smsNumber == "" {
		return fmt.Errorf("SMS_NUMBER not set in environment variables")
	}

	req, err := http.NewRequest("GET", smsAPIURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)

	q := req.URL.Query()
	q.Add("from", smsNumber)
	q.Add("to", phoneNumber)
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMS API error: %s", resp.Status)
	}

	fmt.Printf("SMS sent to mobile: %s\n", phoneNumber)

	return nil
}
