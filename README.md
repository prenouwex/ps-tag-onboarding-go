# PS Tag Onboarding Go

Tag Onboarding Project exercise using Go-Chi, Gorm, Sqlite in-memory


## Build

```
./scripts/build.sh
```

This will generate `ps-tag-onboarding-go` executable.

## Run

```
./scripts/run.sh
```

This will start the `ps-tag-onboarding-go` server.

Note that the URI of this application is `http://localhost:8089/users/` which should provide you access to the CRUD operations.

The above URL when invoked, would return the following user list:

```
[{"id":1,"first_name":"John","last_name":"Doe","email":"john.doe@yahoo.com","age":34},{"id":2,"first_name":"Zenia","last_name":"Brennan","email":"ultrices.vivamus.rhoncus@yahoo.ca","age":34},{"id":3,"first_name":"Branden","last_name":"Spears","email":"non.lobortis@hotmail.net","age":34},{"id":4,"first_name":"Alice","last_name":"Wallace","email":"at@protonmail.couk","age":34},{"id":5,"first_name":"Ira","last_name":"Francis","email":"in.lobortis.tellus@protonmail.ca","age":34}]
```


## Design Decisions
Here is the list of external libraries and the reason why they have been imported: 

### Go-Chi 
Lightweight, idiomatic and composable router for building Go HTTP services

### Gorm
ORM for Go, it helps separate out DB specifics and business logic
Robust abstraction layer between the Go code and the database,

### SQlite In memory 
Easy, fast and handy to use for small project and prot