# File Upload Service

Welcome to the File Upload Service! This project provides a backend service for uploading and managing files with user authentication and file access controls.

## Features

- User authentication with JWT tokens
- Secure password storage using bcrypt
- File upload and download
- Access control based on user authentication

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.16+)
- [Git](https://git-scm.com/)

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/bartzalewski/file-upload-service.git
   cd file-upload-service
   ```

2. **Initialize Go modules**:

   ```sh
   go mod tidy
   ```

3. **Run the application**:

   ```sh
   go run main.go
   ```

The server will start on `http://localhost:8080`.

### Directory Setup

The service will automatically create an `uploads` directory for storing uploaded files.

## API Endpoints

### User Authentication

#### Sign Up

- **Endpoint**: `POST /signup`
- **Request Body**:
  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```
- **Response**:
  ```json
  {
    "status": "User created successfully"
  }
  ```

#### Sign In

- **Endpoint**: `POST /signin`
- **Request Body**:
  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```
- **Response**:
  ```json
  {
    "status": "Signed in successfully"
  }
  ```
  A `token` cookie is set upon successful sign-in.

### File Management

#### Upload File

- **Endpoint**: `POST /upload`
- **Request**: Form-data with file field named `file`
- **Response**:
  ```json
  {
    "filename": "uploaded-file-name.txt",
    "uploaded_at": "2023-05-24T12:34:56Z",
    "uploader": "your-username"
  }
  ```

#### Download File

- **Endpoint**: `GET /files/{filename}`
- **Response**: The requested file is served for download.

## Testing the Service

Use tools like `curl` or Postman to test the endpoints.

### Example Requests

1. **Sign Up**:

   ```sh
   curl -X POST -d '{"username":"testuser", "password":"password"}' -H "Content-Type: application/json" http://localhost:8080/signup
   ```

2. **Sign In**:

   ```sh
   curl -X POST -d '{"username":"testuser", "password":"password"}' -H "Content-Type: application/json" http://localhost:8080/signin -c cookies.txt
   ```

3. **Upload File**:

   ```sh
   curl -X POST -F 'file=@path/to/your/file.txt' http://localhost:8080/upload -b cookies.txt
   ```

4. **Download File**:

   ```sh
   curl -O http://localhost:8080/files/file.txt -b cookies.txt
   ```

## Built With

- [Go](https://golang.org/) - The Go programming language
- [Gorilla Mux](https://github.com/gorilla/mux) - A powerful URL router and dispatcher for Golang
- [JWT-Go](https://github.com/dgrijalva/jwt-go) - A Go implementation of JSON Web Tokens
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - A package for password hashing

## Contributing

Feel free to submit issues or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go community for their invaluable resources and support.
- [Gorilla Mux](https://github.com/gorilla/mux) and [JWT-Go](https://github.com/dgrijalva/jwt-go) for making development easier.

---

Happy coding! 🚀
