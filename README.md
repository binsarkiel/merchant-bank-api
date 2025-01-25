# Merchant-Bank API

This is a simple API for simulating transactions between merchants and customers. It uses JSON files as a database for demonstration purposes.

## Features

-   **Login:** Users (merchants or customers) can log in using their username and password.
-   **Payment:** Logged-in users can transfer money to other registered users using the recipient's username.
-   **Logout:** Users can log out, invalidating their token.
-   **Authentication:** Uses JWT (JSON Web Token) for authentication and authorization.
-   **Database Simulation:** Uses JSON files (`users.json`, `sessions.json`, `transactions.json`) to simulate data storage.

## Project Structure

| Directory/File          | Description                                                                                               |
| ----------------------- | --------------------------------------------------------------------------------------------------------- |
| `api/`                  | Package untuk API handlers                                                                               |
| `api/handlers.go`       | Handlers untuk masing-masing endpoint                                                                     |
| `data/`                 | Direktori untuk file JSON data                                                                            |
| `data/sessions.json`    | File JSON untuk data sesi                                                                                |
| `data/transactions.json` | File JSON untuk data transaksi                                                                           |
| `data/users.json`       | File JSON untuk data user                                                                                |
| `go.mod`                | File Go modules                                                                                           |
| `go.sum`                | Checksum Go modules                                                                                       |
| `main.go`               | File utama aplikasi                                                                                      |
| `models/`               | Package untuk model data                                                                                   |
| `models/models.go`       | Struktur data (User, Transaction, Session, request/response models)                                     |
| `repository/`           | Package untuk akses data                                                                                  |
| `repository/repository.go` | Fungsi-fungsi untuk berinteraksi dengan file JSON                                                        |
| `services/`             | Package untuk business logic                                                                              |
| `services/services.go`   | Logika untuk autentikasi, payment, validasi token, invalidasi token                                  |
| `utils/`                | Package untuk fungsi-fungsi utility                                                                       |
| `utils/utils.go`        | Fungsi-fungsi bantuan (error handling, token extraction)                                                 |
| `.env`                  | Variabel environment (JWT secret key)                                                                    |

## High-Level Architecture

**Components:**

-   **Client:** Any application that can send HTTP requests (e.g., Postman, browser, mobile app).
-   **API (Gin):** Handles incoming HTTP requests, routing, authentication (via middleware), and communication with the `services` layer.
-   **Services:** Contains the core business logic, such as user authentication, payment processing, and token validation.
-   **Repository:** Handles data access to the data store (JSON files in this case). It abstracts the data storage details from the `services` layer.
-   **Data Store:** JSON files used to simulate a database (`users.json`, `sessions.json`, `transactions.json`, `invalidated_tokens.json`).

**Integration:**

1.  The client sends an HTTP request to the API.
2.  The API (Gin) receives the request and routes it to the appropriate handler.
3.  The `AuthMiddleware` (in `api/handlers.go`) validates the JWT token from the `Authorization` header.
4.  The handler calls the relevant function in the `services` layer.
5.  The `services` layer uses the `repository` layer to interact with the data store (JSON files).
6.  The `repository` layer reads from or writes to the JSON files.
7.  The result is returned back up the chain to the client as an HTTP response.

## Technologies Used

-   **Go:** Programming language.
-   **Gin:** Web framework.
-   **JWT:** Authentication.

## Packages Used

-   `github.com/gin-gonic/gin`
-   `github.com/golang-jwt/jwt/v5`
-   `github.com/google/uuid`
-   `github.com/joho/godotenv`
-   `golang.org/x/crypto/bcrypt`
-   `encoding/json`
-   `os`

## Getting Started

### Prerequisites

-   **Go:** Ensure Go is installed on your system (version 1.20 or higher is recommended).
-   **Git:** To clone the repository.

