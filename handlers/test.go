package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Print the request body
	log.Println("Request Body:", string(body))

	// First unmarshal into a string
	var jsonString string
	err = json.Unmarshal(body, &jsonString)
	if err != nil {
		log.Printf("Error decoding JSON string: %v", err)
		http.Error(w, "Error decoding JSON string", http.StatusInternalServerError)
		return
	}

	// Then unmarshal the string into a map
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Printf("Error decoding JSON map: %v", err)
		http.Error(w, "Error decoding JSON map", http.StatusInternalServerError)
		return
	}
	log.Printf("Decoded JSON: %+v", data)

	// Access and print the value of key4 directly from the root level
	if key4Value, ok := data["key4"].(string); ok {
		log.Printf("Value of key4: %s\n", key4Value)
	} else {
		log.Println("key4 is not a string or not present")
	}

	// Send an OK response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
