Golang PoC for the OneDrive user enumeration technique (based on the ["OneDrive to Enum Them All"](https://www.trustedsec.com/blog/onedrive-to-enum-them-all/) post by TrustedSec)

```
./onedrive-enum -h
Usage of ./onedrive-enum:
  -d string
        Domain
  -j string
        Log results to JSON file (default "results.json")
  -u string
        Userlist
  -w int
        Workers (default 3)
```

Usage:

```
./onedrive-enum -d DOMAIN -u USERLIST
```

Selecting valid users (status_code == 403) from JSON file:

```
cat results.json | jq -r 'select(.status_code == 403) | .username'
```