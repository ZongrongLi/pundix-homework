# pundix-homework
- use cosmos-sdk to pack rpc request
- all denom set to default "FX"
- pick useraccount from expolorer randomly 
- validator in exam set to singapore



demand:
```ref
Using Golang, create a simple REST API that will query information from our mainnet blockchain.

eg. http://localhost:3000/query/bank/total will return a json on the browser with the output, just display the information, no css required.

The information query for the API is only for query/bank and query/distribution

Ensure proper error handling if the wrong ‘input’ field is called


Note:

You will want to first download the fxcored CLI tool in your machine/cloud instance to test out some of the following commands. Add the following flag to your command `--node https://fx-json.functionx.io:26657`  (you can follow the guide in our documentation to do so, https://functionx.gitbook.io/home/f-x-core/installation) :

	fxcored query bank

	fxcored query distribution

Eg. to query the total supply of coins on the chain `fxcored query bank total —node https://fx-json.functionx.io:26657` 



You might want to use a free cloud service(eg. AWS, Google cloud etc.) to connect your fxcored commands to your backend golang program.

You can use the addresses (such as the validator’s address) you find on the mainnet explorer, https://explorer.functionx.io/fxcore/validators. You are free to use any packages.

Upload your code into github with proper commit messages.
```

https://lizongrong.club/query/distribution/queryParams
```json
{
    "params":{
        "community_tax":"0.400000000000000000",
        "base_proposer_reward":"0.010000000000000000",
        "bonus_proposer_reward":"0.040000000000000000",
        "withdraw_addr_enabled":true
    }
}
```

https://lizongrong.club/query/distribution/validatorCommission?validator=fxvaloper1a73plz6w7fc8ydlwxddanc7a239kk45jnl9xwj
```json
{
    "commission":{
        "commission":[
            {
                "denom":"FX",
                "amount":"146236838209477174104602.042504533824137437"
            }
        ]
    }
}
```


https://lizongrong.club/query/distribution/validatorOutstandingRewards?validator=fxvaloper1a73plz6w7fc8ydlwxddanc7a239kk45jnl9xwj
```json
{
    "rewards":{
        "rewards":[
            {
                "denom":"FX",
                "amount":"1704072357231688849603690.428541765671723851"
            }
        ]
    }
}
```

https://lizongrong.club/query/distribution/communityPool

```json
{
    "pool":[
        {
            "denom":"FX",
            "amount":"40496599530653244544335396.361926119997781247"
        }
    ]
}
```


https://lizongrong.club/query/bank/balance?address=fx15sy7ph7j6vma607y80cxdc7qg7pgvjdhnql3q6

```json
{
    "balance":{
        "denom":"FX",
        "amount":"1375124410302910720"
    }
}
```


https://lizongrong.club/query/bank/total

```json
{
    "amount":{
        "denom":"FX",
        "amount":"498642746962843784115573074"
    }
}
```