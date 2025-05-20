package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

/*
	"Claims" in JWT = Data embedded in the token

	What Are "Registered Claims"?
	JWT has a set of registered (standard) fields defined in the spec (RFC 7519) that help describe the token's purpose, origin, and validity.
*/

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   string(userID.String()),
	})

	tokenStr, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

/*
	func(token *jwt.Token) (interface{}, error) {
	    return []byte(tokenSecret), nil
	}
	-- This is a function that takes the parsed token and returns the key needed to verify the signature.

    Why a function?
    Because the key may depend on the token header (e.g., different algorithms or key IDs). The JWT lib calls this function after reading the token header but before verifying the signature.

    What it returns:
    It returns the key as an interface{}, typically:
        A []byte for HMAC symmetric signing (like your tokenSecret).
        A public key for asymmetric signing (RSA, ECDSA).

    In your code:
    You return the tokenSecret converted to a byte slice (because HMAC signing uses a secret key as bytes).
*/

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	/*
		we are passing `&claims` ptr to claims cuz both
		jwt.RegisteredClaims{} and it's ptr implements interface jwt.Claims
		and ptr to claims allows for its modification
	*/
	tkn, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userId, err := tkn.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("missing authorization header")
	}

	strippedToken, found := strings.CutPrefix(token, "Bearer ")
	if !found {
		return "", errors.New("misformed authorization header")
	}

	return strippedToken, nil
}
