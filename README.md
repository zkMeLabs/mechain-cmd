# mechain-cmd

---
mechain client cmd tool, supporting commands to make requests to mechain

## Disclaimer

**The software and related documentation are under active development, all subject to potential future change without
notification and not ready for production use. The code and security audit have not been fully completed and not ready
for any bug bounty. We advise you to be careful and experiment on the network at your own risk. Stay safe out there.**

## Cmd usage

mechain is still undergoing rapid development iterations, and mechain-cmd also needs to be continuously updated and adapted. When using it, please do not directly use the master branch or develop branch. If you are using this tool on the mechain Mainnet, please switch to the latest official release version.

### installation

```
git clone https://github.com/zkMeLabs/mechain-cmd.git
cd mechain-cmd
make build
cd build
./mechain-cmd -h
```

### basic config

The command tool supports the "--home" option to specify the path of the config file and the keystore, the default path is a directory called ".mechain-cmd" under the home directory of the system.
When running commands that interact with the mechain, if there is no config/config.toml file under the path and the commands runs without "--config" flag,
the tool will generate the config/config.toml file automatically which is consistent with the testnet configuration under the path.

Below is an example of the config file. The rpcAddr and chainId should be consistent with the mechain network.
For mechain Mainnet, you can refer to [mechain Mainnet RPC Endpoints](https://todo).
The rpcAddr indicates the Tendermint RPC address with the port info.

```env
rpcAddr = "https://todo"
chainId = "mechain_5151-1"
```

The command tool can also support other networks besides the Mainnet.
you can replace the content of a custom config file in the default config directory with config.toml or
run command with "-c filepath" to set the custom config file.

#### Get help

The commands support different kinds of commands, including bucket,object,group,bank,policy,sp,payment-account and account.

```
// get help for supporing commands and basic command format
mechain-cmd -h
   bucket           support the bucket operation functions, including create/update/delete/head/list and so on
   object           support the object operation functions, including put/get/update/delete/head/list and so on
   group            support the group operation functions, including create/update/delete/head/head-member/mirror
   bank             support the bank functions, including transfer in mechain and query balance
   policy           support object,bucket and group policy operation functions
   payment-account  support the payment account operation functions
   sp               support the storage provider operation functions
   account          support the keystore operation functions
   version          print version info

```

The following command can be used to obtain help information for commands. For example, you can use "mechain-cmd object -h" to obtain the subcommand infos under the object command.

```
mechain-cmd [command-name] -h
```

The following command can be used to obtain help information for subcommands. For example, you can use "mechain-cmd object update -h" to obtain the help info to update object.

```
mechain-cmd [command-name][subcommand-name] -h
```

### Precautions

1. The user need to use "account import" or "account new" command to init the keystore before running other commands. The "import" command imports account info from private key file and generate a keystore to manage user's private key and the "new" command create a new account with the keystore.The content of the keystore is the encrypted private key information.
All the other commands need to run with -k if the keystore is not on the default path.

2. The operator account should have enough balance before sending request to mechain.

3. The cmd tool has ability to intelligently select the correct SP by the info of bucket name and object name in command. Users do not need to specify the address of the SP in the command or config file.

4. The "mc://" is a fixed prefix which representing the mechain object or bucket.

### Examples

#### Init Accounts

Users can use "account import [keyfile]" to generate the keystore.  Before importing the key and generate keystore, you should export your private key from MetaMask and write it into a local keyfile as plaintext.

```
// import private key and generate keystore key.json
// The key.txt contain the plaintext private key. After the keystore file has been generated, user can delete the private key file key.txt.
mechain-cmd account import key.txt
```

The keystore will be generated in the path "keystore/keyfile" under the home directory of the system or the directory set by "--home"
and it is also the default path to load keystore when running other commands.
Password info is also needed to run the command. The terminal will prompt user to enter the password information.
Users can also specify the password file path by using the "--passwordfile".
Users are responsible for keeping their password information. If the password is lost, it is needed will need to re-import the key.

The "account new" command can be used to create a new account for executing commands. However, please note that after creating the account, you need to transfer token to the address of this account before you can send transactions or use storage-related functions.

```
// new account and generate keystore key.json
mechain-cmd account new
```

Users can use "account export" or "account ls" to display the keystore information of account.

```
// list the account info
mechain-cmd account ls
// export the account key info
mechain-cmd account export --unarmoredHex --unsafe
```

Users can create multiple accounts using the "account import" or "account new" command. You can use the "set-default" command to specify which account to use for running other commands by default. When executing commands using the default account, there is no need to specify the keystore.
You can also switch to different accounts for sending requests by specifying the --keystore flag.

```
// set the default account.
mechain-cmd account set-default [address-info]

// set keystore flag to use other account which is not default
mechain-cmd --keystore [keystore-path]  bucket create mc://test-bucket
```

#### Bank Operations

```
// transfer to an account in mechain
mechain-cmd bank transfer --toAddress 0xF678C3734F0EcDCC56cDE2df2604AC1f8477D55d --amount 12345

// crosschain transfer some tokens to an account in BSC
mechain-cmd bank bridge --toAddress 0xF678C3734F0EcDCC56cDE2df2604AC1f8477D55d --amount 12345

// query the balance of account
mechain-cmd bank balance --address 0xF678C3734F0EcDCC56cDE2df2604AC1f8477D55d

```

#### Storage Provider Operations

The "sp" command is used to list the SP and query the information of SP.

```
// list storage providers
mechain-cmd sp ls

// get storage provider info
mechain-cmd sp head https://mechain-testnet-sp-1.nodereal.io

// get quota and storage price of storage provider:
mechain-cmd sp get-price https://mechain-testnet-sp-1.nodereal.io
```

#### Bucket Operations

Before creating bucket, it is recommended to first run the "sp ls" command to obtain the SP list information of mechain,
and then select the target SP to which the bucket will be created on.

```
// create bucket. 
// The targt primary SP address to which the bucket will be created on need to be set by --primarySP flag.
// If the primary SP has not been not set, the cmd will choose first SP in the SP list which obtain from chain as the primary SP.
mechain-cmd bucket create mc://mechain-bucket

// update bucket visibility, charged quota or payment address
(1) mechain-cmd bucket update --visibility=public-read mc://mechain-bucket
(2) mechain-cmd bucket update --chargedQuota 50000 mc://mechain-bucket
```

#### Upload/Download Operations

(1) put Object

The "object put" command is used to upload a file from local which is less than 2G. The bucket name and object name should be replaced with specific names and
the file-path should replace by the file path of local system.

```
mechain-cmd object put --contentType "text/xml" --visibility private file-path mc://mechain-bucket/mechain-object
```

if the object name has not been set, the command will use the file name as object name.

(2) download object

The "object get" command is used to download an object to local path. This command will return the local file path where the object will be downloaded and the file size after successful execution.

```
mechain-cmd object get mc://mechain-bucket/mechain-object file-path 
```

The filepath can be a specific file path, a directory path, or not set at all.
If not set, the command will download the content to a file with the same name as the object name in the current directory. If it is set as a directory, the command will download the object file into the directory.

(3) create empty folder

Please note that the object name corresponding to the folder needs to end with "/" as suffix

```
mechain-cmd object put mc://mechain-bucket/folder/
```

(4) upload local folder

To upload a local folder (including all the files in it), you can use --recursive flag and specify the local folder path

```
mechain-cmd object put --recursive local-folder-path mc://mechain-bucket
```

(5) upload multiple files

To upload multiple files by one command, you can specify all the file paths that need to be uploaded one by one.
The files will be uploaded to the same bucket.

```
mechain-cmd object put  filepath1 filepath2 ...  mc://mechain-bucket
```

#### Group Operations

The group commands is used to create group, update group members, delete group and query group info.

```
// create group
mechain-cmd group create mc://groupname

// update group member
mechain-cmd group update --addMembers 0xca807A58caF20B6a4E3eDa3531788179E5bc816b mc://groupname

// head group member
mechain-cmd group head-member  0xca807A58caF20B6a4E3eDa3531788179E5bc816b mc://groupname

// delete group
mechain-cmd group delete mc://group-name
```

#### Policy Operations

The mechain-cmd policy command supports the policy for put/delete resources policy(including objects, buckets, and groups) to the principal.

The principal is need to be set by --grantee which indicates a mechain account or --groupId which indicates group id.

The object policy action can be "create", "delete", "copy", "get" , "execute", "list", "update" or "all".
The bucket policy actions can be "update", "delete",  "list", "update", "createObj", "deleteObj", "copyObj", "getObj", "executeObj" or "all".The actions which
contain Obj means it is an action for the objects in the bucket.
The group policy actions can be "update", "delete" or "all".

The policy effect can set to be "allow" or "deny" by --effect

Put policy examples:

```
// grant object operation permissions to a group
mechain-cmd policy put  --groupId 128  --actions get,delete  grn:o::mechain-bucket/mechain-object

// grant object operation permissions to an account
mechain-cmd policy put --grantee 0x169321fC04A12c16... --actions get,delete grn:o::mechain-bucket/mechain-object

// grant bucket operation permissions to a group
mechain-cmd policy put --groupId 130 --actions delete,update,createObj  grn:b::mechain-bucket

// grant bucket operation permissions to an account
mechain-cmd policy put --grantee 0x169321fC04A12c16... --actions delete,update  grn:b::mechain-bucket

// grant group operation permissions to an account 
mechain-cmd policy put --grantee 0x169321fC04A12c16... --actions update  grn:g:owneraddress:mechain-group
```

Delete policy examples:

```
// delete the bucket policy from a group
mechain-cmd policy delete --groupId 11  grn:b::mechain-bucket

// delete the object policy from an grantee
mechain-cmd policy delete --grantee 0x169321fC04A12c16...  grn:o::mechain-bucket/mechain-object

```

#### List Operations

```
// list buckets
mechain-cmd bucket ls

// list objects of the bucket
mechain-cmd object ls mc://mechain-bucket

// list objects of the bucket in a recursive way
mechain-cmd object ls --recursive mc://mechain-bucket

// list the objects by prefix 
mechain-cmd object ls --recursive mc://mechain-bucket/prefixName
```

#### Delete Operations

```
// delete bucekt
mechain-cmd bucket delete mc://mechain-bucket

//delete object
mechain-cmd object delete mc://mechain-bucket/mechain-object

```

#### Head Operations

```
// head bucekt
mechain-cmd bucket head mc://mechain-bucket

// head object
mechain-cmd object head mc://mechain-bucket/mechain-object

// head Group
mechain-cmd group head mc://groupname
```

#### Payment Operations

```
// create a payment account
mechain-cmd payment-account create

// list payment accounts under owner or a address with optional flag --user 
mechain-cmd payment-account ls --owner 0x5a64aCD8DC6Ce41d824638419319409246A9b41A

// deposit from owner's account to the payment account 
mechain-cmd payment-account deposit --toAddress 0xF678C3734F0EcDCC56cDE2df2604AC1f8477D55d --amount 12345

// witharaw from a payment account to owner's account
mechain-cmd payment-account withdraw --fromAddress 0xF678C3734F0EcDCC56cDE2df2604AC1f8477D55d --amount 12345
```

#### Quota Operations

```
// get quota info
mechain-cmd bucket get-quota mc://mechain-bucket

// buy quota
mechain-cmd bucket buy-quota --chargedQuota 1000000 mc://mechain-bucket
```

#### Resource mirror Operations

```
// mirror a group as NFT to BSC, you might use group id or groupName to identidy the group
mechain-cmd group mirror --id 1
or
mechain-cmd group mirror --groupName yourGroupName

// mirror a bucket as NFT to BSC, you might use bucket id or bucketName to identidy the bucket
mechain-cmd bucket mirror --id 1
or
mechain-cmd bucket mirror --bucketName yourBucketName

// mirror a object as NFT to BSC, you might use object id or (bucketName, objectName) to identidy the object
mechain-cmd object mirror --id 1
or
mechain-cmd object mirror --bucketName yourBucketName --objectName yourObjectName
```

## Reference

- [mechain](https://github.com/zkMeLabs/mechain): the mechain blockchain
- [mechain-Contract](https://github.com/zkMeLabs/mechain-contracts): the cross chain contract for mechain that deployed on BSC network. .
- [mechain-Storage-Provider](https://github.com/zkMeLabs/mechain-storage-provider): the storage service infrastructures provided by either organizations or individuals.
- [mechain-Relayer](https://github.com/zkMeLabs/mechain-relayer): the service that relay cross chain package to both chains.

## Fork Information

This project is forked from [greenfield-cmd](https://github.com/bnb-chain/greenfield-cmd). Significant changes have been made to adapt the project for specific use cases, but much of the core functionality comes from the original project.
