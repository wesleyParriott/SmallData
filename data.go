package SmallData

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

/*
	file data consists of
	[0]       = max data entries
	[1...n]   = key
	[n...n-1] = value
	continue finding key/value till EOF
*/

// HashTableData is the data structure that holds the basic Key/Value information of every entry
type HashTableData struct {
	Key   int
	Value []byte
}

// HashTable is the data structure for the HashTable itself which is the array of HashTableData
type HashTable struct {
	Table []HashTableData

	CurrentEntries int
	MaxTableSize   int

	FileName string
}

type packedInfo struct {
	// TODO timestamp [64]byte
	keyHash     []byte // int32
	keyHashSize byte
	valueSize   byte // int32
	value       []byte
}

func packData(h HashTableData) packedInfo {
	var pki packedInfo

	if h.Key < 0 {
		log.Printf("WARNING trying to pack a HashTableData with a negitive key value!")
		return pki
	}

	// cast that shit like
	// colonists
	// nothing bad will
	// come of this
	keyS := strconv.Itoa(h.Key)

	pki.keyHash = []byte(keyS)
	pki.keyHashSize = byte(len(pki.keyHash))

	pki.value = []byte(h.Value)
	pki.valueSize = byte(len(pki.value))

	return pki
}

var maxTableSize int

// NewTable returns a New HashTable that is the size of
// the given size * the size of a HashTable
func NewTable(size int) *HashTable {
	ht := new(HashTable)
	ht.Table = make([]HashTableData, size)
	ht.MaxTableSize = size
	maxTableSize = size

	return ht
}

// NewTable from file returns a New Hash table based on a dump file given
// it will be the size of the file
func NewTableFromFile(fileName string) *HashTable {
	ht := new(HashTable)

	ht.FileName = fileName

	filePath := "./" + fileName
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Warning When trying to open dump file: %s", err)

		size := 1024
		log.Printf("INFO Now loading file with size: %d bytes\n", size)
		return NewTable(size)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		log.Printf("WARNING When trying to open dump file: %s", err)
		return ht
	}
	fileSize := int(fileInfo.Size())

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("FATAL When trying to get dump file contents: %s", err)
	}

	size := int(fileContents[0])

	log.Printf("INFO loaded file with size: %d max entries\n", size)

	ht.Table = make([]HashTableData, size)
	ht.MaxTableSize = size
	maxTableSize = size

	stream := make([]byte, fileSize)
	n, err := f.Read(stream)
	if err != nil {
		log.Printf("WARNING When trying to read dump file: %s", err)
		return ht
	}
	log.Printf("INFO %d bytes read from dump file", n)
	loadData(stream, ht)

	return ht
}
