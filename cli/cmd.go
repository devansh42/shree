package main

// func initCommands() cli.App{
// 	return nil
// }

//
/*
func initCommands() cli.App {

	initApp()

	app := cli.App{Name: "shree",
		Commands: []*cli.Command{
			&cli.Command{
				Name: "get",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name: "me",
						Subcommands: []*cli.Command{
							&cli.Command{
								Name:   "in",
								Action: signIn},

							&cli.Command{
								Name:   "out",
								Action: signOut}}}},

				&cli.Command{
					Name: "connect",
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:    "src",
							Aliases: []string{"s"},
						},
						cli.StringFlag{
							Name:    "protocol",
							Aliases: "p"},

						cli.UintFlag{
							Name:    "dest",
							Aliases: []string{"d"}}},

					Action: connectLocalTunnel},

				&cli.Command{
					Name: "disconnect",
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:    "port",
							Aliases: []string{"p"}}},

					Action: disconnectLocalTunnel}}}}

	return app
}
*/
