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

To fetch unit details by sub-domain, send a GET request to
`/api/units/by_subdomain?sub_domain=yourname` without authentication. The
response includes the unit's name, logo and sub-domain. If the sub-domain does
not exist, the API returns an error message "sub domain not found".

When logging in, include the unit's sub-domain along with the username and password:

```json
{
  "username": "admin",
  "password": "admin123",
  "sub_domain": "admin"
}
```

The response includes the user's `unit_id` and an `is_admin` flag.
