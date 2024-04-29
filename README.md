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
![User Registration](![register](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/9e42e2eb-0f84-4b10-b1f6-6ec77ff5ab3a))

### User Login
![User Login](![login](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/047f373a-bf78-472d-9e8d-3f42cfb4f065))

### Chat Interface
![Chat Interface](![chat](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/592e8edc-5036-4477-82ca-d0d7a2a91a89))

### Add New Contact
![Add New Contact](![contact-added](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/4978b62b-18e5-4413-86bc-20d181dd15bc))

### User Profile
![User Profile](![profile](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/5e4adb66-77ae-436a-928a-ac9702ad669a))

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
