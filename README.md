# go-payfee

POC for test purposes.

CRUD an REDIS cluster.

## database

Redis Cluster Database (aws MemoryDB)

## Endpoints

+ POST /script/add

        {
        "script": {
                "name": "script.debit",
                "description": "script for credit operation",
                "fee": ["fee-3.5","fee-05"]
            }
        }

+ GET /script/get/script.debit

+ POST /key/add

        {
            "name":"fee-3.5",
            "value": 3.5
        }

+ GET /key/get/fee-3.5
