#!/usr/bin/env sh

curl -v --location 'localhost:8080/' \
--header 'Content-Type: application/json' \
--data '{
    "header": {
        "ADIF_VER": "3.1.5"
    },
    "records": [
        {
            "BAND": "40m",
            "CALL": "K9CTS"
        },
        {
            "BAND": "40m",
            "CALL": "W9PVA"
        }
    ]
}'