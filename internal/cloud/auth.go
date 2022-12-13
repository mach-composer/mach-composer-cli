package cloud

// This module implements the PKCE for OAuth2 flow. It starts a webserver and
// open the user's browser to open a authorization URL using the local http
// endpoint as the redirect URL.

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/grokify/go-pkce"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/authhandler"
)

var (
	LocalPortStart = 49852
	LocalPortEnd   = 49862
)

var ErrAuthNotCompleted = errors.New("Authentication not completed")
var ErrAuthTimeout = errors.New("Authentication aborted (timeout)")

type Options struct {
	Timeout           int
	CompletedRedirect string
	CompletedTemplate string
}

// AuthorizeUser implements the PKCE OAuth2 flow.
func StartAuthentication(ctx context.Context, config oauth2.Config, options *Options) (*oauth2.Token, error) {
	setOptionDefaults(options)

	state := generateStateParam()
	pkceParams := generatePKCEParams()
	handler, redirectURL, err := createAuthHandler(ctx, options)
	if err != nil {
		return nil, err
	}

	config.RedirectURL = redirectURL
	ts := authhandler.TokenSourceWithPKCE(ctx, &config, state, handler, &pkceParams)
	return ts.Token()
}

func createAuthHandler(ctx context.Context, options *Options) (authhandler.AuthorizationHandler, string, error) {
	listener, err := getAvailableListener()
	if err != nil {
		return nil, "", err
	}

	handler := func(authCodeURL string) (string, string, error) {
		var wg sync.WaitGroup
		wg.Add(1)

		authCode := ""
		state := ""

		errs := make(chan error, 1)

		m := http.NewServeMux()
		server := &http.Server{Addr: listener.Addr().String(), Handler: m}
		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			authCode = r.URL.Query().Get("code")
			state = r.URL.Query().Get("state")

			if options.CompletedTemplate != "" {
				w.Header().Add("Content-Type", "text/html")
				if _, err := io.WriteString(w, options.CompletedTemplate); err != nil {
					errs <- err
				}
			} else if options.CompletedRedirect != "" {
				http.Redirect(
					w, r, options.CompletedRedirect, http.StatusTemporaryRedirect)
			}
			wg.Done()
		})

		// Start server
		go func() {
			if waitTimeout(&wg, time.Duration(options.Timeout)*time.Second) {
				errs <- ErrAuthTimeout
				close(errs)
			}
			if err := server.Shutdown(ctx); err != nil {
				errs <- err
				close(errs)
			}
		}()

		// Open user's browser
		c := color.New(color.FgYellow)
		c.Printf("A web browser has been opened at %s. Please continue the login via the web browser.\n", authCodeURL)
		err := open.Start(authCodeURL)
		if err != nil {
			return "", "", fmt.Errorf("unable to open browser")
		}

		// Wait till we the authentication process is either completed or the
		// timeout was triggered.
		if err := server.Serve(listener); err != http.ErrServerClosed {
			panic(err)
		}

		select {
		case err := <-errs:
			return "", "", err
		default:
		}

		if authCode == "" || state == "" {
			return authCode, state, ErrAuthNotCompleted
		}

		return authCode, state, nil
	}

	redirectURI := fmt.Sprintf("http://%s", listener.Addr().String())
	return handler, redirectURI, nil
}

// getAvailableListener returns a working listener by trying all ports between
// the global LocalPortStart and LocalPortEnd range.
func getAvailableListener() (net.Listener, error) {
	for port := LocalPortStart; port <= LocalPortEnd; port++ {

		// Listen on 127.0.0.1, using `localhost` here can result into issues
		// on windows.
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		if listener, err := net.Listen("tcp", addr); err == nil {
			return listener, nil
		}
	}
	return nil, fmt.Errorf(
		"no free TCP port available in range %d-%d", LocalPortStart, LocalPortEnd)
}

func generatePKCEParams() authhandler.PKCEParams {
	pkceParams := authhandler.PKCEParams{
		ChallengeMethod: "S256",
		Verifier:        pkce.NewCodeVerifier(),
	}
	pkceParams.Challenge = pkce.CodeChallengeS256(pkceParams.Verifier)
	return pkceParams
}

// If a state query param is not passed in, generate a random
// base64-encoded nonce so that the state on the auth URL
// is unguessable, preventing CSRF attacks, as described in
//
// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
func generateStateParam() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		panic("source of randomness unavailable: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func setOptionDefaults(options *Options) {
	if options.Timeout == 0 {
		options.Timeout = 120
	}

	if options.CompletedRedirect == "" && options.CompletedTemplate == "" {
		options.CompletedTemplate = strings.TrimSpace(`
			<html>
				<body>
					<h1>Login successful!</h1>
					<h2>You can close this window</h2>
				</body>
			</html>
			`)
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
// Taken from https://stackoverflow.com/a/32843750
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
