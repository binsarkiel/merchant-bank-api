# Merchant-Bank API

API sederhana untuk transaksi antara pedagang (merchant) dan pelanggan (customer) yang disimulasikan dengan interaksi ke database file JSON.

## Fitur

-   **Login:** User (merchant atau customer) dapat login menggunakan username dan password.
-   **Payment:** User yang sudah login dapat melakukan transfer uang ke user lain (yang terdaftar) menggunakan username penerima.
-   **Logout:** User dapat logout dan token akan di-invalidate.
-   **Authentication:** Menggunakan JWT (JSON Web Token) untuk autentikasi dan otorisasi.
-   **Simulasi Database:** Menggunakan file JSON (`users.json`, `sessions.json`, `transactions.json`) untuk simulasi penyimpanan data.

## Teknologi yang Digunakan

-   **Go:** Bahasa pemrograman utama.
-   **Gin:** Web framework untuk routing dan handling HTTP requests.
-   **JWT:** Untuk autentikasi dan otorisasi.
-   **`encoding/json`:** Untuk encoding dan decoding data JSON.
-   **`os`:** Untuk interaksi dengan file system.
-   **`github.com/google/uuid`:** Untuk generate UUID (untuk ID transaksi).
-   **`github.com/joho/godotenv`:** Untuk membaca environment variables dari file `.env`.

## Cara Menjalankan

### Prerequisites

-   **Go:** Pastikan Go sudah terinstall di sistem Anda (versi 1.20 ke atas direkomendasikan).
-   **Git:** Untuk meng-clone repository.

### Langkah-langkah

1.  **Clone repository:**

    ```bash
    git clone <repository_url>
    cd merchant-bank-api
    ```

2.  **Install dependencies:**

    ```bash
    go mod download
    ```

3.  **Buat file `.env`:**

    -   Buat file dengan nama `.env` di root direktori project.
    -   Isi file `.env` dengan `JWT_SECRET_KEY`:

        ```
        JWT_SECRET_KEY=ganti_dengan_secret_key_anda
        ```

        **Penting:** Ganti `ganti_dengan_secret_key_anda` dengan secret key yang kuat dan aman. **Jangan gunakan secret key ini untuk production.**

4.  **Buat file JSON (jika belum ada):**

    -   Jika folder `data` dan file-file JSON di dalamnya belum ada, buat secara manual atau gunakan command berikut di terminal:

        ```bash
        mkdir data
        touch data/users.json data/sessions.json data/transactions.json
        echo "[]" > data/users.json
        echo "[]" > data/sessions.json
        echo "[]" > data/transactions.json
        ```

    -   Isi `data/users.json` dengan data user awal. Contoh:

        ```json
        [
            {
                "id": "1",
                "name": "John Doe",
                "username": "johndoe",
                "password": "password123",
                "account_type": "customer",
                "account_balance": 1000
            },
            {
                "id": "2",
                "name": "Jane Smith",
                "username": "janesmith",
                "password": "password456",
                "account_type": "merchant",
                "account_balance": 5000
            }
        ]
        ```

5.  **Jalankan aplikasi:**

    ```bash
    go run main.go
    ```

    Server akan berjalan di port `8080`.

## Endpoint API

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
    | `username` | string | Username user (required)                      |
    | `password` | string | Password user (required)                      |

-   **Response (Success - 200 OK):**

    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

    | Field   | Type   | Description           |
    | :------ | :----- | :-------------------- |
    | `token` | string | JWT token for authentication |

-   **Response (Error - 401 Unauthorized):**

    ```json
    {
        "error": "invalid password"
    }
    ```

    ```json
    {
        "error": "user not found"
    }
    ```

-   **Response (Error - 400 Bad Request):**

    ```json
    {
        "error": "..." // Error message related to invalid request body
    }
    ```

### 2. `/payment`

-   **Method:** `POST`
-   **Headers:**
    -   `Authorization`: `Bearer <TOKEN>` (gunakan token yang didapat dari `/login`)
-   **Request Body:**

    ```json
    {
        "recipient": "janesmith",
        "amount": 100
    }
    ```

    | Field       | Type    | Description                                       |
    | :---------- | :------ | :------------------------------------------------ |
    | `recipient` | string  | Username penerima (required)                     |
    | `amount`    | number  | Jumlah uang yang ditransfer (required)            |

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
        "error": "Invalid token"
    }
    ```

    ```json
    {
        "error": "Authorization header required"
    }
    ```

    ```json
    {
        "error": "token expired"
    }
    ```

-   **Response (Error - 500 Internal Server Error):**

    ```json
    {
        "error": "insufficient balance"
    }
    ```

    ```json
    {
        "error": "recipient not found"
    }
    ```

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

    ```json
    {
        "error": "failed to record transaction"
    }
    ```

-   **Response (Error - 400 Bad Request):**

    ```json
    {
        "error": "..." // Error message related to invalid request body
    }
    ```

### 3. `/logout`

-   **Method:** `POST`
-   **Headers:**
    -   `Authorization`: `Bearer <TOKEN>` (gunakan token yang didapat dari `/login`)
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
        "error": "Invalid token"
    }
    ```

-   **Response (Error - 500 Internal Server Error):**

    ```json
    {
        "error": "Failed to get user balance"
    }
    ```

## Testing

Anda dapat menggunakan Postman, curl, atau tools lain untuk testing API ini.

**Contoh menggunakan Postman:**

1.  Import collection Postman (disediakan di jawaban sebelumnya).
2.  Buat environment variable `token` di Postman.
3.  Lakukan request ke `/login` untuk mendapatkan token. Token akan otomatis tersimpan di environment variable `token`.
4.  Lakukan request ke `/payment` dan `/logout` dengan menyertakan header `Authorization: Bearer {{token}}`.
