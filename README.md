# go-payfee

POC for test purposes

CRUD an REDIS cluster.

## database

Redis Cluster Database (aws MemoryDB)

## Endpoints

+ GET /info

+ POST /add/script

        {
        "script": {
                "name": "script.debit",
                "description": "script for credit operation",
                "fee": ["fee-3.5","fee-05"]
            }
        }

+ GET /script/script.debit

+ POST /add/key

        {
            "name":"fee-3.5",
            "value": 3.5
        }

+ GET /key/fee-3.5
