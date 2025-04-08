# Golang API with Gin and MongoDB

## Description

This is a RESTful API built using the [Gin](https://github.com/gin-gonic/gin) framework and MongoDB. The API supports user authentication, role-based access control (RBAC), and CRUD operations for todos. Each user can manage their own todos, with additional functionality like pagination, counting todos, and searching by user ID.

This template can be extended to other APIs.

---

## Features

- User Authentication and Authorization with JWT.
- Role-Based Access Control (RBAC).
- CRUD Operations for Users and Todos.
- Pagination and Filtering for Todos.
- MongoDB indexing for improved query performance.
- Secure handling of sensitive data with `.env` file.
- Middleware for CORS, Compression, and Authentication.

---

## Project Structure

```plaintext
├── config/
│   └── config.go          # Database connection and environment variables
├── controllers/
│   ├── user.go # CRUD for users
│   └── todo.go # CRUD for todos
├── helpers/
│   └── jwt.go             # Helper functions for JWT handling
├── middleware/
│   └── auth.go # Middleware for authentication and authorization
├── models/
│   ├── user.go            # User model
│   └── todo.go            # Todo model
├── routes/
│   └── routes.go          # All API routes
├── utils/
│   └── validations.go     # Validation logic # coming soon
├── main.go                # Entry point of the application
├── .env                   # Environment variables (MongoDB URL, JWT Secret)
└── README.md              # Documentation
```

---

## Endpoints

### **Authentication**

| HTTP Method | Endpoint    | Description            |
| ----------- | ----------- | ---------------------- |
| `POST`      | `/login`    | Login and generate JWT |
| `POST`      | `/register` | Register a new user    |

### **User Management**

| HTTP Method | Endpoint    | Description         |
| ----------- | ----------- | ------------------- |
| `GET`       | `/user/:id` | Get a user by ID    |
| `PUT`       | `/user/:id` | Update user details |
| `DELETE`    | `/user/:id` | Delete a user       |

### **Todo Management**

| HTTP Method | Endpoint                     | Description                         |
| ----------- | ---------------------------- | ----------------------------------- |
| `POST`      | `/todo`                      | Create a new todo                   |
| `GET`       | `/todo/:id`                  | Get a todo by ID                    |
| `GET`       | `/todo/user/:userID`         | Get all todos for a specific user   |
| `GET`       | `/todo/paginated?page&limit` | Get todos with pagination           |
| `GET`       | `/todos/count/:userID`       | Count all todos for a specific user |
| `PUT`       | `/todo/:id`                  | Update a todo                       |
| `DELETE`    | `/todo/:id`                  | Delete a todo                       |

---

## Environment Variables

Create a `.env` file in the root directory and add the following variables:

```plaintext
MONGO_URL=mongodb://localhost:27017
JWT_SECRET=your_jwt_secret
```

---

## Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/clinton-mwachia/Go-starter-template.git
   cd Go-starter-template
   ```

2. **Install dependencies**  
   Make sure you have Go installed. Then run:

   ```bash
   go mod tidy
   ```

3. **Run the application**  
   Start the application using:

   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`.

---

## Usage

### **Postman**

- Import the API endpoints into Postman using the provided routes.
- Use the `POST /login` endpoint to generate a JWT token.
- Include the token in the `Authorization` header for all protected endpoints:
  ```plaintext
  Authorization: Bearer <your_jwt_token>
  ```

### **MongoDB Indexing**

Ensure MongoDB has the required indexes for optimal performance:

coming soon

---

## Testing

Comming soon

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## Update all packages to latest

```go
go get -u
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch.
3. Commit your changes.
4. Push to the branch.
5. Open a pull request.
