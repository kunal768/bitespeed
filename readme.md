Code for [Bitespeed challenge](https://www.notion.so/bitespeed/Bitespeed-Backend-Task-Identity-Reconciliation-53392ab01fe149fab989422300423199)

# To Build

```shell
go build -o bitespeed
```

# To Run 

```shell
./bitespeed
```

# Example Request  

```json
{
	"email": "babu1@digital.com",
	"phoneNumber": "898991"
}
```

# Example Response

```json
{
    "primaryContatctId": 46,
    "emails": [
        "babu1@digital.com"
    ],
    "phoneNumbers": [
        "898991"
    ],
    "secondaryContactIds": []
}
```

