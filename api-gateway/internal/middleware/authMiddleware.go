package middleware

import (
	"log"
	"net/http"
	"os"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Println("Error getting cookie: ", err)
			http.Error(rw, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("Validating token: ", cookie.Value)

		req, err := http.NewRequest("GET", os.Getenv("AUTH_SERVICE_URL")+"/validate-token", nil)
		if err != nil {
			log.Println("Error creating validation request: ", err)
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", "Bearer "+cookie.Value)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("Invalid token or error validating: ", err)
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	}
}
