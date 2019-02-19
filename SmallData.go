package SmallData

import (
    "fmt"
    "hash/adler32"
    "io/ioutil"
    "log"
    "os"
    "strconv"
    "time"
    _"unsafe"
)

/* 
TODO:
   take file structure stuff to its own file 
       file data consists of
       [0]       = max data entries
       [1...n]   = key
       [n...n-1] = value
       continue finding key/value till EOF
*/

// HashTableData is the data structure that holds the basic Key/Value information of every entry
type HashTableData struct {
    Key int
    Value []byte
}

// HashTable is the data structure for the HashTable itself which is the array of HashTableData
type HashTable struct {
    Table []HashTableData

    CurrentEntries int
    MaxTableSize int

    FileName string
}

type packedInfo struct {
    // TODO timestamp [64]byte
    keyHash   []byte // int32 
    keyHashSize byte
    valueSize byte // int32
    value     []byte
}

func packData(h HashTableData) packedInfo{
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

func loadData(stream []byte, ht *HashTable) {
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
            log.Fatalf("Error when loading data: %s", err)
        }
        insert(ht, key, v)
    }
}

var maxTableSize int

// NewTable returns a New HashTable that is the size of 
// the given size * the size of a HashTable
func NewTable(size int) (*HashTable) {
    ht := new(HashTable)
    ht.Table = make([]HashTableData, size)
    ht.MaxTableSize = size
    maxTableSize = size

    return ht
}

// NewTable from file returns a New Hash table based on a dump file given
// it will be the size of the file
func NewTableFromFile(fileName string) (*HashTable) {
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


func hashString(s string) int {
    h := adler32.New()
    h.Write([]byte(s))
    return int(h.Sum32())
}

func hashCode(k int ) int {
    return k % maxTableSize
}


// just a wrapper for testing some stuff
func BunkHashString(s string) int {
    return hashString(s)
}
func BunkHashCode(k int ) int {
    return hashCode(k)
}
//

func insert(hashTable *HashTable, k int , v []byte) error{

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
    for hashTable.Table[index].Key != 0  && hashTable.Table[index].Value != nil {
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

//Remove removes a value from the hash table by nulling the value and assigning a 0 to the key 
func (ht *HashTable) Remove(k string) error {
    item := ht.SearchForData(hashString(k))
    if item == nil {
        return fmt.Errorf("removal didn't occur because item with key of %s was null when removing", k)
    }

    item.Key   = 0
    item.Value = nil
    ht.CurrentEntries -= 1

    return nil
}

// Dump dumps the contents of the table into a file 
// the name of the dump file is determined by the name given
// upon intialization of the hashtable
func (ht HashTable) Dump() {
    if ht.CurrentEntries == 0 {
        log.Print("ERROR not dumping because there's nothing to dump!")
        return
    }

    fileName := ht.FileName
    f, err:= os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 777)
    if err != nil {
        log.Printf("WARNING couldn't open the file for dumping because: %s", err)
        return
    }

    maxTableSizeBuff := make([]byte, 1)
    maxTableSizeBuff[0] = byte(ht.MaxTableSize)
    _, err = f.Write(maxTableSizeBuff)
    if err != nil {
        log.Printf("WARNING Not dumping contents because couldn't write maxTable data: %s", err)
        return
    }
    for i := 0; i < ht.MaxTableSize; i++ {
        h := ht.Table[i]

        if h.Key == 0 {
            continue
        }

        dat := packData(h)

        buff := make([]byte, 2 + int(dat.keyHashSize) + int(dat.valueSize))

        buff[0] = dat.keyHashSize
        buff[1] = dat.valueSize
        buffIndex := 2

        for i := 0; i < int(dat.keyHashSize); i++ {
            // TODO add fail safes
            buff[buffIndex] = dat.keyHash[i]
            buffIndex++
        }

        for i := 0; i < int(dat.valueSize); i++ {
            // TODO add fail safes
            buff[buffIndex] = dat.value[i]
            buffIndex++
        }

        _, err := f.Write(buff)
        if err != nil {
            log.Printf("WARNING breaking out of loop")
            break
        }
    }
}

func (ht HashTable) Print() {
    for _, h := range ht.Table {
        fmt.Printf("%d:%s\n", h.Key, h.Value)
    }
}

func (ht HashTable) Stats() {
    fmt.Printf("Max Entries: %d\nCurrent Entries:%d\n", ht.MaxTableSize, ht.CurrentEntries)
}

func search(hashTable *HashTable, k int) *HashTableData{
    index := hashCode(k)

    for _, _ = range hashTable.Table{
        if hashTable.Table[index].Key == k {
            return &hashTable.Table[index]
        }

        index++
        index %= maxTableSize
    }

    return nil
}

// Search searchs for the HashTableData entry of the value of the given int key
func (ht *HashTable) SearchForData(k int) *HashTableData{
    return search(ht, k)
}

// SearchString searchs for the string Value of the given string Key
func (ht *HashTable) SearchString(k string) string{
    dat := ht.SearchForData(hashString(k))

    if(dat == nil) {
        log.Printf("[Warning] value for key: %s not found", k)
        return ""
    }

    return string(dat.Value)
}

// TODO change to both the key and value being strings
func (ht *HashTable) StoreString(key string, value string) error {
    return insert(ht, hashString(string(key)), []byte(value))
}

func (ht *HashTable) StoreBytes(key []byte, value []byte) error {
    return insert(ht, hashString(string(key)),value)
}

func (ht *HashTable) StoreBytesWithTimeStamp(input []byte) (int, error) {
    t := time.Now().UTC()
    hashed := hashString(t.String())
    err := insert(ht, hashed, input)

    return hashed, err
}
