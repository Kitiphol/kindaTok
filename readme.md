# SCA2024-T3-P2-Dying

This project is a scalable web application with a microservices backend (Go) and a Next.js frontend.

---

## Prerequisites

- [Node.js](https://nodejs.org/) (v18 recommended)
- [npm](https://www.npmjs.com/) (comes with Node.js)
- [Go](https://go.dev/) (v1.20+ recommended)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

---

## Frontend

The frontend is located in [`frontend/toktik_frontend`](frontend/toktik_frontend/).

### Install dependencies

```sh
cd frontend/toktik_frontend
npm install
```

### Run development server

```sh
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

---

## Backend

Each backend service is a Go microservice in [`Backend/`](Backend/).  
You need to install dependencies and run each service separately.

### Install Go dependencies

For each service (e.g., `Authentication`, `UserService`, `VideoService`, etc.):

```sh
cd Backend/<ServiceName>
go mod tidy
```

### Run a backend service

```sh
cd Backend/<ServiceName>
go run ./cmd/main.go
```

Replace `<ServiceName>` with the actual service folder name, such as `Authentication`, `UserService`, `VideoService`, etc.

---

## Docker (Optional)

You can also use Docker to build and run the services.  
See the provided `Dockerfile` in each service directory and the [`docker-compose.yml`](docker-compose.yml) file for orchestration.

---

## Notes

- Make sure to configure environment variables as needed (see `.env` files).
- For more details, check the individual `README.md` files in each service (if available).

---