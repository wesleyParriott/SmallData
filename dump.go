package SmallData

import (
	"os"
)

// Dump dumps the contents of the table into a file
// the name of the dump file is determined by the name given
// upon intialization of the hashtable
func (ht *HashTable) Dump() {
	if ht.CurrentEntries == 0 {
		warning("Not dumping because there's nothing to dump!")
		return
	}

	totalSize := byte(ht.MaxTableSize)
	if totalSize == 0 {
		warningf("Not dumping the contents of the table because the size %d is larger than 255 bytes!", ht.MaxTableSize)
		return
	}

	fileName := ht.FileName
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 777)
	if err != nil {
		warningf("error when dumping: %s", err)
		return
	}
	defer f.Close()

	maxTableSizeBuff := make([]byte, 1)
	maxTableSizeBuff[0] = totalSize
	_, err = f.Write(maxTableSizeBuff)
	if err != nil {
		warningf("Not dumping contents because couldn't write maxTable data: %s", err)
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
			if i > len(buff) {
				break
			}
			buff[buffIndex] = dat.keyHash[i]
			buffIndex++
		}

		for i := 0; i < int(dat.valueSize); i++ {
			if i > len(buff) {
				break
			}
			buff[buffIndex] = dat.value[i]
			buffIndex++
		}

		_, err := f.Write(buff)
		if err != nil {
			warning("breaking out of loop")
			break
		}
	}
}
