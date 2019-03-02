# SmallData

A simple, easy to use and build (no 3rd party dependencies!) Key/Value hash list to faciliate an in memory database.

## Getting Started

```Golang 
SmallData.NewTable(1024)

key   := "some key"
value := "some value"
err := htv.StoreString(key, value)
if err != nil {
    // handle err...
}

ret := htv.SearchString(key)
fmt.Println(ret)
// $> some value
```

The way collisions are handled is a set of semi-colon delimited values.

```Golang 
key := "some_key"
value01 := "some value"
value02 := "some other value"

err := htv.StoreString(key, value01)
if err != nil {
    // handle err ...
}
err := htv.StoreString(key, value02)
if err != nil {
    // handle err ...
}

ret := htv.SearchString(key)
fmt.Println(ret)
// $> some value;some other value
```

# Example Programs

Example programs that show you how the library could be used.

## cmd/cli: a command line interface 

### Installation
``` Bash
$ cd cmd/cli
$ go build
$ ./cli
```
### Usage
exit: exits out of the cli.

help: prints this helpful message.

insert: inserts into the hash table

        Example: $> insert key_name value [value...]


search: searchs for a value in the hash table based on a key

        Example: $> search key

        Return : $> value...

search_key: searchs for a value in the hash table based on a key

        Example: $> search_key 42402122

        Return : $> value... 


remove: removes a key and its associated value from the hash table

        Example: $> remove key

dump  : dumps to a config based file (default is dump.dat)


print : prints out the hash table and all it's values

        NOTE: 0's are null values in the hash table

stats : prints the max amount of entries and current amount of entries

byte_me: returns the go based array of bytes of a given string

hash_me: returns a hashed value based on the string given

## Tests

The testing is currently only in a happy path state. These tests are meant for regression testing for a continous
integreation tool. All the tests are located in SmallData_test.go.

To run all the tests:

```
go test -v 
```

### TestCoverage

Test coverage as of now is 80.2% according to 

```
go test -cover 
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* [It Looks Sad](https://itlookssad.bandcamp.com) for making dope music for me to program to.
