# GedditQL - A DB built for educational purposes on golang

GedditQL is a database built with Golang for our database class. It isn't a fully fledged database management system where it is reliable for production use but was built with the aim of learning how database management systems work well.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You need need golang installed on your system that has a verion of >= 1.8.

### Installing it on your system

Make sure you have the minimum go libraries installed properly in your `$GOPATH`.
In addition, for the client side application, you would need to install a separate repo listed below.

` go get -u github.com/Yong-L/go-texttable `

### Running it on your system

You need to run two terminal instances to run both the server and the client. The server will be served on `port 8888`. 

To run the server run: `go run main.go` from the root directory.

To run the client run: `go run client/client.go` from the root directory.

### Supported SQL Queries

This database supports very minimal database queries such as 

```
SELECT
CREATE
INSERT
UPDATE
DELETE
```

It also supports minimal function with select statements such as 

```
SUM
COUNT
```

### TODO 

- [ ] Mutexes for race conditions
