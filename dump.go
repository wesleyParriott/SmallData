package SmallData

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// TODO: INDEX OUT OF RANGE if file is empty
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

// Dump dumps the contents of the table into a file
// the name of the dump file is determined by the name given
// upon intialization of the hashtable
func (ht *HashTable) Dump() {
	if ht.CurrentEntries == 0 {
		log.Print("ERROR not dumping because there's nothing to dump!")
		return
	}

	fileName := ht.FileName
	fmt.Println(fileName)
	fmt.Println(ht.FileName)
	fmt.Println(ht.CurrentEntries)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 777)
	if err != nil {
		log.Printf("WARNING (when dumping) %s", err)
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

		buff := make([]byte, 2+int(dat.keyHashSize)+int(dat.valueSize))

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
