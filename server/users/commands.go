package users

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

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

func registerListUsers(m *Module) *service.Command {
	return &service.Command{
		Keyword:    "users:list",
		Usage:      "List users.",
		ShortUsage: "List users",
		Run: func(ctx *service.CommandContext) {
			lst, err := m.DB.ListUsers()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, user := range lst {
				fmt.Println(user.GetEmail(), user.GetHashedPassword())
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

func registerGenerateToken(m *Module) *service.Command {
	return &service.Command{
		Keyword:    "users:token <email>",
		Usage:      "Generate user API token.",
		ShortUsage: "Generate user API token. Replaces any existing token.",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			email := ctx.Args[0]
			u, err := m.DB.GetUserByEmail(email)
			if err != nil {
				fmt.Println(err)
				return
			}

			b := make([]byte, 15)
			_, err = rand.Read(b)
			if err != nil {
				fmt.Println(err)
				return
			}

			t := strings.ToLower(base64.StdEncoding.EncodeToString(b))
			u.SetToken(&t)
			err = m.DB.UpdateUser(u)
			if err != nil {
				panic(err)
			}
			fmt.Println("generated token:", t)
		},
	}
}
