package tls

import (
	"fmt"

	"github.com/octavore/naga/service"
)

func (m *Module) registerCommands(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "tls:provision <example.com> <my@email.com>",
		ShortUsage: `Provision an ssl cert for the given domain and email`,
		Usage: `Provision an ssl cert for the given domain.
Required params: domain to provision a cert for; contact email for Let's Encrypt.`,
		Flags: []*service.Flag{{Key: "agree"}},
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(2)
			if !ctx.Flags["agree"].Present() {
				fmt.Print("Please provide the --agree flag to indicate that you agree to Let's Encrypt's TOS. \n")
				return
			}
			err := m.ObtainCert(ctx.Args[1], ctx.Args[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("success!")
		},
	})

	c.AddCommand(&service.Command{
		Keyword:    "tls:renew",
		ShortUsage: `Renew SSL certs`,
		Usage:      `Renew installed SSL certs.`,
		Run: func(ctx *service.CommandContext) {
			err := m.renewExpiredCerts()
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println("success!")
		},
	})
}
