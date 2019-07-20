package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wesleyParriott/SmallData"
)

const usage = `
Small Data: a tiny hash table cli
    ==COMMANDS==
    [MAIN]
    exit: exits out of the cli.

    help: prints this helpful message.

    insert: inserts into the hash table
            Example: $> insert key_name value [value...] 

    search:     searchs for a value in the hash table based on a given string
            Example: $> search some_key
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

    [UTIL]
    byte_me: returns the go based array of bytes of a given string  

    hash_me: returns a hashed value based on the string given

`

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("FATAL %s", err)
	}
}

func warnOnError(err error) {
	if err != nil {
		log.Printf("WARNING %s", err)
	}
}

func clean(input string) string {
	var ret []byte

	for i := 0; i < len(input); i++ {
		if input[i] == '\r' || input[i] == '\n' {
			continue
		}

		ret = append(ret, input[i])
	}

	return string(ret)
}

func combine(inputs []string) (ret string, err error) {
	if len(inputs) == 0 {
		err = fmt.Errorf("length of given inputs was 0 when combining")
		return
	}
	if len(inputs) == 1 {
		ret = inputs[0]
		return
	}

	for _, input := range inputs {
		if ret == "" {
			ret = input
			continue
		}
		ret = fmt.Sprintf("%s %s", ret, input)
	}
	return
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var htv *SmallData.HashTable

	htv = SmallData.NewTableFromFile("test_data.dat", 8)

LOOP:
	for {
		fmt.Printf("[SmallData] ")
		input, err := reader.ReadString('\n')
		fatalOnError(err)

		input = clean(input)
		fatalOnError(err)

		inputs := strings.Split(input, " ")
		switch inputs[0] {
		case "exit":
			break LOOP
		case "help":
			fmt.Printf(usage)
		case "insert":
			// $> insert makir_shada_kazim 2nd level barbarian
			// it's space delimited making it thus:
			//  // [insert,makir_shada_kazim,2nd,level,barbarian]
			//  // command = [0] insert
			//  // key     = [1] makir_shada_kazim
			//  // value1  = [2] 2nd
			//  // value2  = [3] level
			//  // value3  = [4] barbarian
			if len(inputs) < 3 {
				log.Printf("improper amount of inputs given. Example: insert key_name value and stuff")
				break
			}
			key := inputs[1]
			value, err := combine(inputs[2:])
			warnOnError(err)
			err = htv.StoreString(key, value)
			if err != nil {
				log.Printf("WARNING %s", err)
			}
		case "search":
			if len(inputs) < 2 {
				log.Printf("improper amount of search parameters given. Example: insert key_name")
			}
			// $> search makir_shada_kazim
			// it's space delimited making it thus:
			//  // [search,makir_shada_kazim]
			//  // command = [0] search
			//  // key     = [1] makir_shada_kazim
			//  we should then get a value back from the database
			key := inputs[1]
			value := htv.SearchString(key)
			fmt.Println(value)
		case "search_key":
			if len(inputs) < 2 {
				log.Printf("improper amount of search parameters given. Example: insert key_name")
			}
			// $> search_key 1045563120
			// it's space delimited making it thus:
			//  // [search,1045563120]
			//  // command = [0] search_key
			//  // key     = [1] 1045563120
			//  we should then get a value back from the database
			key, err := strconv.Atoi((inputs[1]))
			if err != nil {
				log.Fatalf("FATAL couldn't convert key to int: %d", input[1])
			}
			value := htv.SearchKey(key)
			fmt.Println(value)
		case "remove":
			if len(inputs) < 2 {
				log.Printf("improper amount of search parameters given. Example: insert key_name")
			}
			// $> search makir_shada_kazim
			// it's space delimited making it thus:
			//  // [search,makir_shada_kazim]
			//  // command = [0] remove
			//  // key     = [1] makir_shada_kazim
			//  if we print the table we should see
			//  that makir_shada_kazim is no longer there
			key := inputs[1]
			htv.Remove(key)
		case "dump":
			htv.Dump()
		case "print":
			htv.Print()
		case "stats":
			htv.Stats()
		case "time":
			fmt.Println(time.Now().Unix())
		case "byte_me":
			if len(inputs) < 2 {
				log.Printf("Need Something to byte! Example: byte_me some_string_value")
			}
			n, _ := strconv.Atoi(inputs[1])
			fmt.Printf("%v\n", byte(n))
		case "hash_me":
			if len(inputs) < 2 {
				log.Printf("Need Something to byte! Example: byte_me some_string_value")
			}
			hash := SmallData.BunkHashString(string(inputs[1]))
			hashCode := SmallData.BunkHashCode(hash)
			fmt.Printf("Hash: %d HashCode: %d\n", hash, hashCode)
		}
	}
}
