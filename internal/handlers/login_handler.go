package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/miguel-martins/multicloud-storage-go/internal/repository"
	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

func LoginHandler(userRepository *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Authenticate user
		authenticated, err := userRepository.Authenticate(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
			return
		}
		if !authenticated {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := util.GenerateToken(credentials.Username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Return JWT token
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