### Steps

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd merchant-bank-api
    ```

2.  **Install dependencies:**

    ```bash
    go mod download
    ```

3.  **Create the `.env` file:**

    -   Create a file named `.env` in the root directory of the project.
    -   Add the `JWT_SECRET_KEY` to the `.env` file:

        ```
        JWT_SECRET_KEY=YOUR_STRONG_SECRET_KEY
        ```

        **Important:** Replace `YOUR_STRONG_SECRET_KEY` with a strong, secure secret key. **Do not use this example key in production.**

4.  **Create JSON data files and add dummy data:**

    -   Create a `data` directory if it doesn't exist: `mkdir data`
    -   Create the necessary JSON files: `touch data/users.json data/sessions.json data/transactions.json data/invalidated_tokens.json`
    -   Initialize `users.json`, `transactions.json`, and `sessions.json` with an empty array `[]` and `invalidated_tokens.json` with empty object `{}`
    -   Add the following dummy data to `data/users.json` (passwords are hashed using bcrypt - you can generate your own using an online bcrypt generator or the provided `repository/repository_test.go` file):

        ```json
        [
            {
                "id": "user-id-1",
                "name": "John Doe",
                "username": "johndoe",
                "password": "$2a$10$S3Ej/D92x.gWj/pWnGyKDu8XguMGbUDाबनातीचाआता.wVp9O44a",
                "account_type": "customer",
                "account_balance": 1000
            },
            {
                "id": "user-id-2",
                "name": "Jane Smith",
                "username": "janesmith",
                "password": "$2a$10$vc6y.96y27a/L.h94r4v/uOk8x9l.C0e/h1vj/t.wzN/Hro.4k5mG",
                "account_type": "merchant",
                "account_balance": 5000
            }
        ]
        ```

5.  **Run the application:**

    ```bash
    go run main.go
    ```

    The server will start on port `8080`.

## API Endpoints

### 1. `/login`

-   **Method:** `POST`
-   **Request Body:**

    ```json
    {
        "username": "johndoe",
        "password": "password123"
    }
    ```

    | Field      | Type   | Description                                   |
    | :--------- | :----- | :-------------------------------------------- |
    | `username` | string | Username of the user (required)               |
    | `password` | string | Password of the user (required)               |

-   **Response (Success - 200 OK):**

    ```json
    {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

    | Field        | Type   | Description                     |
    | :----------- | :----- | :------------------------------ |
    | `access_token` | string | JWT access token for authentication |

-   **Response (Error - 401 Unauthorized):**

    ```json
    {
        "error": "invalid username or password"
    }
    ```

-   **Response (Error - 400 Bad Request):**

    ```json
    {
        "error": "..." // Error message related to invalid request body
    }
    ```

    -   **Response (Error - 500 Internal Server Error):**
        ```json
        {
            "error": "..." // Error message related to server error
        }
        ```

### 2. `/payment`

-   **Method:** `POST`
-   **Headers:**
    -   `Authorization`: `Bearer <TOKEN>` (use the access token received from `/login`)
-   **Request Body:**

    ```json
    {
        "recipient": "janesmith",
        "amount": 100
    }
    ```

    | Field       | Type    | Description                                       |
    | :---------- | :------ | :------------------------------------------------ |
    | `recipient` | string  | Username of the recipient (required)              |
    | `amount`    | number  | Amount of money to transfer (required)            |

-   **Response (Success - 200 OK):**

    ```json
    {
        "transaction": {
            "activity": "transfer_money",
            "transaction_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
            "sender": "John Doe",
            "recipient": "Jane Smith",
            "amount": 100,
            "created_at": "2023-10-27T14:00:00Z"
        }
    }
    ```

-   **Response (Error - 401 Unauthorized):**

    ```json
    {
        "error": "invalid token"
    }
    ```

    ```json
    {
        "error": "Authorization header required"
    }
    ```

    ```json
    {
        "error": "invalid token format. Use 'Bearer <token>'"`
    }
    ```

-   **Response (Error - 422 Unprocessable Entity):**

    ```json
    {
        "error": "insufficient balance" // Example: Not enough money in the account
    }
    ```

    ```json
    {
        "error": "recipient not found" // Example: Recipient user doesn't exist
    }
    ```

-   **Response (Error - 500 Internal Server Error):**
    ```json
    {
        "error": "failed to update sender's balance"
    }
    ```
    ```json
    {
        "error": "failed to update recipient's balance"
    }
    ```

### 3. `/logout`

-   **Method:** `DELETE`
-   **Headers:**
    -   `Authorization`: `Bearer <TOKEN>` (use the access token received from `/login`)
-   **Response (Success - 200 OK):**

    ```json
    {
        "message": "Logout successful",
        "remaining_balance": 900
    }
    ```

-   **Response (Error - 401 Unauthorized):**

    ```json
    {
        "error": "invalid token"
    }
    ```
-   **Response (Error - 500 Internal Server Error):**
    ```json
    {
        "error": "Failed to get user balance"
    }
    ```

## Testing

You can use Postman, curl, or other tools for testing this API.

**Example using Postman:**

1.  Import the Postman collection (provided in a previous response or create your own based on the API documentation).
2.  Create an environment variable named `token` in Postman.
3.  Send a request to `/login` to get an access token. Copy the `access_token` value to the `token` environment variable.
4.  Send requests to `/payment` and `/logout` including the header `Authorization: Bearer {{token}}`.

## Deployment with Docker

Here's a basic guide on how to deploy this application using Docker:

1.  **Create a `Dockerfile`:**

    ```dockerfile
    # Use the official Golang image as the base image
    FROM golang:1.21

    # Set the working directory inside the container
    WORKDIR /app

    # Copy go.mod and go.sum files to the container
    COPY go.mod go.sum ./

    # Download all dependencies
    RUN go mod download

    # Copy the rest of the application code to the container
    COPY . .

    # Build the Go application
    RUN go build -o main .

    # Expose port 8080
    EXPOSE 8080

    # Define the command to run the executable
    CMD ["./main"]
    ```

2.  **Build the Docker image:**

    ```bash
    docker build -t merchant-bank-api .
    ```

3.  **Run the Docker container:**

    ```bash
    docker run -p 8080:8080 -d --name merchant-bank-container -v $(pwd)/data:/app/data -e JWT_SECRET_KEY=YOUR_STRONG_SECRET_KEY merchant-bank-api
    ```

    -   `-p 8080:8080`: Maps port 8080 of the container to port 8080 on the host machine.
    -   `-d`: Runs the container in detached mode (in the background).
    -   `--name merchant-bank-container`: Assigns a name to the container.
    -   `-v $(pwd)/data:/app/data`: Mounts the local `data` directory to the `/app/data` directory inside the container (for data persistence).
    -   `-e JWT_SECRET_KEY=YOUR_STRONG_SECRET_KEY`: Sets the `JWT_SECRET_KEY` environment variable inside the container. **Replace `YOUR_STRONG_SECRET_KEY` with your actual secret key.**

    **Note:** You might need to adjust the Docker commands based on your operating system and environment.
