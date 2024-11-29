package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	sdktypes "github.com/zkMeLabs/mechain-go-sdk/types"
)

// cmdBuyQuota buy the read quota of the bucket
func cmdBuyQuota() *cli.Command {
	return &cli.Command{
		Name:      "buy-quota",
		Action:    buyQuotaForBucket,
		Usage:     "update bucket quota info",
		ArgsUsage: "BUCKET-URL",
		Description: `
Update the read quota metadata of the bucket, indicating the target quota of the bucket.
The command need to set the target quota with --chargedQuota 

Examples:
$ mechain-cmd payment buy-quota  --chargedQuota 1000000  mechain://bucket-name`,
		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Name:     chargeQuotaFlag,
				Usage:    "indicate the target quota to be set for the bucket",
				Required: true,
			},
		},
	}
}

func cmdGetQuotaInfo() *cli.Command {
	return &cli.Command{
		Name:      "get-quota",
		Action:    getQuotaInfo,
		Usage:     "get quota info of the bucket",
		ArgsUsage: "BUCKET-URL",
		Description: `
Get charged quota, free quota and consumed quota info from storage provider 

Examples:
$ mechain -c config.toml payment quota-info  mechain://bucket-name`,
	}
}

// buyQuotaForBucket set the charged quota meta of bucket on chain
func buyQuotaForBucket(ctx *cli.Context) error {
	bucketName, err := getBucketNameByUrl(ctx)
	if err != nil {
		return toCmdErr(err)
	}

	client, privateKey, err := NewClient(ctx, ClientOptions{IsQueryCmd: false})
	if err != nil {
		return toCmdErr(err)
	}

	c, cancelBuyQuota := context.WithCancel(globalContext)
	defer cancelBuyQuota()

	// if bucket not exist, no need to buy quota
	_, err = client.HeadBucket(c, bucketName)
	if err != nil {
		return toCmdErr(ErrBucketNotExist)
	}

	targetQuota := ctx.Uint64(chargeQuotaFlag)
	if targetQuota == 0 {
		return toCmdErr(errors.New("target quota not set"))
	}

	txnHash, err := client.BuyQuotaForBucket(c, bucketName, targetQuota, sdktypes.BuyQuotaOption{TxOpts: &TxnOptionWithSyncMode}, privateKey)
	if err != nil {
		fmt.Println("buy quota error:", err.Error())
		return nil
	}

	fmt.Printf("buy quota for bucket: %s \n", bucketName)
	fmt.Println("transaction hash: ", txnHash)
	return nil
}

// getQuotaInfo query the quota price info of sp from mechain chain
func getQuotaInfo(ctx *cli.Context) error {
	bucketName, err := getBucketNameByUrl(ctx)
	if err != nil {
		return toCmdErr(err)
	}

	client, _, err := NewClient(ctx, ClientOptions{IsQueryCmd: false})
	if err != nil {
		return toCmdErr(err)
	}

	c, cancelGetQuota := context.WithCancel(globalContext)
	defer cancelGetQuota()

	// if bucket not exist, no need to get info of quota
	_, err = client.HeadBucket(c, bucketName)
	if err != nil {
		return toCmdErr(ErrBucketNotExist)
	}

	quotaInfo, err := client.GetBucketReadQuota(c, bucketName)
	if err != nil {
		return toCmdErr(err)
	}

	nameMaxLen := len("consumed charged quota:")
	format := fmt.Sprintf("%%-%ds %%-%dd   \n", nameMaxLen, 50)
	firstLineFormat := fmt.Sprintf("%%-%ds %%-%ds  \n", nameMaxLen, 50)
	fmt.Printf(firstLineFormat, "quota name", "quota value")
	fmt.Printf(format, "charged quota:", quotaInfo.ReadQuotaSize)
	fmt.Printf(format, "remained free quota:", quotaInfo.SPFreeReadQuotaSize)
	fmt.Printf(format, "consumed charged quota:", quotaInfo.ReadConsumedSize)
	fmt.Printf(format, "consumed free quota:", quotaInfo.FreeConsumedSize)

	return nil
}
