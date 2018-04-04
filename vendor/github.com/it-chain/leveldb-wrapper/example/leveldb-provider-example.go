package main

import (
	"github.com/it-chain/leveldb-wrapper"
	"fmt"
	"os"
)

func main(){

	path := "./leveldb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	defer os.RemoveAll(path)

	studentDB := dbProvider.GetDBHandle("Student")
	studentDB.Put([]byte("20164403"),[]byte("JUN"),true)

	name, _ := studentDB.Get([]byte("20164403"))

	fmt.Printf("%s",name)
}