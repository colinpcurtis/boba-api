# Boba API

[![GitHub Super-Linter](https://github.com/colinpcurtis/boba-api/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

This is the backend to the [Boba App](https://github.com/colinpcurtis/boba-app).  

## Setup
To initialize a go module, run
```bash
go mod init server
```

Make a `.env` file to hold environment variables required to run the app.
```bash
MONGO_URL=(mongo cluster url)
```

Then to install the required dependencies run 
```bash
go get
```

## Run the Server
Run 
```bash
go build
```
to compile the code, then 
```bash
./server
```
to run the server
