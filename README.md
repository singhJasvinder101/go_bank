# Go Bank - Backend System for Transaction Handling

## Overview
Go Bank is a backend system designed to handle transactions efficiently. This project is built using **Golang**, **PostgreSQL**, and **Docker**, with a strong focus on scalability and robustness. Future plans to include adding advanced features such as **gRPC**, **Kubernetes**, and **CI/CD automation**.

## Features 
- [x] **Database Schema Design** - Structured using **PostgreSQL** with migrations.
- [x] **Dockerized PostgreSQL Instance** - Running the database inside a container.
- [x] **Database Migrations** - Using **migrate** to handle schema vrsion control.
- [x] **SQL Queries with sqlc** - Generating type-safe queries.
- [x] **API Development** - Implementing APIs.
- [ ] **User Authentication** - Adding JWT-based authentication.
- [x] **Transaction Handling** - Ensuring atomicity and consistency.
- [x] **Unit & Integration Testing** - Writing test cases for reliability.
- [ ] **gRPC Support** - Implementing efficient communication between microservices.
- [ ] **Kubernetes Deployment** - Managing scalability and orchestration.
- [x] **CI/CD Pipelines** - Automating testing and deployments.
- [ ] **Monitoring & Logging** - Integrating tools like Prometheus and Grafana.

## Setup Instructions 
1. **Clone the repository:**
   ```sh
   git clone https://github.com/singhJasvinder101/go_bank
   cd go_bank
   ```
2. **Start PostgreSQL container:**
   ```sh
   make postgres
   ```
3. **Create the database:**
   ```sh
   make createDB
   ```
4. **Run database migrations:**
   ```sh
   make migrateUp
   ```
5. **Generate SQL queries:**
   ```sh
   make sqlc
   ```

Stay tuned for more updates! 🚀
