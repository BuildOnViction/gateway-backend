package jwt

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

type UserClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

func UserClaimFactory() jwt.Claims {
	return &UserClaims{}
}

const expiration = 86400

func GenerateToken(signingKey []byte, userAddress string) (string, error) {
	claims := UserClaims{
		userAddress,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expiration).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// for attaching jwt verification for each Handler inside transport_grpc or transport http
func VerifyToken(keyFunc jwt.Keyfunc, method jwt.SigningMethod, newClaims ClaimsFactory) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// tokenString is stored in the context from the transport handlers.
			tokenString, ok := ctx.Value(JWTTokenContextKey).(string)
			if !ok {
				return nil, errors.WithStack(JWTError{Violates: map[string][]string{
					"access_token": {
						"ACCESS_TOKEN.MISSING",
						ErrTokenContextMissing.Error(),
					},
				}})
			}

			// Parse takes the token string and a function for looking up the
			// key. The latter is especially useful if you use multiple keys
			// for your application.  The standard is to use 'kid' in the head
			// of the token to identify which key to use, but the parsed token
			// (head and claims) is provided to the callback, providing
			// flexibility.
			token, err := jwt.ParseWithClaims(tokenString, newClaims(), func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if token.Method != method {
					return nil, errors.WithStack(JWTError{Violates: map[string][]string{
						"access_token": {
							"ACCESS_TOKEN.UNEXPECTED_SIGNING_METHOD",
							ErrUnexpectedSigningMethod.Error(),
						},
					}})
				}

				return keyFunc(token)
			})
			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						// Token is malformed
						return nil, errors.WithStack(JWTError{Violates: map[string][]string{
							"access_token": {
								"ACCESS_TOKEN.MALFORMED",
								ErrTokenMalformed.Error(),
							},
						}})
					case e.Errors&jwt.ValidationErrorExpired != 0:
						// Token is expired
						return nil, errors.WithStack(JWTError{Violates: map[string][]string{
							"access_token": {
								"ACCESS_TOKEN.EXPIRED",
								ErrTokenExpired.Error(),
							},
						}})
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						// Token is not active yet
						return nil, errors.WithStack(JWTError{Violates: map[string][]string{
							"access_token": {
								"ACCESS_TOKEN.INACTIVE",
								ErrTokenNotActive.Error(),
							},
						}})
					case e.Inner != nil:
						// report e.Inner
						return nil, e.Inner
					}
					// We have a ValidationError but have no specific Go kit error for it.
					// Fall through to return original error.
				}
				return nil, err
			}

			if !token.Valid {
				return nil, errors.WithStack(JWTError{Violates: map[string][]string{
					"access_token": {
						"ACCESS_TOKEN.INVALID",
						ErrTokenInvalid.Error(),
					},
				}})
			}

			claims, ok := token.Claims.(*UserClaims)
			if !ok {
				return nil, errors.WithStack(JWTError{Violates: map[string][]string{
					"access_token": {
						"ACCESS_TOKEN.INVALID",
						ErrTokenInvalid.Error(),
					},
				}})
			}

			ctx = context.WithValue(ctx, "User", claims.Address)
			ctx = context.WithValue(ctx, JWTClaimsContextKey, token.Claims)

			return next(ctx, request)
		}
	}
}

// Verify verifies the access token string and return a user claim if the token is valid
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
