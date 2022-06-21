# bounce

a simple redirect server written in go

## running

clone the repo and run `go build`

bounce requires a `bounce.ini` config file ([example](./bounce.ini))

next, just run the binary and visit on the configured `port`.

## api spec

`POST /new` - creates a new link. Requires `url` attribute in body with a valid uri. Returns `json` with a `id` attribute.
`GET /r/<id>` - redirect to a link given `id`. Returns 404 if invalid
