package SmallData

import (
	"fmt"
	"hash/adler32"
)

func hashString(s string) int {
	h := adler32.New()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func hashCode(k int) int {
	return k % maxTableSize
}

// just a wrapper for testing some stuff
func BunkHashString(s string) int {
	return hashString(s)
}
func BunkHashCode(k int) int {
	return hashCode(k)
}

func (ht HashTable) Print() {
	for _, h := range ht.Table {
		fmt.Printf("%d:%s\n", h.Key, h.Value)
	}
}

func (ht HashTable) Stats() {
	debugf("Max Entries: %d\nCurrent Entries:%d\n", ht.MaxTableSize, ht.CurrentEntries)
}
