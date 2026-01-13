# Securing a REST API with JWT in Go

This tutorial explains how to secure a REST API in Go using **JSON Web Tokens (JWT)**.

---

## Prerequisites

- Go 1.20+
- Basic knowledge of HTTP in Go
- `github.com/golang-jwt/jwt/v5`

Install the dependency:

```bash
go get github.com/golang-jwt/jwt/v5
```

---

## How JWT Authentication Works

1. User logs in with credentials  
2. Server validates the credentials  
3. Server generates a signed JWT  
4. Client stores the token  
5. Client sends the token with each request  
6. Middleware validates the token  

JWT authentication is **stateless**, meaning no session storage is required.

---

## Creating a JWT on Login

```go
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Username != OurUser.Username || req.Password != OurUser.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized) // validating the credentials
		return
	}

	claims := jwt.MapClaims{
		"sub": req.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // generating a new token
	signed, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": signed,
		"token_type":   "bearer",
	})
}
```

---

## JWT Claims Explained

| Claim | Description |
|------|------------|
| `sub` | User identifier |
| `iat` | Issued at |
| `exp` | Expiration time |


### What are JWT claims?

Claims are pieces of information stored inside a JWT.
They are key–value pairs that describe the authenticated user and the token itself.


---

## Sending the Token from the Client

```
Authorization: Bearer <JWT_TOKEN>
```

---

## JWT Authentication Middleware

A **middleware function** is a piece of code that runs between an incoming request and the final handler, allowing you to inspect, modify, or block the request before it reaches the endpoint (for example, to handle authentication, logging, or validation).
In our case we use the next middleware function to check if the token is present and if it's valid. In other case, we will reject the request and send an `401 Unauthorized` response.

```go
type contextKey string

const ClaimsContextKey contextKey = "claims"

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { // parsing the token using the JWT key
			if t.Method.Alg()[:2] != "HS" {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

---

## Accessing Claims in Handlers

```go
/*
	Here we extract the claims from the request and converting them to jwt.MapClaims type.
	After that, we extract the username from the claims. (i.e. claims["sub"])
*/
claims := r.Context().Value(ClaimsContextKey).(jwt.MapClaims)
username := claims["sub"].(string)
```

---

## Protecting Routes (IMPORTANT)

It is important to note that a middleware does **not** run automatically.
If it’s not attached to the route, it won’t execute.

### ❌ Wrong (middleware not applied)

```go
http.HandleFunc("/protected", ProtectedHandler)
```

### ✅ Correct

```go
http.Handle(
	"/protected",
	jwtAuthMiddleware(http.HandlerFunc(ProtectedHandler)),
)
```

---

## JWT Secret

Use the same secret everywhere: (and for the better, DON'T STORE IT IN THE SOURCE CODE)

```go
var JwtKey = []byte("supersecretkey")
```

Store secrets in environment variables in production.

---

## Other Security Best Practices

- Use HTTPS
- Set expiration times
- Never store sensitive data in JWT claims
- Keep secrets out of source control
