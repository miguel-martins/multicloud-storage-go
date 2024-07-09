package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/miguel-martins/multicloud-storage-go/internal/handlers"
	"github.com/miguel-martins/multicloud-storage-go/internal/middleware"
	"github.com/miguel-martins/multicloud-storage-go/internal/repository"
	"github.com/miguel-martins/multicloud-storage-go/internal/util"
)

// Application encapsulates the configuration and dependencies of the application.
type Application struct {
	Client *http.Client
	DB     *sql.DB
}

type GlobalPublicKeyResponse struct {
	GlobalPublicKey string `json:"global_public_key"`
}

// NewApplication creates a new instance of the Application.
func NewApplication(clientCert, clientKey string) (*Application, error) {
	clientTLSConfig, err := loadClientTLSConfig(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientTLSConfig,
		},
	}

	db, err := util.InitDB()
	if err != nil {
		return nil, err
	}

	return &Application{
		Client: client,
		DB:     db,
	}, nil
}

// Start initializes routes and starts the server.
func (app *Application) Start() error {
	// Initialize UserRepository
	userRepository := repository.NewUserRepository(app.DB)

	// Define HTTP routes
	http.HandleFunc("/register", handlers.RegisterHandler(userRepository))
	http.HandleFunc("/login", handlers.LoginHandler(userRepository))
	http.Handle("/upload", middleware.JWTMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))

	// Start server
	log.Println("Starting Multicloud Storage API service...")
	log.Println(`
	    __  ___      ____  _      __                __   _____ __
	   /  |/  /_  __/ / /_(_)____/ /___  __  ______/ /  / ___// /_____  _________  ____  ___
	  / /|_/ / / / / / __/ / ___/ / __ \/ / / / __  /   \__ \/ __/ __ \/ ___/ __ \/ __ \/ _ \
	 / /  / / /_/ / / /_/ / /__/ / /_/ / /_/ / /_/ /   ___/ / /_/ /_/ / /  / /_/ / /_/ /  __/
	/_/  /_/\__,_/_/\__/_/\___/_/\____/\__,_/\__,_/   /____/\__/\____/_/   \__,_/\__, /\___/
	                                                                            /____/		`)
	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", nil)
}

// loadClientTLSConfig loads the client TLS configuration.
func loadClientTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		ServerName:   "localhost",
		RootCAs:      caCertPool,
	}, nil
}

func (app *Application) getGlobalPublicKey(taURL string) (string, error) {
	resp, err := app.Client.Get(taURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("unexpected status code: " + resp.Status)
	}

	var globalPublicKeyResponse GlobalPublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&globalPublicKeyResponse); err != nil {
		return "", err
	}

	return globalPublicKeyResponse.GlobalPublicKey, nil
}

// storeGlobalPublicKey stores the global public key securely.
func storeGlobalPublicKey(globalPublicKey string) {
	// Store the Global Public Key securely (e.g., in a configuration file, database, or encrypted storage)
	log.Println("Received Global Public Key:", globalPublicKey)
}

func main() {
	app, err := NewApplication("client.crt", "client.key")
	if err != nil {
		log.Fatal("Error initializing application:", err)
	}

	globalPublicKey, err := app.getGlobalPublicKey("https://localhost:443/gpk")
	if err != nil {
		log.Fatal("Error getting Global Public Key from TA: ", err)
	}

	storeGlobalPublicKey(globalPublicKey)

	if err := app.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
