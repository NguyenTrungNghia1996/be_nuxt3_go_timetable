# Go Fiber Template

This project provides a minimal REST API built with
[Fiber](https://github.com/gofiber/fiber). Only user login and image
uploads to an S3 compatible service remain.

## Running locally

```bash
go run main.go
```

Create an `.env` file (see `env` for an example) containing your database
credentials.

## Postman Collection

To quickly explore the API you can import
`postman/go-fiber-template.postman_collection.json` into Postman. The collection
assumes two variables:

- `baseUrl` – base address of your running server, e.g. `http://localhost:4000`
- `token` – JWT obtained from the `Login` request

The collection contains examples for logging in and obtaining presigned URLs
for uploading images.
