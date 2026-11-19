# TeleEcho Messenger

![Go 1.23+](https://img.shields.io/badge/go-1.23%2B-00ADD8.svg)
![React 18](https://img.shields.io/badge/react-18-61DAFB.svg)
![TypeScript](https://img.shields.io/badge/typescript-3178C6.svg)

TeleEcho is a real-time messenger with a Go backend and a React frontend. Users
register, add contacts, and chat one-to-one or in groups. Messages travel over
WebSockets, and the rest of the API is plain HTTP with JWT authentication.

**Contents:** [Features](#features) · [Screenshots](#screenshots) · [Setup](#setup) ·
[API](#api) · [Project layout](#project-layout) · [Authors](#authors)

## Features

- Registration and login. Passwords are hashed before they are stored, and the
  API is protected with JWT.
- User profiles with a profile picture and a short bio.
- Contacts: add, list and remove contacts, and change a contact's status.
- Direct chats with real-time messages over WebSockets.
- Groups with a name, a description and a picture, with members you can add and
  remove, and group chats with real-time messages.
- Profile pictures are stored in S3-compatible object storage.

| | |
|---|---|
| Backend | Go, [Echo](https://echo.labstack.com/), GORM, PostgreSQL, Gorilla WebSocket, JWT |
| Frontend | React 18, TypeScript, Vite |
| Storage | S3-compatible object storage for profile pictures |

## Screenshots

### User registration
![User Registration](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/90c9786d-7fad-4194-ba5d-23ca7564043e)

### User login
![User Login](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/7aae7cf5-1d14-455e-b95b-d886f01426ff)

### User profile
![User Profile](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/1468cd71-d2ae-4e4d-9237-f901276819ec)

### Chat interface
![Chat Interface](https://github.com/Mohammad-Rahmanian/TeleEcho/assets/78559411/75413591-0246-48d6-b3f5-e3af627721f5)

## Setup

You need Go 1.23 or newer, Node.js and Docker (or any PostgreSQL server).

### 1. Clone the repository

```bash
git clone https://github.com/amirerfantim/TeleEcho
cd TeleEcho
```

### 2. PostgreSQL

```bash
docker run --name postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=pass \
  -e POSTGRES_DB=messenger-db -p 5432:5432 -d postgres
```

The backend creates its tables when it starts.

### 3. Object storage

Registration uploads the profile picture, so the backend needs an S3-compatible
storage service and a bucket. The storage client addresses the bucket as
`<bucket>.<endpoint host>` (virtual-hosted style) and has no option for
path-style addresses, so the service must support that. Amazon S3 and most cloud
object storage services do.

### 4. Backend

Create `Backend/configs/config.yaml`. The backend reads this exact path from the
`Backend` directory, and every key is required:

```yaml
port: 8020
address: localhost
token-key: 'your-secret-jwt-key-here'
database-port: 5432
database-address: localhost
database-user: 'admin'
database-password: 'pass'
database-name: 'messenger-db'
storage-service-id: 'your-storage-service-id-here'
storage-service-secret: 'your-storage-service-secret-here'
storage-service-endpoint: 'https://your-storage-service-endpoint-here'
storage-service-bucket: 'your-storage-service-bucket-name-here'
```

Then install the dependencies and run it:

```bash
cd Backend
go get ./...
go run main.go
```

`go get ./...` is needed: `go.mod` does not list every package the code imports.
It also needs Go 1.23 or newer, because the current `gorilla/sessions` requires it.

Keep `port: 8020`. The frontend calls `http://127.0.0.1:8020`, and the backend
only accepts cross-origin requests from `http://localhost:5173`.

### 5. Frontend

```bash
cd Frontend
npm install
npm run dev
```

Open <http://localhost:5173>. `npm run build` makes a production build.

## API

Routes marked JWT need the token from `/login` in the `Authorization` header,
**without** a `Bearer ` prefix (the backend rejects the prefix with 401). Request
data is sent as form fields or in the query string.

| Route | Auth | Parameters |
|---|---|---|
| `POST /register` | none | form: `username`, `firstname`, `lastname`, `phone`, `password` (8 or more characters), `bio` (up to 200 characters), `profile` (a picture file, required, up to 10 MB) |
| `POST /login` | none | form: `username`, `password`. Returns a JWT. |
| `GET`, `PATCH`, `DELETE /users` | JWT | `GET` takes `username`; `PATCH` takes the fields of `/register`; `DELETE` removes the signed-in user |
| `POST`, `GET`, `DELETE`, `PATCH /contacts` | JWT | `phone` to add; `id` to list; `username` to remove; `username` and `status` to change status |
| `POST`, `GET`, `PATCH`, `DELETE /group` | JWT | `name`, `description`, `profile` to create; `name` to search; `groupID` and `username` to add or remove a member |
| `GET /group/all` | JWT | `groupID`: all users of a group |
| `POST`, `GET /group-chat` | JWT | `groupID` to create; `chatID` to read messages |
| `POST /chat` | JWT | `receiverID`: start a direct chat |
| `GET /chat/message`, `DELETE /chat/message` | JWT | `chatID` to read; `id` to delete a message |
| `DELETE /chat/:chatID` | JWT | Delete a direct chat |

WebSocket routes take the token in the query string, as `?token=<JWT>`:

| Route | Purpose |
|---|---|
| `/all-chat` | The user's direct chats |
| `/write-message?chatID=<id>`, `/read-message?chatID=<id>` | Send and receive direct messages |
| `/all-group-chat` | The user's group chats |
| `/write-group-message`, `/read-group-message` | Send and receive group messages |

## Project layout

```
Backend/
  main.go            start-up: config, database, storage, server
  configs/           config file loader
  router/            routes and CORS
  middleware/        JWT check
  api/handlers/      HTTP and WebSocket handlers
  api/database/      PostgreSQL access (tables are created at start)
  api/services/      S3 storage client
  model/             database models
Frontend/
  src/components/    pages and forms (login, register, contacts, groups, chats)
```

## Authors

TeleEcho was built by [Amirerfan Teimouri](https://github.com/amirerfantim) and
[Mohammad Rahmanian](https://github.com/Mohammad-Rahmanian). This repository is a
fork of [Mohammad-Rahmanian/TeleEcho](https://github.com/Mohammad-Rahmanian/TeleEcho).

Contact: [LinkedIn](https://www.linkedin.com/in/amirerfantim/) ·
[amirerfantim@gmail.com](mailto:amirerfantim@gmail.com)
