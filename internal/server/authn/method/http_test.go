package method

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.flipt.io/flipt/internal/config"
	"go.flipt.io/flipt/rpc/flipt/auth"
)

func TestForwardCookies(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/foo", nil)
	req.AddCookie(&http.Cookie{Name: stateCookieKey, Value: "state-value"})
	req.AddCookie(&http.Cookie{Name: tokenCookieKey, Value: "token-value"})

	md := ForwardCookies(t.Context(), req)

	assert.Equal(t, []string{"state-value"}, md.Get(stateCookieKey))
	assert.Equal(t, []string{"token-value"}, md.Get(tokenCookieKey))
}

func TestForwardResponseOptionHostOnlyCookie(t *testing.T) {
	m := NewHTTPMiddleware(config.AuthenticationSessionConfig{TokenLifetime: time.Hour})
	w := httptest.NewRecorder()

	require.NoError(t, m.ForwardResponseOption(t.Context(), w, &auth.CallbackResponse{ClientToken: "token-value"}))

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, tokenCookieKey, cookies[0].Name)
	assert.Equal(t, "token-value", cookies[0].Value)
	assert.NotContains(t, w.Header().Get("Set-Cookie"), "Domain=")
}

func TestForwardPrefix(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected []string
	}{
		{"forward", map[string]string{"X-Forwarded-Prefix": "/my-prefix"}, []string{"/my-prefix"}},
		{"none", map[string]string{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "GET", "http://example.com/foo", nil)
			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}
			md := ForwardPrefix(t.Context(), req)
			assert.Equal(t, tt.expected, md.Get(ForwardedPrefixKey))
		})
	}
}
