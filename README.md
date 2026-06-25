# Shopping Cart API

A demo REST API for a shopping cart in a microservices architecture in `go`

## Navigate

1. [About](#about)
2. [Stack](#stack)
3. [API](#api)
4. [Content](#content)
5. [Restrictions](#restrictions)

## About

This project was created to update the portfolio and practice of working with microservice architecture.

## Stack

- **Language:** Go
- **HTTP-Framework:** Gin
- **Main storage:** PostgreSQL
- **Cart storage:** Valkey (open Redis analog)

## API

| Method | Path                 | Descpription                        |
|--------|----------------------|-------------------------------------|
| POST   | /users/              | Sign up                             |
| POST   | /users/auth          | Log in                              |
| GET    | /users/{login}       | Get user profile by login           |
| DELETE | /users/{login}       | Delete your profile                 |
| POST   | /products            | Create product                      |
| GET    | /products/{id}       | Get product information             |
| DELETE | /products/{id}       | Delete product if you owner         |
| GET    | /products/search     | Search product by title in query    |
| GET    | /cart                | Get information for your cart       |
| POST   | /cart/items          | Add product to cart                 |
| PUT    | /cart/items/quantity | Update quantity for product in cart |
| PUT    | /cart/items/remove   | Remove item from cart               |
| POST   | /cart/clear          | Full clean cart                     |
| GET    | /cart/ttl            | Get time to live for cart           |

## Content

- Authorization and Authentication with JWT token in cookie. But there is no email confirmation.
- Create and adding product to storage in PostgreSQL.
- Adding product to the cart.
- Get information for profile users, for products, for cart.

## Restrictions

- Absent order service.
- Not enough test. Now tests only data validation, need write mocks for services.
- Services communicate directly without a message broker.
- No monitoring and health check.
