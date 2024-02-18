package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Harmew/brokerMail/models"
	"github.com/Harmew/brokerMail/utils"
)

// SendGridConfig struct to hold the SendGrid configuration and implement the ServeHTTP method
type SendGridConfig struct {
	Authorization string
	ContentType   string
	URL           string
	Sender        string
}

func main() {
	// Load .env file (Development)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Handle the /send route with the SendGridConfig struct
	mux.Handle("/send", SendGridConfig{
		// Get the SendGrid API key and URL from the environment variables
		Authorization: "Bearer " + os.Getenv("SENDGRID_API_KEY"),
		ContentType:   "application/json",
		URL:           os.Getenv("SENDGRID_API_URL"),
		Sender:        os.Getenv("SENDER"),
	})

	// Start the server
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}

// ServeHTTP method to handle the /send route
func (sendGridConfig SendGridConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Switch to handle the request method
	switch r.Method {
	case "POST":
		sendMail(sendGridConfig, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// sendMail function to send the mail using the SendGrid API
func sendMail(sendGridConfig SendGridConfig, w http.ResponseWriter, r *http.Request) {
	// Create a new context with a timeout of 30 seconds
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the request body into the SendGridInternal struct
	var interfaceBody models.SendGridInternal
	err = json.Unmarshal(body, &interfaceBody)
	if err != nil {
		http.Error(w, "Error unmarshalling the request body", http.StatusInternalServerError)
		return
	}

	// Validate the request body using the ValidateJSON function (internal package)
	err = utils.ValidateJSON(interfaceBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send the request to SendGrid using the requestSendGrid function
	response, err := requestSendGrid(ctx, sendGridConfig, interfaceBody)
	if err != nil {
		http.Error(w, "Error sending the request to SendGrid", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, "Error closing the response body", http.StatusInternalServerError)
			return
		}
	}(response.Body)

	// Write the response to the client
	w.WriteHeader(response.StatusCode)
	// Copy the response body to the client
	_, err = io.Copy(w, response.Body)
	if err != nil {
		http.Error(w, "Error writing the response", http.StatusInternalServerError)
		return
	}
}

// requestSendGrid function to send the request to SendGrid
func requestSendGrid(ctx context.Context, sendGridConfig SendGridConfig, interfaceBody models.SendGridInternal) (*http.Response, error) {
	// Create the SendGridRequest struct
	sendGridRequest := models.SendGridRequest{
		Personalizations: []models.Personalizations{},
		From: models.From{
			Email: sendGridConfig.Sender,
		},
		Subject: interfaceBody.Subject,
		Content: []models.Content{{Type: "text/plain", Value: interfaceBody.Content}},
	}

	sendGridRequest.Personalizations = append(sendGridRequest.Personalizations, models.Personalizations{
		To: []models.To{},
	})

	for _, recipient := range interfaceBody.Recipients {
		sendGridRequest.Personalizations[0].To = append(sendGridRequest.Personalizations[0].To, models.To{
			Email: recipient,
		})
	}

	// Marshal the SendGridRequest struct into JSON
	jsonData, err := json.Marshal(sendGridRequest)
	if err != nil {
		return nil, errors.New("error marshalling the request body")
	}

	// Create a new request with the context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendGridConfig.URL, nil)
	if err != nil {
		return nil, errors.New("error creating the request")
	}

	// Set the request headers
	req.Header.Set("Authorization", sendGridConfig.Authorization)
	req.Header.Set("Content-Type", sendGridConfig.ContentType)

	// Set the request body
	req.Body = io.NopCloser(bytes.NewReader(jsonData))

	// Send the request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("error sending the request")
	}

	// Return the response and nil error
	return response, nil
}
