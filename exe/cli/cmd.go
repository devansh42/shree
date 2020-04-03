package main

import (
	"time"

	cli "github.com/urfave/cli/v2"
)

func welComeMsg() {
	println("Hello, I am Shree, your partner in tunneling!!\nI can make local and remote tunnels for you")
	println("Type \"help\" for listing available options")
	println("Happy Tunneling!!")

}

func getCliApp() cli.App {

	app := cli.App{Name: "shree",
		Version: "0.0.1",

		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{Name: "devansh42"},
		},
		Description: "Cli for Shree - A Tunneling Solution for developers",
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "exit",
				Aliases: []string{"quit"},

				Action: exitApp,
				Usage:  "Exit gracefully",
			},

			&cli.Command{
				Name:  "get",
				Usage: "Command for login/logout",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name: "me",
						Subcommands: []*cli.Command{
							&cli.Command{
								Name:   "in",
								Usage:  "Authenticates yourself",
								Action: signIn},

							&cli.Command{
								Name:   "out",
								Usage:  "Sign out yourself from current session",
								Action: signOut}}}}},
			&cli.Command{ //Local port forwarding
				Name:  "connect",
				Usage: "Makes a local tunnel from src -> dest ",
				Flags: []cli.Flag{

					&cli.UintFlag{
						Name:     "src",
						Usage:    "specifies port to be open for tunneling",
						Required: true,
					},

					&cli.UintFlag{
						Name:     "dest",
						Usage:    "specifies port to which tunnel points to ",
						Required: true,
					}},

				Action: connectLocalTunnel},
			&cli.Command{ //Remote port forwarding
				Name:  "expose",
				Usage: "Makes a remote tunnel from src (on remote machine) -> dest (on local machine)",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:     "dest",
						Required: true,
						Usage:    "specifies port on local machine to tunnel",
					},
				},
				Action: exposeRemoteTunnel,
			},
			&cli.Command{
				Name:  "list",
				Usage: "Lists local/remote tunnels",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:   "connections", //List local port forwardings
						Action: listLocalTunnels,
						Usage:  "lists local tunnels",
					},
					&cli.Command{
						Name:   "exposure",
						Action: listRemoteTunnels,
						Usage:  "lists remote tunnels",
					},
				},
			},
			&cli.Command{
				Name:  "unexpose",
				Usage: "turn off the remote tunnel",

				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:  "port",
						Usage: "port on local machine to turn off remote tunnel",
					},
				},
				Action: disconnectRemoteTunnel,
			},
			&cli.Command{
				Name:  "disconnect",
				Usage: "turn off the local tunnel",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:  "port",
						Usage: "port on local machine to turn off local tunnel",
					},
				},

				Action: disconnectLocalTunnel}}}

	return app
}
