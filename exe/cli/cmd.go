package main

import (
	"strings"
	"time"

	cli "github.com/urfave/cli/v2"
)

func welComeMsg() {
	print(COLOR_GREEN)

	println("\t", strings.Repeat("-=", 50))
	println("\tHello, I am Shree, your partner in tunneling!!")
	println("\tI can make local and remote tunnels for you")
	println("\tType \"help\" for listing available options")
	println("\tYou can change backend and ssh server address with 'set' command, type 'help set' for more info")
	println("\tHappy Tunneling!!")
	println("\t", strings.Repeat("=-", 50))
	resetConsoleColor()
}

func getCliApp() cli.App {

	app := cli.App{Name: "shree",
		Version: "0.0.1",
		CommandNotFound: func(c *cli.Context, cmd string) {
			println("Command ", cmd, " not found")
		},
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{Name: "devansh42", Email: "devanshguptamrt@gmail.com"},
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
				Name:    "whoami",
				Aliases: []string{"who"},
				Usage:   "Reveals current user",
				Action:  whoAmI,
			},
			&cli.Command{
				Name:  "set",
				Usage: "Sets cli properties like remote server address",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "backend",
						Aliases: []string{"b"},
						Usage:   "sets the backend server address to the given value",
					},
					&cli.StringFlag{
						Name:    "remote",
						Aliases: []string{"r"},
						Usage:   "sets the remote server address to the given value",
					},
				},
				Action: setProps,
			},
			&cli.Command{
				Name:   "get",
				Usage:  "Gets any specific property e.g get --backend",
				Action: getProps,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "backend",
						Usage:   "gets backend address",
						Aliases: []string{"b"},
					},
					&cli.BoolFlag{
						Name:    "remote",
						Usage:   "gets remote address",
						Aliases: []string{"r"},
					},
				},
			},
			&cli.Command{
				Name:  "sign",
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
						Aliases:  []string{"s"},
					},

					&cli.UintFlag{
						Name:     "dest",
						Usage:    "specifies port to which tunnel points to ",
						Required: true,
						Aliases:  []string{"d"},
					}},

				Action: connectLocalTunnel},
			&cli.Command{ //Remote port forwarding
				Name:  "expose",
				Usage: "Makes a remote tunnel from src (on remote machine) -> dest (on local machine)",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:     "dest",
						Required: true,
						Aliases:  []string{"d"},
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
