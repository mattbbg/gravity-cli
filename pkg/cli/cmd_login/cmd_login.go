package cmd_login

import (
	"fmt"

	"github.com/spf13/cobra"
)

var PasswordFlag string
var AccountFlag string

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login gravity user",
	Long:  `Login user by Password and Account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=========Gravity Account=========")
		fmt.Println("Account: " + AccountFlag)
		fmt.Println("Password: " + PasswordFlag)
		fmt.Println("=================================")
	},
}

func NewLoginCmd() *cobra.Command {

	LoginCmd.Flags().StringVarP(&PasswordFlag, "gravity-passowrd", "p", "", "Gravity password")
	LoginCmd.Flags().StringVarP(&AccountFlag, "gravity-account", "a", "", "Gravity account")

	return LoginCmd
}
