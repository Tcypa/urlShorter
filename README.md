# URL Shortener

A simple URL shortening api.  
Supports **In-Memory** and **PostgreSQL** storage options.

##  API Endpoints

- **`POST /shorten`** — Accepts JSON `{ "url": "https://example.com" }`, returns shorten Url.
- **`GET /{shortURL}`** — Redirects to the original Url or returns 404.

## Run the Service

**Using In-Memory Storage:**

Run by default or use:
```sh
go run main.go -stgType=memory
```


**Using Postgres Storage:**
```sh
go run main.go -stgType=postgres
```