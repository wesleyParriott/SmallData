package SmallData

import (
	"io/ioutil"
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

// HashTable is the data structure for the HashTable itself
// which holds the array of HashTableData, the current entry count, the max table size, and the dump file name
type HashTable struct {
	Table []HashTableData

	CurrentEntries int
	MaxTableSize   int

	FileName string
}

type packedInfo struct {
	keyHash     []byte // int32
	keyHashSize byte
	valueSize   byte // int32
	value       []byte
}

func loadData(stream []byte, ht *HashTable) {
	// early out in case of an empty stream
	if len(stream) == 0 {
		warning("stream was empty. Not loading data")
		return
	}

	index := 1
	for {
		if index >= len(stream) {
			break
		}
		khsb := stream[index]
		if khsb == '\x00' {
			// the cases where there's
			// uninitialized hashtable data inserted
			// because the whole thing wasn't filled up yet
			index++
			continue
		}
		khs := int(khsb)
		index++

		vs := int(stream[index])
		index++

		kh := stream[index:(khs + index)]
		index += khs

		v := stream[index:(vs + index)]
		index += vs

		key, err := strconv.Atoi(string(kh))
		if err != nil {
			// that's pretty baaaaaad if this happens
			fatalf("Error when loading data: %s", err)
		}
		insert(ht, key, v)
	}
}

func packData(h HashTableData) packedInfo {
	var pki packedInfo

	if h.Key < 0 {
		warning("WARNING trying to pack a HashTableData with a negitive key value!")
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

// NewTable from file returns a New Hash table based on a dump file given it will be the size of the file.
// If the file cannot be found a clean table will be produced with the default size given.
func NewTableFromFile(fileName string, defaultSize int) *HashTable {
	if defaultSize > 255 {
		fatalf("Default size cannont be larger than 255 bytes. Recieved a default size of %d", defaultSize)
	}

	filePath := "./" + fileName

	f, err := os.Open(filePath)
	if err != nil {
		warningf("When trying to open dump file: %s", err)

		infof("Now loading file with default size: %d bytes\n", defaultSize)
		ht := NewTable(defaultSize)
		ht.FileName = filePath

		f.Close()

		return ht
	}
	defer f.Close()

	ht := new(HashTable)
	ht.FileName = filePath

	fileInfo, err := f.Stat()
	if err != nil {
		warningf("When trying to open dump file: %s", err)
		return ht
	}
	// NOTE: I think file size is a bad name here. Because it's the
	//       it's the table size given from the file read into the program
	fileSize := int(fileInfo.Size())
	if fileSize == 0 {
		fatalf("the file size of the given file %s was 0!", filePath)
	}

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fatalf("When trying to get dump file contents: %s", err)
	}

	size := int(fileContents[0])
	if size <= 0 {
		fatalf("Size of Hash table from file %s is an invalid size. The size found was %d", fileName, size)
	}

	infof("loaded file with size: %d max entries\n", size)

	ht.Table = make([]HashTableData, size)
	ht.MaxTableSize = size
	maxTableSize = size

	stream := make([]byte, fileSize)
	n, err := f.Read(stream)
	if err != nil {
		warningf("When trying to read dump file: %s", err)
		return ht
	}
	infof("%d bytes read from dump file", n)
	loadData(stream, ht)

	return ht
}
