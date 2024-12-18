package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var globalContext, _ = context.WithCancel(context.Background())

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, err = os.Getwd()
		if err != nil {
			fmt.Println("fail to get home dir or local dir")
		}
	}

	flags := []cli.Flag{
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:  "host",
				Usage: "host name of request",
			},
		),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:  "rpcAddr",
				Usage: "mechain chain client rpc address",
			},
		),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:  "chainId",
				Usage: "mechain chainId",
			},
		),

		&cli.StringFlag{
			Name:    passwordFileFlag,
			Aliases: []string{"p"},
			Usage:   "password file for encrypting and decoding the private key",
		},
		&cli.StringFlag{
			Name:    configFlag,
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringFlag{
			Name:    keyStoreFlag,
			Aliases: []string{"k"},
			Usage:   "keystore file path",
		},
		&cli.StringFlag{
			Name:  homeFlag,
			Usage: "directory for config and keystore",
			Value: filepath.Join(homeDir, DefaultConfigDir),
		},
	}

	app := &cli.App{
		Name:  "mechain-cmd",
		Usage: "cmd tool for supporting making request to mechain",
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:  "bucket",
				Usage: "support the bucket operation functions, including create/update/delete/head/list and so on",
				Subcommands: []*cli.Command{
					cmdCreateBucket(),
					cmdUpdateBucket(),
					cmdMigrateBucket(),
					cmdDelBucket(),
					cmdHeadBucket(),
					cmdListBuckets(),
					cmdBuyQuota(),
					cmdGetQuotaInfo(),
					cmdMirrorBucket(),
					cmdSetTagForBucket(),
				},
			},
			{
				Name:  "object",
				Usage: "support the object operation functions, including put/get/update/delete/head/list and so on",
				Subcommands: []*cli.Command{
					cmdPutObj(),
					cmdGetObj(),
					cmdDelObject(),
					cmdHeadObj(),
					cmdCancelObjects(),
					cmdListObjects(),
					cmdCalHash(),
					cmdUpdateObject(),
					cmdGetUploadProgress(),
					cmdMirrorObject(),
					cmdSetTagForObject(),
				},
			},
			{
				Name:  "group",
				Usage: "support the group operation functions, including create/update/delete/head/head-member/mirror/renew",
				Subcommands: []*cli.Command{
					cmdCreateGroup(),
					cmdUpdateGroup(),
					cmdHeadGroup(),
					cmdHeadGroupMember(),
					cmdDelGroup(),
					cmdMirrorGroup(),
					cmdRenewGroup(),
					cmdListGroup(),
					cmdListGroupMember(),
					cmdListGroupBelong(),
					cmdSetTagForGroup(),
				},
			},

			{
				Name:  "bank",
				Usage: "support the bank functions, including transfer in mechain and query balance",
				Subcommands: []*cli.Command{
					cmdTransfer(),
					cmdGetAccountBalance(),
					cmdBridge(),
				},
			},
			{
				Name:  "policy",
				Usage: "support object,bucket and group policy operation functions",
				Subcommands: []*cli.Command{
					cmdPutPolicy(),
					cmdDelPolicy(),
					cmdListPolicy(),
				},
			},

			{
				Name:  "payment-account",
				Usage: "support the payment account operation functions",
				Subcommands: []*cli.Command{
					cmdCreatePaymentAccount(),
					cmdPaymentDeposit(),
					cmdPaymentWithdraw(),
					cmdListPaymentAccounts(),
				},
			},
			{
				Name:  "sp",
				Usage: "support the storage provider operation functions",
				Subcommands: []*cli.Command{
					cmdListSP(),
					cmdGetSP(),
					cmdGetQuotaPrice(),
				},
			},
			{
				Name:  "account",
				Usage: "support the keystore operation functions",
				Subcommands: []*cli.Command{
					cmdImportAccount(),
					cmdListAccount(),
					cmdCreateAccount(),
					cmdExportAccount(),
					cmdSetDefaultAccount(),
				},
			},
			{
				Name:  "fee",
				Usage: "support fee grant operation functions",
				Subcommands: []*cli.Command{
					cmdGrantAllowance(),
				},
			},
			{
				Name:  "task",
				Usage: "support the batch upload file",
				Subcommands: []*cli.Command{
					cmdTaskStatus(),
					cmdTaskDelete(),
					cmdTaskRetry(),
				},
			},
			cmdShowVersion(),
		},
	}
	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("config"))

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
