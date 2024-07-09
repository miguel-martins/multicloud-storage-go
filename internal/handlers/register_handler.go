package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/miguel-martins/multicloud-storage-go/internal/models"
	"github.com/miguel-martins/multicloud-storage-go/internal/repository"
	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

func RegisterHandler(userRepository *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		// Hash the password before saving
		hashedPassword, _ := util.HashPassword(user.Password)
		user.Password = hashedPassword

		// Save user to the database
		err := userRepository.Save(&user)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
