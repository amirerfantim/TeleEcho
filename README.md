# TeleEcho Messenger

## Overview

TeleEcho is a messaging platform built with a Golang backend and a React frontend. It features user registration, login, contact management, group creation, and real-time messaging.

## Technologies Used

- Backend: Golang
- Frontend: React, TypeScript
- Database: PostgreSQL
- Storage: Integration with object storage for media
- Authentication: JWT for secure API access

## Key Features

- User registration with profile information
- Secure user login
- Real-time one-on-one and group messaging
- Contact management
- Group chat creation and management

## Screenshots

### User Registration
<div style="display: flex; align-items: center;">
  <img src="https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/28a64be2-f839-4cd7-9f38-94c27473a6c2" alt="Chat" width="800">
</div>

### User Login
<div style="display: flex; align-items: center;">
  <img src="https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/7aae7cf5-1d14-455e-b95b-d886f01426ff" alt="Chat" width="800">
</div>

### User Profile
<div style="display: flex; align-items: center;">
  <img src="https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/2e574838-c399-47ed-8303-413207d099e9" alt="Profile" width="1000">
</div>

### Chat Interface
<div style="display: flex; align-items: center;">
  <img src="https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/75413591-0246-48d6-b3f5-e3af627721f5" alt="Chat" width="1000">
</div>

## Setup Instructions

1. **Clone the Repository**
  ```
  git clone https://github.com/Mohammad-Rahmanian/TeleEcho
  ```

2. **Database Setup**
  ```
  docker run --name postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=messenger-db -p 5432:5432 -d postgres
  ```

3. **Configuration File**

   Create a `config.yaml` in the `Backend/configs` with the necessary configurations. A template for the `config.yaml` file:
  ```yaml
  port: 8020
  address: localhost
  token-key: 'your-secret-jwt-key-here'
  database-port: 5432
  database-address: localhost
  database-user: 'database-username-here'
  database-password: 'database-password-here'
  database-name: 'messenger-db'
  storage-service-id: 'your-storage-service-id-here'
  storage-service-secret: 'your-storage-service-secret-here'
  storage-service-endpoint: 'https://your-storage-service-endpoint-here'
  storage-service-bucket: 'your-storage-service-bucket-name-here'
  ```

4. **Backend Setup**
  ```
  cd Backend
  go get ./...
  ```

5. **Frontend Setup**
  ```
  cd ../Frontend
  npm install
  ```

6. **Running the Application**
- For Backend:
  ```
  go run main.go
  ```
- For Frontend:
  ```
  npm run dev
  ```
