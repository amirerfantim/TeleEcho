# TeleEcho Messenger

## Overview

TeleEcho is a comprehensive messaging platform built with a Golang backend and a React frontend. It features user registration, login, contact management, group creation, and real-time messaging, all secured with JWT for authentication.

## Technologies Used

- Backend: Golang
- Frontend: React, TypeScript
- Database: PostgreSQL
- Storage: Integration with object storage for media and message data
- Authentication: JWT for secure API access

## Features

- User registration with profile information
- Secure user login
- Real-time one-on-one and group messaging
- Contact management
- Group chat creation and management
- Media file sharing

## Screenshots

### User Registration
![User Registration](![register](![register](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/28a64be2-f839-4cd7-9f38-94c27473a6c2)))

### User Login
![User Login](![login](![login](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/7aae7cf5-1d14-455e-b95b-d886f01426ff)))

### User Profile
![User Profile](![profile](![profile](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/2e574838-c399-47ed-8303-413207d099e9)))

### Chat Interface
![Chat Interface](![chat](![chat](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/75413591-0246-48d6-b3f5-e3af627721f5)))


## Setup Instructions

1. **Clone the Repository**
git clone <repository-url>

2. **Database Setup**
docker run --name postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=messenger-db -p 5432:5432 -d postgres

4. **Configuration File**
- Create a `config.yaml` in the `Backend/configs` with the necessary configurations.

4. **Backend Setup**
cd Backend
go get ./...

5. **Frontend Setup**
cd ../Frontend
npm install

6. **Running the Application**
- For Backend:
  ```
  go run main.go
  ```
- For Frontend:
  ```
  npm start
  ```
\
