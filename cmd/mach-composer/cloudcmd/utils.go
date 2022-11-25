package cloudcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/labd/mach-composer/internal/cloud"
	"github.com/labd/mach-composer/internal/utils"
)

func getClient(cmd *cobra.Command) (*mccsdk.APIClient, context.Context) {
	ctx := cmd.Context()
	cfg := &cloud.ClientConfig{
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
			Scopes:       []string{"foobar"},
			TokenURL:     endpoints.TokenURL,
		}
		cfg.HTTPClient = oauth2Config.Client(
			context.WithValue(context.TODO(), oauth2.HTTPClient, cfg.HTTPClient))
	} else {
		oauth2Config := &oauth2.Config{
			Endpoint: endpoints,
		}

		token := &oauth2.Token{
			AccessToken:  viper.GetString("token.access"),
			RefreshToken: viper.GetString("token.refresh"),
			Expiry:       viper.GetTime("token.expiry"),
		}

		cfg.HTTPClient = oauth2Config.Client(ctx, token)
	}

	client := cloud.NewClient(cfg)
	return client, ctx
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

func handleError(err error) error {
	if openApiErr, ok := err.(*mccsdk.GenericOpenAPIError); ok {
		remoteErr := openApiErr.Model()
		if svcErr, ok := remoteErr.(mccsdk.Error); ok {
			fmt.Printf("error: %s\n", svcErr.GetError())
		} else {
			var errorMsg string

			errorData := struct {
				Message string `json:"message"`
			}{}
			if err := json.Unmarshal(openApiErr.Body(), &errorData); err == nil {
				errorMsg = errorData.Message
			} else {
				errorMsg = openApiErr.Error()
			}
			fmt.Println("Server returned an error:", errorMsg)
		}
	} else {
		fmt.Println("Internal error: ", err)
	}

	os.Exit(1)
	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustGetString(cmd *cobra.Command, key string) string {
	value, err := cmd.Flags().GetString(key)
	if err != nil {
		panic(err)
	}
	return value
}

func writeTable(writer io.Writer, header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render() // Send output
	fmt.Println()
}
