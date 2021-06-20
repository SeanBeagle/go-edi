package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Interchange struct {
	Header            Segment
	Trailer           Segment
	FunctionalGroups  []FunctionalGroup
	segmentTerminator byte
}

type FunctionalGroup struct {
	Header          Segment
	Trailer         Segment
	TransactionSets []TransactionSet
}

type TransactionSet struct {
	Header   Segment
	Trailer  Segment
	Segments []Segment
}

type Segment struct {
	Id                string
	Elements          []Element
	elementSeparator  string
	segmentTerminator byte
}

type Element struct {
	Id    string
	Value string
}

func NewInterchange(file string) *Interchange {
	segmentTerminator := GetSegmentTerminator(file)
	segments := GetSegments(file, segmentTerminator)

	header, trailer := segments[0], segments[len(segments)-1]
	if header.Id != "ISA" {
		panic("Expected ISA as first segment, not " + header.Id)
	} else if trailer.Id != "IEA" {
		panic("Expected ISE as last segment, not " + trailer.Id)
	}

	functionalGroups := GetFunctionalGroups(segments[1 : len(segments)-1])
	interchange := Interchange{Header: header, Trailer: trailer, FunctionalGroups: functionalGroups,
		segmentTerminator: segmentTerminator}
	return &interchange
}

func NewFunctionalGroup(segments []Segment) *FunctionalGroup {
	header, trailer := segments[0], segments[len(segments)-1]
	if header.Id != "GS" {
		panic("Expected GS as first segment, not " + header.Id)
	} else if trailer.Id != "GE" {
		panic("Expected GE as last segment, not " + trailer.Id)
	}

	transactionSets := GetTransactionSets(segments[1 : len(segments)-1])
	functionalGroup := FunctionalGroup{Header: header, Trailer: trailer, TransactionSets: transactionSets}
	return &functionalGroup
}

func NewTransactionSet(segments []Segment) *TransactionSet {
	header, trailer := segments[0], segments[len(segments)-1]
	if header.Id != "ST" {
		panic("Expected ST as first segment, not " + header.Id)
	} else if trailer.Id != "SE" {
		panic("Expected SE as last segment, not " + trailer.Id)
	}

	segments = segments[1 : len(segments)-2]
	transactionSet := TransactionSet{Header: header, Trailer: trailer, Segments: segments}
	return &transactionSet
}

func NewSegment(line string) *Segment {
	id := GetSegmentId(line, '*')
	elements := GetElements(line, "*")
	segment := Segment{Id: id, Elements: elements}
	return &segment
}

func NewElement(id, value string) *Element {
	element := Element{Id: id, Value: value}
	return &element
}

/*
   PARSER
*/

func GetSegments(file string, segmentTerminator byte) []Segment {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	var segments []Segment

	for bytes, err := r.ReadBytes(segmentTerminator); err == nil; bytes, err = r.ReadBytes(segmentTerminator) {
		segments = append(segments, *NewSegment(string(bytes[:len(bytes)-1])))
	}
	return segments
}

func GetElements(segment string, elementSeparator string) []Element {
	var elements []Element
	tokens := strings.Split(segment, elementSeparator)
	for _, token := range tokens[1:] {
		elements = append(elements, *NewElement("", token))
	}
	return elements
}

func GetFunctionalGroups(segments []Segment) []FunctionalGroup {
	var functionalGroups []FunctionalGroup

	var segmentArr []Segment
	for _, segment := range segments {
		segmentArr = append(segmentArr, segment)
		if segment.Id == "GE" {
			functionalGroups = append(functionalGroups, *NewFunctionalGroup(segmentArr))
			segmentArr = []Segment{}
		}
	}

	if len(segmentArr) > 0 {
		panic(fmt.Sprintf("Found %d segments outside of functional group.", len(segmentArr)))
	}

	return functionalGroups
}

func GetTransactionSets(segments []Segment) []TransactionSet {
	var transactionSets []TransactionSet

	var arr []Segment
	for _, segment := range segments {
		arr = append(arr, segment)
		if segment.Id == "SE" {
			transactionSets = append(transactionSets, *NewTransactionSet(arr))
			arr = []Segment{}
		}
	}

	if len(arr) > 0 {
		panic(fmt.Sprintf("Found %d segments outside of transactionSet.", len(arr)))
	}

	return transactionSets
}

func GetSegmentId(segment string, elementSeparator byte) string {
	return segment[:strings.IndexByte(segment, elementSeparator)]
}

func (s Segment) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func (element Element) String() string {
	return fmt.Sprintf("Element{Id: \"%s\", Value: \"%s\"}", element.Id, element.Value)
}

func GetSegmentTerminator(file string) byte {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for b, err := r.ReadByte(); err == nil; b, err = r.ReadByte() {
		if b == '>' {
			b, err = r.ReadByte()
			if err == nil {
				return b
			} else {
				panic("could not identify segment terminator")
			}
		}
	}
	panic("could not identify segment terminator")
}

/*
   WEB API
*/

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnEdi810(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: edi810")
	edi := NewInterchange("data/sample810.txt")
	json.NewEncoder(w).Encode(edi)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/edi810", returnEdi810)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
	edi := NewInterchange("data/sample810.txt")
	j, _ := json.Marshal(edi)
	fmt.Println(string(j))
}
