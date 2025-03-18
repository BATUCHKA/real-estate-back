package main

import (
	"fmt"
	"net/mail"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	db := database.Database

	var SuperUserEmail string
	var SuperUserPassword string

	var cmdCreateSuperUser = &cobra.Command{
		Use:   "create-superuser",
		Short: "Create superuser for esan",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, _args []string) {
			if len(SuperUserEmail) > 0 && len(SuperUserPassword) > 0 {
				if _, err := mail.ParseAddress(SuperUserEmail); err != nil {
					fmt.Println("email address invalid")
					return
				}
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(SuperUserPassword), bcrypt.DefaultCost)
				if err != nil {
					panic(err)
				}
				superuser := &models.Users{
					Email:    SuperUserEmail,
					Password: string(hashedPassword),
				}
				result := db.GormDB.Create(&superuser)
				if result.Error != nil {
					panic(result.Error)
				}
				fmt.Println("Succesfully created superuser")
			} else {
				fmt.Println("you must fill email or password (--email, --password)")
			}
		},
	}
	cmdCreateSuperUser.Flags().StringVarP(&SuperUserEmail, "email", "e", "", "superuser email address")
	cmdCreateSuperUser.Flags().StringVarP(&SuperUserPassword, "password", "p", "", "superuser password")

	var cmdPreparePermissions = &cobra.Command{
		Use:   "prepare-permission",
		Short: "insert preloaded permissions to database",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, _args []string) {
			util.Permissions.Flush()
			util.Permissions.FlushPreRoles()
		},
	}

	var cmdWarmUp = &cobra.Command{
		Use:   "warm-up",
		Short: "warm up for static reference tables",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			models.RoleFlush()
		},
	}

	var cmdEnableDBExtensions = &cobra.Command{
		Use:   "enable-extensions",
		Short: "enable postgres required extensions",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, _args []string) {
			db.GormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdCreateSuperUser)
	rootCmd.AddCommand(cmdEnableDBExtensions)
	rootCmd.AddCommand(cmdPreparePermissions)
	rootCmd.AddCommand(cmdWarmUp)
	rootCmd.Execute()
}
