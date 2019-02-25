package SmallData

import (
	"log"
)

func search(hashTable *HashTable, k int) *HashTableData {
	index := hashCode(k)

	for _, _ = range hashTable.Table {
		if hashTable.Table[index].Key == k {
			return &hashTable.Table[index]
		}

		index++
		index %= maxTableSize
	}

	return nil
}

// Search searchs for the HashTableData entry of the value of the given int key
func (ht *HashTable) SearchForData(k int) *HashTableData {
	return search(ht, k)
}

// SearchString searchs for the string Value of the given string Key
func (ht *HashTable) SearchString(k string) string {
	dat := ht.SearchForData(hashString(k))

	if dat == nil {
		log.Printf("[Warning] value for key: %s not found", k)
		return ""
	}

	return string(dat.Value)
}

// SearchString searchs for the string Value of the given string Hashed Key Int
func (ht *HashTable) SearchKey(k int) string {

	dat := search(ht, k)

	if dat == nil {
		log.Printf("[Warning] value for key: %d not found", k)
		return ""
	}

	return string(dat.Value)
}
