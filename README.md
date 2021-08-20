# Boba API

This is the backend to the [Boba App](https://github.com/colinpcurtis/boba-api).  

## Setup
To initialize a go module, run
```
go mod init server
```

Make a `.env` file to hold environment variables required to run the app.
```
MONGO_URL=(mongo cluster url)
```

Then to install the required dependencies run 
```
go get
```

## Run the Server
Run 
```
go build
```
to compile the code, then 
```
./server
```
to run the server
