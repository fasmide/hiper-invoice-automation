# hipermads
Invoicing webcrawler for hal9k's cashier. 

As of 2022-03-08, the hal9k hackerspace changed internet provider from hiper.dk to kviknet.dk and as such, this 
piece of hacked togeather webcrawler/invoice-emailer needed to be changed to at least support kviknet. 
Now, theres seperate cli's, one for hiper and one for kviknet. 

## Usage
hipermads needs a lot of environment secrets in order to do its work, start of by creating a `.env` file
```
declare -x HIPER_EMAIL_FROM=
declare -x HIPER_EMAIL_PASSWORD=
declare -x HIPER_EMAIL_TO=
declare -x HIPER_ID_DB_PATH=
declare -x HIPER_MIDES_ACCOUNTNO=
declare -x HIPER_PASSWORD=
declare -x HIPER_USERNAME=
declare -x KVIKNET_ID_DB_PATH=
declare -x KVIKNET_PASSWORD=
declare -x KVIKNET_USERNAME=
```

then run eigher `go run cli/hiper/main.go` or `go run cli/kviknet/main.go`

