// Copyright (c) 2021 Wireleap

package initcmd

import (
	"flag"
	"log"
	"text/tabwriter"

	"github.com/wireleap/client/clientcfg"
	"github.com/wireleap/client/filenames"
	"github.com/wireleap/client/sub/configcmd"
	"github.com/wireleap/client/sub/initcmd/embedded"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd() *cli.Subcmd {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	force := fs.Bool("force-unpack-only", false, "Overwrite embedded files only")
	r := &cli.Subcmd{
		FlagSet: fs,
		Desc:    "Initialize wireleap directory",
		Run: func(fm fsdir.T) {
			if err := cli.UnpackEmbedded(embedded.FS, fm, *force); err != nil {
				log.Fatalf("error while unpacking embedded files: %s", err)
			}
			if !*force {
				if err := fm.Set(clientcfg.Defaults(), filenames.Config); err != nil {
					log.Fatalf("could not write initial config.json: %s", err)
				}
				for k, v := range map[string]string{
					"address.socks":           "127.0.0.1:13491",
					"address.tun":             "10.13.49.0:13493",
					"address.h2c":             "127.0.0.1:13492",
					"accesskey.use_on_demand": "true",
				} {
					configcmd.Run(fm, k, v)
				}
			}
		},
	}
	r.Writer = tabwriter.NewWriter(r.FlagSet.Output(), 0, 8, 6, ' ', 0)
	return r
}
