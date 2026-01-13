# Securing a REST API with JWT in Go

This tutorial explains how to secure a REST API in Go using JSON Web Tokens (JWT). We’ll cover:
- Creating a JWT on login
- Understanding JWT claims
- Protecting routes using middleware
- A common mistake where middleware appears not to work (and how to fix it)

---

## Prerequisites
Install the JWT library

```bash
go get github.com/golang-jwt/jwt/v5
```
---

## Overview: How JWT Authentication Works
1. User logs in with credentials
2. Server validates credentials
3. Server generates a signed JWT
4. Client stores the token
5. Client sends the token with each request
6. Middleware validates the token before allowing access

JWTs are **stateless**, meaning the server does not store session data.

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
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": req.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

---

## Sending the Token from the Client

```
Authorization: Bearer <JWT_TOKEN>
```

---

## JWT Authentication Middleware

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
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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
claims := r.Context().Value(ClaimsContextKey).(jwt.MapClaims)
username := claims["sub"].(string)
```

---

## Protecting Routes (IMPORTANT)

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

## Common Pitfall

Middleware does **not** run automatically.  
If it’s not attached to the route, it won’t execute.

---

## JWT Secret

Use the same secret everywhere:

```go
var JwtKey = []byte("supersecretkey")
```

Store secrets in environment variables in production.

---

## Security Best Practices

- Use HTTPS
- Set expiration times
- Never store sensitive data in JWT claims
- Keep secrets out of source control

---

## Conclusion

JWT allows you to build secure, stateless APIs in Go.  
Correct middleware attachment is critical for proper authorization.

---

## Next Steps

- Add role-based authorization  
- Implement refresh tokens  
- Add middleware tests  
