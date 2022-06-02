package app

import (
	"net/http"
	"os"
	"strings"

	u "git.gonkar.com/gonkar/infra-cmd/utils"
)

var (
	mmToken, _ = os.LookupEnv("MM_TOKEN")
)

// JwtAuthentication - handeling of token creation
func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/v1/home/"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path            //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			//next.ServeHTTP(w, r)
			return
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		// Validate if received token is legit

		if tokenPart != mmToken {
			response = u.Message(false, "Invalid token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		// Proceed in the middleware chain!
		next.ServeHTTP(w, r)
	})
}
