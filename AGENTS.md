When you change API routes or controllers, update the Postman collection file `postman/go-fiber-template.postman_collection.json` so it reflects all available endpoints. Include example request bodies for create and update operations when relevant.
Always run `go vet ./...` and `go test ./...` after making changes and report the results.

All GET list API endpoints must accept optional `page`, `limit` and `search` query parameters. If `page` or `limit` are omitted or set to `0`, the full list should be returned without pagination. Ensure the Postman collection includes these parameters for each list endpoint.
