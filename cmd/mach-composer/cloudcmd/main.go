package cloudcmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/labd/mach-composer/internal/cloud"
)

var CloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Manage your Mach Composer Cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var cloudLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to mach composer cloud",
	RunE: func(cmd *cobra.Command, args []string) error {

		authConfig := oauth2.Config{
			Scopes:   []string{"openid"},
			ClientID: "b0b9ccbd-0613-4ccf-86a1-dab07b8b5619",
			Endpoint: getAuthEndpoint(),
		}

		successURL, err := url.JoinPath(viper.GetString("auth-url"), "/success")
		if err != nil {
			cmd.Println("Invalid authentication url set")
			os.Exit(1)
		}

		token, err := cloud.StartAuthentication(cmd.Context(), authConfig, &cloud.Options{
			CompletedRedirect: successURL,
		})
		if err != nil {
			if errors.Is(err, cloud.ErrAuthNotCompleted) {
				cmd.Println("Authentication not completed")
				os.Exit(1)
			}
			if errors.Is(err, cloud.ErrAuthTimeout) {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			cmd.Println("Error occured during authentication:", err)
			os.Exit(1)
		}

		if token.AccessToken != "" {
			viper.Set("token.access", token.AccessToken)
			viper.Set("token.refresh", token.RefreshToken)
			viper.Set("token.expiry", token.Expiry)
			if err := viper.WriteConfig(); err != nil {
				cmd.PrintErrln(err)
			}
			cmd.Println("Successfully authenticated to mach composer cloud")
		}
		return nil
	},
}

var cloudConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure mach composer cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		hasValue := false

		if cmd.Flags().Changed("set-organization") {
			viper.Set("organization", MustGetString(cmd, "set-organization"))
			hasValue = true
		}

		if cmd.Flags().Changed("set-project") {
			viper.Set("project", MustGetString(cmd, "set-project"))
			hasValue = true
		}

		if cmd.Flags().Changed("set-auth-url") {
			viper.Set("auth-url", MustGetString(cmd, "set-auth-url"))
			hasValue = true
		}

		if cmd.Flags().Changed("set-api-url") {
			viper.Set("api-url", MustGetString(cmd, "set-api-url"))
			hasValue = true
		}

		if hasValue {
			if err := viper.WriteConfig(); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		fmt.Println("Auth URL     : ", viper.GetString("auth-url"))
		fmt.Println("API URL      : ", viper.GetString("api-url"))
		fmt.Println("Organization : ", viper.GetString("organization"))
		fmt.Println("Project      : ", viper.GetString("project"))
		return nil
	},
}

func init() {
	// Config command
	CloudCmd.AddCommand(cloudConfigCmd)
	cloudConfigCmd.Flags().String("set-organization", "", "Set default organization")
	cloudConfigCmd.Flags().String("set-project", "", "Set default project")
	cloudConfigCmd.Flags().String("set-api-url", "", "Set api url")
	cloudConfigCmd.Flags().String("set-auth-url", "https://auth.mach.cloud", "Authentication URL")

	// Login
	CloudCmd.AddCommand(cloudLoginCmd)

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configPath := path.Join(xdg.ConfigHome, "mach-composer")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			fmt.Printf("Error encountered while creating configuration file: %s", err)
			os.Exit(1)
		}
	}

	viper.SetConfigName("mach-composer")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("MCC")

	err := viper.SafeWriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			fmt.Printf("Error encountered while writing configuration file: %s", err)
			os.Exit(1)
		}
	}

	viper.SetDefault("api-url", "https://api.mach.cloud")
	viper.SetDefault("auth-url", "https://auth.mach.cloud")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Invalid config file found at: %s\n", viper.GetViper().ConfigFileUsed())
		os.Exit(1)
	}

	// Copy the values from Viper to the matching flag values.
	// TODO: make recursive
	for _, cmd := range CloudCmd.Commands() {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if cmd.Flags().Changed(f.Name) {
				return
			}

			if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
				Must(cmd.Flags().Set(f.Name, viper.GetString(f.Name)))
			}
		})
	}
}
