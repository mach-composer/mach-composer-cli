package cloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/labd/mach-composer/internal/utils"
)

const cliClientID = "b0b9ccbd-0613-4ccf-86a1-dab07b8b5619"

type ClientConfig struct {
	URL        string
	HTTPClient *http.Client
}

func NewClient(ctx context.Context) (*mccsdk.APIClient, error) {
	cfg := &ClientConfig{
		URL: viper.GetString("api-url"),
		HTTPClient: &http.Client{
			Transport: utils.DebugTransport,
		},
	}

	clientID := os.Getenv("MCC_CLIENT_ID")
	clientSecret := os.Getenv("MCC_CLIENT_SECRET")
	endpoints := getAuthEndpoint()

	if clientID != "" && clientSecret != "" {
		oauth2Config := &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{},
			TokenURL:     endpoints.TokenURL,
		}
		cfg.HTTPClient = oauth2Config.Client(
			context.WithValue(ctx, oauth2.HTTPClient, cfg.HTTPClient))
	} else {
		oauth2Config := &oauth2.Config{
			ClientID: cliClientID,
			Endpoint: endpoints,
		}

		token := &oauth2.Token{
			AccessToken:  viper.GetString("token.access"),
			RefreshToken: viper.GetString("token.refresh"),
			Expiry:       viper.GetTime("token.expiry"),
		}

		tokenSource := oauth2Config.TokenSource(ctx, token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, err

		}
		if newToken.AccessToken != token.AccessToken {
			viper.Set("token.access", newToken.AccessToken)
			viper.Set("token.refresh", newToken.RefreshToken)
			viper.Set("token.expiry", newToken.Expiry)
			if err := viper.WriteConfig(); err != nil {
				return nil, err
			}
		}

		ctx := context.WithValue(ctx, oauth2.HTTPClient, cfg.HTTPClient)
		cfg.HTTPClient = oauth2.NewClient(ctx, tokenSource)
	}

	configuration := mccsdk.NewConfiguration()
	configuration.Servers = mccsdk.ServerConfigurations{
		{
			URL:         cfg.URL,
			Description: "Mach Composer Cloud Server",
		},
	}
	configuration.HTTPClient = cfg.HTTPClient

	result := mccsdk.NewAPIClient(configuration)
	return result, nil
}

func Login(ctx context.Context) error {
	authConfig := oauth2.Config{
		Scopes:   []string{"openid"},
		ClientID: cliClientID,
		Endpoint: getAuthEndpoint(),
	}

	successURL, err := url.JoinPath(viper.GetString("auth-url"), "/success")
	if err != nil {
		return fmt.Errorf("Invalid authentication url set")
	}

	token, err := StartAuthentication(ctx, authConfig, &Options{
		CompletedRedirect: successURL,
	})
	if err != nil {
		if errors.Is(err, ErrAuthNotCompleted) {
			return fmt.Errorf("Authentication not completed")
		}
		if errors.Is(err, ErrAuthTimeout) {
			return err
		}
		return fmt.Errorf("Error occured during authentication: %w", err)
	}

	if token.AccessToken != "" {
		viper.Set("token.access", token.AccessToken)
		viper.Set("token.refresh", token.RefreshToken)
		viper.Set("token.expiry", token.Expiry)
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil

}

func getAuthEndpoint() oauth2.Endpoint {
	result := oauth2.Endpoint{}
	baseURL := viper.GetString("auth-url")

	if url, err := url.JoinPath(baseURL, "/authorize"); err == nil {
		result.AuthURL = url
	}
	if url, err := url.JoinPath(baseURL, "/oauth/token"); err == nil {
		result.TokenURL = url
	}
	return result
}
