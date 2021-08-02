package cli

import (
	"fmt"
	"time"

	app "github.com/BrobridgeOrg/gravity-cli/pkg/app"
	accessKeyCmd "github.com/BrobridgeOrg/gravity-cli/pkg/cli/commands/accessKey"
	rootCmd "github.com/BrobridgeOrg/gravity-cli/pkg/cli/commands/root"
	setConfigCmd "github.com/BrobridgeOrg/gravity-cli/pkg/cli/commands/setConfig"
	auth "github.com/BrobridgeOrg/gravity-sdk/authenticator"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/BrobridgeOrg/gravity-sdk/core/keyring"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cli struct {
	app        app.App
	configFile string
}

func NewCli(a app.App) *Cli {

	cli := &Cli{
		app: a,
	}

	return cli

}

func (cli *Cli) Init() error {

	//New Root Command
	rCmd := rootCmd.NewRootCmd(cli)
	cmd := rCmd.Init()

	// New setConfig Command
	setConfigCmd := setConfigCmd.NewSetConfigCmd(cli)
	setConfig := setConfigCmd.Init()

	// New AccesKey Command
	acckeyCmd := accessKeyCmd.NewAccessKeyCmd(cli)
	accessKey := acckeyCmd.Init()

	// Root Combine
	cmd.AddCommand(setConfig)
	cmd.AddCommand(accessKey)

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}

func (cli *Cli) RunCli() error {
	// Nothing to do

	return nil
}

func (cli *Cli) GetAuthClient() (*auth.Authenticator, error) {
	//Read config
	viper.SetConfigType("toml")
	viper.SetConfigFile(cli.configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Reading confing error: ", err)
		log.Error("Please check confing file: ", cli.configFile)
		return nil, err
	}

	// Set viper value
	appID := viper.GetString("gravity.appID")
	host := viper.GetString("gravity.host")
	domain := viper.GetString("gravity.domain")
	accessKey := viper.GetString("gravity.accessKey")

	// connect to gravity server
	options := core.NewOptions()
	options.PingInterval = time.Duration(10) * time.Second
	options.MaxPingsOutstanding = 3
	options.MaxReconnects = -1
	client := core.NewClient()
	if err := client.Connect(host, options); err != nil {
		return nil, err
	}

	// Initializing authenticator and Connect to server
	authOptions := auth.NewOptions()
	authOptions.Domain = domain
	authOptions.Key = keyring.NewKey(appID, accessKey)
	authOptions.Channel = fmt.Sprintf("%s.authentication_manager", domain)

	authClient := auth.NewAuthenticatorWithClient(client, authOptions)
	if err := authClient.Connect(host, options); err != nil {
		return nil, err
	}

	return authClient, nil
}

func (cli *Cli) SetConfigFile(file string) {
	cli.configFile = file
}

func (cli *Cli) GetConfigFile() string {
	return cli.configFile
}
