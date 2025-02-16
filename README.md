# Chirpy - A Twitter Clone Backend API

Chirpy is a backend API built using Golang, SQLC, and PostgreSQL. It serves as the foundation for a Twitter-like social platform where users can create and interact with short text-based posts called "chirps."

## Features

- User authentication (signup, login, JWT-based authentication)
- CRUD operations for chirps (create, read, update, delete)
- Follow/unfollow functionality
- Like and reply to chirps
- Feed generation based on followed users
- Secure and efficient SQL queries using SQLC

## Tech Stack

- **Golang** - The core language for backend development
- **PostgreSQL** - The relational database for storing data
- **SQLC** - Type-safe, efficient query generation for PostgreSQL
- **JWT** - JSON Web Token for authentication
- **Gin** - HTTP router and middleware for Go

## Getting Started

### Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/dl/) (v1.19+ recommended)
- [PostgreSQL](https://www.postgresql.org/download/) (v13+ recommended)
- [SQLC](https://sqlc.dev/) for generating query bindings

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Prodigy00/chirpy.git
   cd chirpy
   ```

2. Set up environment variables:
   ```sh
   cp .env.example .env
   ```
   Then update `.env` with your PostgreSQL credentials.

3. Install dependencies:
   ```sh
   go mod tidy
   ```

4. Generate SQL bindings:
   ```sh
   sqlc generate
   ```

5. Run database migrations:
   ```sh
   make migrate-up
   ```

6. Start the server:
   ```sh
   go run main.go
   ```

The API will be available at `http://localhost:8080`.

## API Endpoints

### Authentication

| Method | Endpoint        | Description          |
|--------|---------------|----------------------|
| POST   | `/signup`     | Register a new user |
| POST   | `/login`      | Authenticate a user |

### Chirps

| Method | Endpoint         | Description               |
|--------|-----------------|---------------------------|
| POST   | `/chirps`       | Create a new chirp       |
| GET    | `/chirps/:id`   | Get a chirp by ID        |
| DELETE | `/chirps/:id`   | Delete a chirp           |
| GET    | `/feed`         | Get a feed of chirps     |

More detailed API documentation will be available soon.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m 'Add new feature'`).
4. Push to your branch (`git push origin feature-name`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries or suggestions, feel free to reach out:

- GitHub: [Prodigy00](https://github.com/Prodigy00)
- Email: [your-email@example.com] (replace with actual email if needed)

---

Happy chirping! 🐦

