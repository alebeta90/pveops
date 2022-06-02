package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	status := data["status"].(bool)
	//message := data["message"].(string)
	if !status {
		http.Error(w, "", http.StatusBadRequest)
	}
	fmt.Printf("this is data: %v\n", data)
	json.NewEncoder(w).Encode(data)
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
