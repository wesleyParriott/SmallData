package SmallData

import (
	"fmt"
	"time"
)

// TODO:
//      batch inserts
//      remove by hashed key

func insert(hashTable *HashTable, k int, v []byte) error {

	if hashTable.CurrentEntries >= maxTableSize {
		return fmt.Errorf("insertion of key/value %d %s not completed due to the datastore being full!", k, v)
	}

	var item HashTableData
	item.Key = k
	item.Value = v

	index := hashCode(item.Key)

	/*
	   Right now this is how collisons are handled. We delimit the value with a semicolon.
	   Now, this may have some problems in the future but we'll figure those out as we need to.
	*/
	for hashTable.Table[index].Key != 0 && hashTable.Table[index].Value != nil {
		if hashTable.Table[index].Key == item.Key {
			hashTable.Table[index].Value = append(hashTable.Table[index].Value, ';')
			for i := 0; i < len(item.Value); i++ {
				hashTable.Table[index].Value = append(hashTable.Table[index].Value, item.Value[i])
			}
			return nil
		}
		index++
		index %= maxTableSize
	}

	hashTable.CurrentEntries++
	hashTable.Table[index] = item

	return nil
}

// StoreString stores the a given value into the table index by a hash created by the given key
func (ht *HashTable) StoreString(key string, value string) error {
	return insert(ht, hashString(string(key)), []byte(value))
}

// StoreBytes is the same as StoreString but use type []byte for the key and value
func (ht *HashTable) StoreBytes(key []byte, value []byte) error {
	return insert(ht, hashString(string(key)), value)
}

// StoreBytesWithTimeStamp will use time of storage to create a hash key and will store the given byte array
// it will then return the hashed integer value for searching later on
func (ht *HashTable) StoreBytesWithTimeStamp(input []byte) (int, error) {
	// NOTE I'm not sure how useful it is to utilize the returning of the hash key
	//      along with how useful this function even is
	t := time.Now().UTC()
	hashed := hashString(t.String())
	err := insert(ht, hashed, input)

	return hashed, err
}

//Remove removes a value from the hash table by nulling the value and assigning a 0 to the key
func (ht *HashTable) Remove(k string) error {
	item := ht.SearchForData(hashString(k))
	if item == nil {
		return fmt.Errorf("removal didn't occur because item with key of %s was null when removing", k)
	}

	item.Key = 0
	item.Value = nil
	ht.CurrentEntries -= 1

	return nil
}
