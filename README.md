# golang-project
Learning golang, created this project using https://youtu.be/tA27VcHUK48 as a guide.

## Overview
* A database for in-memory store.
* Also provides an option to save the key-value to a db.json file.

## Instructions to run the project: 

1. Either run the project directly
```go
go run main.go
```
OR build a .exe file and run the executable ( for Windows only )
```go
go build
```

2. The execution file starts the server 

3. You need to initiate the server in a different terminal
```
telnet localhost:8080
```

4. That should return a message in the current terminal saying: 
> --- Welcome to RuntimeDB server
