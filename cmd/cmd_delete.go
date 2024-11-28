package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/zkMeLabs/mechain-go-sdk/client"
	sdktypes "github.com/zkMeLabs/mechain-go-sdk/types"
)

// cmdDelBucket delete an existed Bucket, the bucket must be empty
func cmdDelBucket() *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Action:    deleteBucket,
		Usage:     "delete an existed bucket",
		ArgsUsage: "BUCKET-URL",
		Description: `
Send a deleteBucket txn to mechain chain, the bucket must be empty before deleting

Examples:
# Delete an existed bucket called mechain-bucket
$ mechain-cmd bucket rm mechain://mechain-bucket/mechain-object`,
	}
}

// cmdDelObject delete an existed object in bucket
func cmdDelObject() *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Action:    deleteObject,
		Usage:     "delete existed object",
		ArgsUsage: "OBJECT-URL",
		Description: `
Send a deleteObject txn to mechain chain

Examples:
# Delete an existed object called mechain-object
$ mechain-cmd object rm mechain://mechain-bucket/mechain-object`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  recursiveFlag,
				Value: false,
				Usage: "performed on all files or objects under the specified directory or prefix in a recursive way",
			},
		},
	}
}

// cmdDelGroup delete an existed group
func cmdDelGroup() *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Action:    deleteGroup,
		Usage:     "delete an existed group",
		ArgsUsage: "GROUP-NAME",
		Description: `
Send a deleteGroup txn to mechain chain

Examples:
# Delete an existed group
$ mechain-cmd group rm group-name`,
	}
}

// deleteBucket send the deleteBucket msg to mechain
func deleteBucket(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return toCmdErr(fmt.Errorf("args number more than one"))
	}
	bucketName, err := getBucketNameByUrl(ctx)
	if err != nil {
		return toCmdErr(err)
	}

	client, privateKey, err := NewClient(ctx, ClientOptions{IsQueryCmd: false})
	if err != nil {
		return toCmdErr(err)
	}

	c, cancelDelBucket := context.WithCancel(globalContext)
	defer cancelDelBucket()

	_, err = client.HeadBucket(c, bucketName)
	if err != nil {
		fmt.Printf("bucket %s not exist or already deleted\n", bucketName)
	}

	txnHash, err := client.DeleteBucket(c, bucketName, sdktypes.DeleteBucketOption{TxOpts: &TxnOptionWithSyncMode}, privateKey)
	if err != nil {
		fmt.Println("delete bucket error:", err.Error())
		return nil
	}

	err = waitTxnStatus(client, c, txnHash, "DeleteBucket")
	if err != nil {
		return toCmdErr(err)
	}

	fmt.Printf("delete_bucket: %s \ntransaction hash: %s\n", bucketName, txnHash)
	return nil
}

// deleteObject send the deleteBucket msg to mechain
func deleteObject(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return toCmdErr(fmt.Errorf("args number more than one"))
	}
	var (
		deleteAll              bool
		bucketName, objectName string
		prefixName             string
		err                    error
		paramErr               error
	)

	urlInfo := ctx.Args().Get(0)
	bucketName, objectName, err = getObjAndBucketNames(urlInfo)
	supportRecursive := ctx.Bool(recursiveFlag)
	if err != nil {
		// if delete all the object in a recursive way, just need to parse bucket name
		if supportRecursive {
			bucketName, paramErr = getBucketNameByUrl(ctx)
			if paramErr != nil {
				return toCmdErr(err)
			} else {
				deleteAll = true
			}
		} else {
			return toCmdErr(err)
		}
	}

	client, privateKey, err := NewClient(ctx, ClientOptions{IsQueryCmd: false})
	if err != nil {
		return toCmdErr(err)
	}

	c, cancelDelObject := context.WithCancel(globalContext)
	defer cancelDelObject()
	if supportRecursive {
		if !deleteAll {
			// if it is a folder and set the --recursive flag , list all the objects and delete them one by one
			prefixName = objectName
			if !strings.HasSuffix(prefixName, "/") {
				prefixName = objectName + "/"
			}
			err = deleteObjectByPage(client, c, bucketName, prefixName, privateKey)
		} else {
			// list all the objects in the bucket and delete them
			err = deleteObjectByPage(client, c, bucketName, prefixName, privateKey)
		}
		if err != nil {
			return toCmdErr(err)
		}

	} else {
		deleteObjectAndWaitTxn(client, c, bucketName, objectName, privateKey)
	}

	return nil
}

func deleteObjectByPage(cli client.IClient, c context.Context, bucketName, prefixName, privateKey string) error {
	var (
		listResult        sdktypes.ListObjectsResult
		continuationToken string
		err               error
	)

	for {
		listResult, err = cli.ListObjects(c, bucketName, sdktypes.ListObjectsOptions{
			ShowRemovedObject: false,
			MaxKeys:           defaultMaxKey,
			ContinuationToken: continuationToken,
			Prefix:            prefixName,
		})
		if err != nil {
			return toCmdErr(err)
		}

		// TODO use one txn to broadcast multi delete object messages
		for _, object := range listResult.Objects {
			// no need to return err if some objects delete failed
			deleteObjectAndWaitTxn(cli, c, bucketName, object.ObjectInfo.ObjectName, privateKey)
		}

		if !listResult.IsTruncated {
			break
		}

		continuationToken = listResult.NextContinuationToken
	}
	return nil
}

func deleteObjectAndWaitTxn(cli client.IClient, c context.Context, bucketName, objectName, privateKey string) {
	txnHash, err := cli.DeleteObject(c, bucketName, objectName, sdktypes.DeleteObjectOption{TxOpts: &TxnOptionWithSyncMode}, privateKey)
	if err != nil {
		fmt.Printf("failed to delete object %s err:%v\n", objectName, err)
		return
	}

	err = waitTxnStatus(cli, c, txnHash, "DeleteObject")
	if err != nil {
		fmt.Printf("failed to query the txn of deleting object %s, err:%v\n", objectName, err)
		return
	}

	fmt.Printf("delete: %s\n", objectName)
}

// deleteGroup send the deleteGroup msg to mechain
func deleteGroup(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return toCmdErr(fmt.Errorf("args number more than one"))
	}
	groupName, err := getGroupNameByUrl(ctx)
	if err != nil {
		return toCmdErr(err)
	}

	client, privateKey, err := NewClient(ctx, ClientOptions{IsQueryCmd: false})
	if err != nil {
		return toCmdErr(err)
	}

	c, cancelDelGroup := context.WithCancel(globalContext)
	defer cancelDelGroup()
	txnHash, err := client.DeleteGroup(c, groupName, sdktypes.DeleteGroupOption{TxOpts: &TxnOptionWithSyncMode}, privateKey)
	if err != nil {
		return toCmdErr(err)
	}

	err = waitTxnStatus(client, c, txnHash, "DeleteGroup")
	if err != nil {
		return toCmdErr(err)
	}

	fmt.Printf("delete_group: %s \ntransaction hash: %s\n", groupName, txnHash)
	return nil
}
