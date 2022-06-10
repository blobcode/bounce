# bounce

a simple redirect server written in go

## api spec

`POST /new` - creates a new link. Requires `url` attribute in body with a valid uri. Returns `json` with a `id` attribute.
`GET /r/<id>` - redirect to a link given `id`. Returns 404 if invalid
