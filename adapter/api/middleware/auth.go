package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Vractos/dolly/pkg/contexttools"
	"github.com/Vractos/dolly/pkg/metrics"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"go.uber.org/zap"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	InternalID string `json:"internal_id,omitempty"`
	Scope      string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	if c.InternalID == "" {
		return errors.New("store id not provided")
	}
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken(logger metrics.Logger) func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		logger.Fatal("Failed to parse the issuer url", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		logger.Fatal("Failed to set up the jwt validator", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Warn("Encountered error while validating JWT", zap.NamedError("reason", err))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

// Extract StoreID (internal_id) from CustomClaims and add to request context.
//
// Must only be used after the EnsureValidToken method on the middleware chain
func AddStoreIDToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		storeId := claims.CustomClaims.(*CustomClaims).InternalID
		r = r.Clone(context.WithValue(r.Context(), contexttools.ContextKeyStoreId, storeId))
		next.ServeHTTP(w, r)
	})
}
