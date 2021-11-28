package main

import (
	"fmt"
)

type Person struct {
	Name    string
	Age     uint8
	Address Addr
}

type Addr struct {
	city     string
	district string
}

func testTranslateStruct() {
	personChan := make(chan Person, 1)
	person := Person{"xiaoluo", 18, Addr{"shenzhen", "bao'an"}}
	personChan <- person
	fmt.Printf("%v\n", person)

	person.Address = Addr{"guangzhou", "tianhe"}
	newPerson := <-personChan
	fmt.Printf("%v", newPerson)
}

func main() {
	testTranslateStruct()
}
