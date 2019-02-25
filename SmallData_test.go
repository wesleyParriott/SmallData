package SmallData

import (
    "os"
    "testing"
    "unsafe"
)

func TestNewTable(t *testing.T) {
    nominalTestTableSize := 8
    testHT := NewTable(nominalTestTableSize)
    testHTLen := len(testHT.Table)

    if testHTLen != nominalTestTableSize{
        t.Fatalf("The size of the test table was not %d but %d", nominalTestTableSize, testHTLen)
    }

    if testHT.MaxTableSize != testHTLen {
        t.Fatalf("The max table size was not %d but %d", nominalTestTableSize, testHTLen)
    }

    if testHT.FileName != "" {
        t.Fatalf("A dump file name was present. Was: %s", testHT.FileName)
    }

    t.Logf("PageSize %d", os.Getpagesize())
    t.Logf("Hash List Nominal Size  %d hash tables", nominalTestTableSize)
    t.Logf("Hash List Size in bytes %d bytes", unsafe.Sizeof(testHT.Table))
}

func TestNewTableFromNoFile(t *testing.T) {
    nominalTestTableSize := 8
    nominalFileName := "test_file.dat"

    testHT := NewTableFromFile(nominalFileName, nominalTestTableSize)
    testHTLen := len(testHT.Table)

    //only because we don't actuall have a test file should this be true
    if testHTLen != nominalTestTableSize{
        t.Fatalf("The size of the test table was not %d but %d", nominalTestTableSize, testHTLen)
    }


    if testHT.FileName != ("./" + nominalFileName) {
        t.Fatalf("nominal test file name [%s] and name found [%s]wasn't the same", testHT.FileName, nominalFileName)
    }

    t.Logf("PageSize %d", os.Getpagesize())
    t.Logf("Hash List Nominal Size  %d hash tables", nominalTestTableSize)
    t.Logf("Hash List Size in bytes %d bytes", unsafe.Sizeof(testHT.Table))
}

func TestNewTableFromAFile(t *testing.T) {
    nominalTestTableSize := 8
    nominalFileName := "test_data.dat"

    testHT := NewTableFromFile(nominalFileName, nominalTestTableSize)

    if testHT.FileName != ("./" + nominalFileName) {
        t.Fatalf("nominal test file name [%s] and name found [%s]wasn't the same", testHT.FileName, nominalFileName)
    }

    t.Logf("PageSize %d", os.Getpagesize())
    t.Logf("Hash List Nominal Size  %d hash tables", nominalTestTableSize)
    t.Logf("Hash List Size in bytes %d bytes", unsafe.Sizeof(testHT.Table))
}

func TestInsertion(t *testing.T) {
    nominalTestTableSize := 8
    testHT := NewTable(nominalTestTableSize)
    err := testHT.StoreString("key", "value")
    if err != nil {
        t.Fatalf("[STORE STRING] %s", err)
    }

    err = testHT.StoreBytes([]byte("key"), []byte("value"))
    if err != nil {
        t.Fatalf("[STORE BYTES] %s", err)
    }

    // the value will be a time stamp
    _, err = testHT.StoreBytesWithTimeStamp([]byte("value"))
    if err != nil {
        t.Fatalf("[STORE TIMESTAMP] %s", err)
    }

    t.Logf("PageSize %d", os.Getpagesize())
    t.Logf("Hash List Nominal Size  %d hash tables", nominalTestTableSize)
    t.Logf("Hash List Size in bytes %d bytes", unsafe.Sizeof(testHT.Table))
}

func TestRemoval(t *testing.T) {
    nominalTestTableSize := 8
    testHT := NewTable(nominalTestTableSize)
    testHT.StoreString("key", "value")

    err := testHT.Remove("key")
    if err != nil {
        t.Fatal(err)
    }

    if testHT.CurrentEntries != 0 {
        t.Fatal("the Current Entries was not detrimented")
    }
}

func TestSearch(t *testing.T) {
    nominalTestTableSize := 8
    testHT := NewTable(nominalTestTableSize)

    testHT.StoreString("Lester", "Vindaloo Addict")
    hashedKey := hashString("Lester")

    testData          := testHT.SearchForData(hashedKey)
    if testData == nil {
        t.Fatal("testData was returned nil")
    }

    keysearchvalue    := testHT.SearchKey(hashedKey) 
    if keysearchvalue != "Vindaloo Addict" {
        t.Fatal("doing a search with key didn't get the right value")
    }

    stringsearchvalue := testHT.SearchString("Lester")
    if stringsearchvalue != "Vindaloo Addict" { 
        t.Fatal("doing a search with string didn't get the right value")
    }
}

func TestDump(t *testing.T) {
    nominalTestTableSize := 8
    nominalFileName := "should_die.dat"

    testHT := NewTableFromFile(nominalFileName, nominalTestTableSize)

    if testHT.FileName == nominalFileName {
        t.Fatal("nominal file name was not congruent with the test file name")
    }

    testHT.StoreString("Lester", "Vindaloo Addict")

    testHT.Dump()

	_, err := os.Stat("./" + nominalFileName)
	if err != nil {
        t.Fatal("The nominal file didn't get created")
	}

    err = os.Remove("./should_die.dat")
    if err != nil {
        t.Log(err)
    }
}
