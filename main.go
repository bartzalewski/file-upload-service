package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FileInfo struct {
	Filename   string    `json:"filename"`
	UploadedAt time.Time `json:"uploaded_at"`
	Uploader   string    `json:"uploader"`
}

type Store struct {
	sync.RWMutex
	users map[string]User
	files map[string]FileInfo
}

var store = Store{
	users: make(map[string]User),
	files: make(map[string]FileInfo),
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	store.Lock()
	store.users[creds.Username] = User{Username: creds.Username, Password: string(hashedPassword)}
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	store.RLock()
	user, exists := store.users[creds.Username]
	store.RUnlock()

	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	filePath := filepath.Join("uploads", filename)

	err = ioutil.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	fileInfo := FileInfo{
		Filename:   filename,
		UploadedAt: time.Now(),
		Uploader:   username,
	}

	store.Lock()
	store.files[filename] = fileInfo
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileInfo)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	params := mux.Vars(r)
	filename := params["filename"]

	store.RLock()
	fileInfo, exists := store.files[filename]
	store.RUnlock()

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Log the file access details
	log.Printf("File %s accessed by %s at %s", fileInfo.Filename, username, time.Now().Format(time.RFC3339))

	filePath := filepath.Join("uploads", filename)
	http.ServeFile(w, r, filePath)
}

func getUsernameFromRequest(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}

	tokenStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}
	if !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	return claims.Username
}

func main() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.HandleFunc("/upload", UploadFile).Methods("POST")
	r.HandleFunc("/files/{filename}", DownloadFile).Methods("GET")

	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
