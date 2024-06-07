<!-- Code for [Bitespeed challenge](https://www.notion.so/bitespeed/Bitespeed-Backend-Task-Identity-Reconciliation-53392ab01fe149fab989422300423199) -->

Linking multiple contact informations to the same person using primary and secondary contacts

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

# Example Request  

```json
{
	"email": "random@pennyworth.edu",
	"phoneNumber": "898991"
}
```

# Example Response

```json
{
    "primaryContatctId": 51,
    "emails": [
        "babu1@digital.com",
        "random@pennyworth.edu"
    ],
    "phoneNumbers": [
        "898991",
        "898991"
    ],
    "secondaryContactIds": [
        52
    ]
}
```

# Example Request  

```json
{
	"email": "random@pennyworth.edu",
	"phoneNumber": "123456"
}
```

# Example Response

```json
{
    "primaryContatctId": 52,
    "emails": [
        "random@pennyworth.edu",
        "random@pennyworth.edu"
    ],
    "phoneNumbers": [
        "898991",
        "123456"
    ],
    "secondaryContactIds": [
        52,
        53
    ]
}
```


# Example Request  

```json
{
	"phoneNumber": "123456"
}
```

# Example Response

```json
{
    "primaryContatctId": 53,
    "emails": [
        "random@pennyworth.edu"
    ],
    "phoneNumbers": [
        "123456"
    ],
    "secondaryContactIds": [
        53
    ]
}
```

# Example Request  

```json
{
	"email": "random@pennyworth.edu"
}
```

# Example Response

```json
{
    "primaryContatctId": 51,
    "emails": [
        "babu1@digital.com",
        "random@pennyworth.edu"
    ],
    "phoneNumbers": [
        "898991",
        "898991"
    ],
    "secondaryContactIds": [
        52
    ]
}
```

