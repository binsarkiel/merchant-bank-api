# Merchant Bank API

This document provides a comprehensive guide to the Merchant Bank API, detailing its functionalities, specifications, and usage instructions.

## Table of Contents

*   [Technical Specifications](#technical-specifications)
*   [Software Architecture](#software-architecture)
*   [Project Structure](#project-structure)
*   [Getting Started](#getting-started)
    *   [Prerequisites](#prerequisites)
    *   [Installation](#installation)
    *   [Setting up Initial Users](#setting-up-initial-users)
*   [Running the Application](#running-the-application)
*   [API Endpoints](#api-endpoints)
    *   [Login (POST /login)](#login-post-login)
    *   [Payment (POST /payment)](#payment-post-payment)
    *   [Logout (DELETE /logout)](#logout-delete-logout)
*   [Testing with Postman](#testing-with-postman)
*   [Deployment with Docker](#deployment-with-docker)

## Technical Specifications

*   **Programming Language:** Go (Golang)
*   **Framework:** Gin Gonic
*   **Database:** JSON file-based (users.json, transactions.json, sessions.json)
*   **Authentication:** JWT (JSON Web Token) with Bearer scheme
*   **Password Hashing:** bcrypt
*   **UUID Generator:** [github.com/google/uuid](https://github.com/google/uuid)

## Software Architecture

The Merchant Bank API follows a layered architecture, promoting separation of concerns and maintainability. The components interact as follows:

![Architecture](https://i.imgur.com/Lsun93Z.png)

*   **Client:** Any client that can send HTTP requests (e.g., Postman, web browser, mobile app).
*   **API Layer (handlers.go):**
    *   Handles incoming HTTP requests.
    *   Uses Gin Gonic framework for routing.
    *   Performs initial request validation (e.g., checking for required parameters, validating the `Bearer` token format).
    *   Calls the appropriate functions in the Service Layer.
    *   Returns HTTP responses to the client.
*   **Service Layer (service.go):**
    *   Contains the core business logic of the application.
    *   Handles user authentication (login, logout, token validation, invalidating tokens).
    *   Implements payment processing logic.
    *   Interacts with the Repository Layer to access and manipulate data.
    *   Implements password hashing using bcrypt.
*   **Repository Layer (repository.go):**
    *   Responsible for direct interaction with the data store (JSON files in this case).
    *   Provides functions to load, save, and query data from the JSON files.
    *   Uses mutexes to ensure data consistency during concurrent access.
*   **Data Store:**
    *   Consists of three JSON files:
        *   `users.json`: Stores user information (ID, name, username, hashed password, account type, account balance).
        *   `transactions.json`: Stores transaction history.
        *   `sessions.json`: Stores user session activity (login/logout timestamps).


## Project Structure

| Directory/File                | Description                                                                                                                      |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `api/`                        | Package for API handlers.                                                                                                         |
| `api/handlers.go`             | Handlers for each API endpoint.                                                                                                  |
| `data/`                       | Directory for JSON data files.                                                                                                   |
| `data/sessions.json`          | JSON file for session data.                                                                                                      |
| `data/transactions.json`      | JSON file for transaction data.                                                                                                  |
| `data/users.json`             | JSON file for user data.                                                                                                         |
| `go.mod`                      | Go modules file.                                                                                                                |
| `go.sum`                      | Go modules checksums.                                                                                                             |
| `main.go`                     | Main application file.                                                                                                           |
| `models/`                     | Package for data models.                                                                                                          |
| `models/models.go`            | Data structures (User, Transaction, Session, request/response models).                                                            |
| `repository/`                 | Package for data access.                                                                                                          |
| `repository/repository.go`    | Functions for interacting with JSON files.                                                                                       |
| `services/`                   | Package for business logic.                                                                                                       |
| `services/service.go`         | Logic for authentication, payment, token validation, token invalidation, and password hashing using bcrypt.                        |
| `services/service_test.go`    | Unit tests for the service layer.                                                                                                  |
| `utils/`                      | Package for utility functions.                                                                                                    |
| `utils/utils.go`              | Helper functions (error handling).                                                                                          |
| `.env.example`                | Example environment variables file (JWT secret key).                                                                              |
## Getting Started

### Prerequisites

*   Go (version 1.20 or later) installed on your system.
*   Git (for cloning the repository).
*   Postman (for API testing).

### Installation

1.  **Clone the Repository:**

    ```bash
    git clone <repository_url>
    cd merchant-bank-api
    ```

2.  **Install Dependencies:**

    ```bash
    go mod tidy
    go mod download
    ```

### Setting up Initial Users

This application uses JSON files as a database. When you first run the application, it will check if `data/users.json` is empty. If it is, it will automatically create three initial users for testing purposes:

*   **John Doe** (customer) - username: `johndoe`, password: `password123`
*   **Jane Smith** (merchant) - username: `janesmith`, password: `password456`
*   **Peter Jones** (customer) - username: `peterjones`, password: `password789`

Their passwords are automatically hashed using bcrypt during this initialization. You can find this logic in the `setupInitialUsers()` function in `main.go`.

**Important:** These are example users for testing only. For a production environment, you should implement a proper user registration and management system.

## Running the Application

1.  **Create `.env` file:**

    *   Copy the contents of `.env.example` to a new file named `.env` in the project's root directory.
    *   **Important:** Replace `JWT_SECRET_KEY=your_secret_key_here` with a strong, secret key. **Do not use `your_secret_key_here` in a production environment.**

    ```
    JWT_SECRET_KEY=your_strong_secret_key
    ```

2.  **Ensure Data Directory and Files Exist:**

    *   Make sure the `data` directory exists in the project's root directory.
    *   Inside the `data` directory, ensure the following files exist and contain an empty JSON array (`[]`):
        *   `data/users.json`
        *   `data/transactions.json`
        *   `data/sessions.json`

    If these files don't exist, the application will create them on startup.

3.  **Run the Application:**

    ```bash
    go run main.go
    ```

    The application will start running on `http://localhost:8080`.

## API Endpoints

### Login (POST /login)

Authenticates a user and returns a JWT token.

**Request:**

*   **Method:** `POST`

*   **URL:** `/login`

*   **Headers:**

    *   `Content-Type`: `application/json`

*   **Body:**

    ```json
    {
        "username": "<username>",
        "password": "<password>"
    }
    ```

    *   `username`: User's username (string, required).
    *   `password`: User's password (string, required).

**Response:**

*   **Success (200 OK):**

    ```json
    {
        "token": "<jwt_token>"
    }
    ```

    *   `token`: JWT token for authentication.

*   **Unauthorized (401 Unauthorized):**

    ```json
    {
        "error": "Invalid username or password"
    }
    ```

*   **Bad Request (400 Bad Request):**
    * Returned if the request body is invalid or missing required fields.

    ```json
    {
        "error": "invalid request body"
    }
    ```

*   **Internal Server Error (500 Internal Server Error):**
     * Returned if there is an error on the server side.

     ```json
     {
         "error": "internal server error"
     }
     ```

### Payment (POST /payment)

Performs a money transfer between two users.

**Request:**

*   **Method:** `POST`

*   **URL:** `/payment`

*   **Headers:**

    *   `Content-Type`: `application/json`
    *   `Authorization`: `Bearer <jwt_token>`

*   **Body:**

    ```json
    {
        "recipient": "<recipient_username>",
        "amount": <amount>
    }
    ```

    *   `recipient`: Username of the recipient (string, required).
    *   `amount`: Amount of money to transfer (float64, required).

**Response:**

*   **Success (200 OK):**

    ```json
    {
        "message": "Transfer success!",
        "transaction": {
            "activity": "transfer_money",
            "transaction_id": "<transaction_uuid>",
            "sender": "<sender_username>",
            "recipient": "<recipient_username>",
            "amount": <amount>,
            "created_at": "<timestamp>"
        }
    }
    ```

*   **Unauthorized (401 Unauthorized):**

    ```json
    {
        "error": "<error_message>"
    }
    ```

    Examples:

    ```json
    {
        "error": "Authorization header required"
    }
    ```

    ```json
    {
        "error": "Invalid Authorization header format"
    }
    ```

    ```json
    {
        "error": "Invalid token"
    }
    ```
    
    ```json
    {
       "error": "Token has expired"
    }
    ```

*   **Bad Request (400 Bad Request):**
     * Returned if the request body is invalid or missing required fields, or if the sender and recipient are the same, or if the sender has insufficient balance, or if the sender or recipient is not found.

    ```json
    {
        "error": "<error_message>"
    }
    ```

    Examples:

    ```json
    {
        "error": "invalid request body"
    }
    ```

    ```json
    {
        "error": "sender and recipient cannot be the same"
    }
    ```

    ```json
    {
        "error": "insufficient balance"
    }
    ```

    ```json
    {
        "error": "sender not found: user not found"
    }
    ```

    ```json
    {
        "error": "recipient not found: user not found"
    }
    ```

*   **Internal Server Error (500 Internal Server Error):**

    ```json
    {
        "error": "<error_message>"
    }
    ```

    Examples:

    ```json
    {
        "error": "failed to update sender's balance: user not found"
    }
    ```

    ```json
    {
       "error": "failed to update recipient's balance: user not found"
    }
    ```

    ```json
    {
        "error": "failed to record transaction: ..."
    }
    ```

### Logout (DELETE /logout)

Logs out the user and invalidates the JWT token.

**Request:**

*   **Method:** `DELETE`
*   **URL:** `/logout`
*   **Headers:**
    *   `Authorization`: `Bearer <jwt_token>`

**Response:**

*   **Success (200 OK):**

    ```json
    {
        "message": "Logout success!",
        "remaining_balance": <remaining_balance>
    }
    ```

    *   `remaining_balance`: The user's remaining balance after logout (float64).

*   **Unauthorized (401 Unauthorized):**

    ```json
    {
       "error": "<error_message>"
    }
    ```

    Examples:

    ```json
    {
        "error": "Authorization header required"
    }
    ```

    ```json
    {
        "error": "Invalid Authorization header format"
    }
    ```

    ```json
    {
        "error": "Invalid token"
    }
    ```
    
    ```json
    {
       "error": "Token has expired"
    }
    ```

*   **Internal Server Error (500 Internal Server Error):**

    ```json
    {
        "error": "<error_message>"
    }
    ```

## Testing with Postman

1.  **Login:**

    *   Create a `POST` request to `/login` with a valid username and password in the request body.
    *   Copy the `token` from the response.

2.  **Payment:**

    *   Create a `POST` request to `/payment`.
    *   Add an `Authorization` header with the value `Bearer <token>` (replace `<token>` with the token from the login response).
    *   Provide the recipient's username and the amount in the request body.

3.  **Logout:**

    *   Create a `DELETE` request to `/logout`.
    *   Add an `Authorization` header with the value `Bearer <token>` (use the same token).

4.  **Test Token Invalidation:**

    *   Try to make another `POST` request to `/payment` using the same token after logging out. You should receive a `401 Unauthorized` error with the message `Invalid token`.

## Deployment with Docker

To deploy this application using Docker, you can follow these basic steps:

1.  **Create a Dockerfile:**

    Create a file named `Dockerfile` in the root directory of your project with the following content:

    ```dockerfile
    # Use the official Golang base image.
    FROM golang:1.20

    # Set the working directory inside the container.
    WORKDIR /app

    # Copy the local package files to the container's workspace.
    COPY . .

    # Download all dependencies.
    RUN go mod download

    # Build the Go app.
    RUN go build -o main .

    # Expose port 8080 to the outside.
    EXPOSE 8080

    # Command to run the executable.
    CMD ["./main"]
    ```

2.  **Build the Docker Image:**

    ```bash
    docker build -t merchant-bank-api .
    ```

    This command creates a Docker image named `merchant-bank-api`.

3.  **Run the Docker Container:**

    ```bash
    docker run -p 8080:8080 -e JWT_SECRET_KEY="your_actual_secret_key" --name merchant-bank-container merchant-bank-api
    ```

    *   `-p 8080:8080`: Maps port 8080 of the container to port 8080 on your host machine.
    *   `-e JWT_SECRET_KEY="your_actual_secret_key"`: Sets the `JWT_SECRET_KEY` environment variable inside the container. **Replace `"your_actual_secret_key"` with a strong secret key.**
    *   `--name merchant-bank-container`: Set name for your container.

4. **Verify**
   * You can verify if the container running or not by running this command `docker ps`

5. **Stop Container**
   * If you want to stop the container, you can run this command `docker stop merchant-bank-container`

6. **Remove Container**
    *   If you want to remove the container, you can run this command `docker rm merchant-bank-container`

7.  **Remove Docker Image**
    *   If you want to remove the docker image, you can run this command `docker rmi merchant-bank-api`
