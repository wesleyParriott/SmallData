package SmallData

import (
    "fmt"
    "testing"
)
/*
TODO:
    search with one key and get all values
*/

/*
CASES:
    
*/

func testPrint(htv *HashTable) {
    fmt.Printf("================\n")
    htv.Print()
}

func TestCollisions(t *testing.T) {
    n := 4
    htv := NewTable(n)

    for i := 0; i < n; i++ {
        htv.StoreString("smeg", fmt.Sprintf("%d", i))
        testPrint(htv)
    }
    testPrint(htv)
    htv.StoreString("smeg", "lastOne")
    testPrint(htv)

    // v := SearchString("smeg")
}
