package main

import (
	"github.com/it-chain/it-chain-Engine/legacy/db/leveldbhelper"
	"log"
	"fmt"
	"os"
)

func main(){

	path := "./leveldb"
	defer os.RemoveAll(path)
	db := leveldbhelper.CreateNewDB(path)
	db.Open()

	err := db.Put([]byte("20164403"),[]byte("JUN"), true)

	if err != nil{
		log.Fatalf("error occured [%s]",err.Error())
	}

	name, _ := db.Get([]byte("20164403"))

	fmt.Printf("%s",name)
}
