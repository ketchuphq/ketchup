package users

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/users/databaseauth"
	"github.com/octavore/nagax/util/token"
)

func registerUserAdd(m *Module) *service.Command {
	return &service.Command{
		Keyword:    "users:add <email>",
		Usage:      "Add a new user.",
		ShortUsage: "Add a new user",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			email := ctx.Args[0]
			fmt.Println("enter a password:")
			pass, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			_, err = m.DBAuth.Create(email, string(pass))
			if err != nil {
				panic(err)
			}
		},
	}
}

func registerSetPassword(m *Module) *service.Command {
	return &service.Command{
		Keyword:    "users:password <email>",
		Usage:      "Set user password.",
		ShortUsage: "Set user password",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			email := ctx.Args[0]
			u, err := m.DB.GetUserByEmail(email)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("enter a password:")
			pass, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			hashedPass := databaseauth.HashPassword(string(pass), token.New32())
			u.SetHashedPassword(&hashedPass)
			err = m.DB.UpdateUser(u)
			if err != nil {
				panic(err)
			}
		},
	}
}
