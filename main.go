package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DNSRecord struct {
    Question []struct {
        Name string `json: "name"`
        Type int `json: "type"`
    } `json: "Question"`
    Answer []struct {
        Name string `json: "name"`
	    Type int `json: "type"`
	    TTL  int `json: "TTL"`
	    Data string `json: "data"`
    } `json: "Answer"`
}

func resolveDNSGoogle(recordName string) (record_name string, record_type string, record_ttl int, record_value string) {

    resolveQuery := "https://dns.google/resolve?name=" + recordName
	resp, err := http.Get(resolveQuery)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dnsRecord DNSRecord

	if err := json.Unmarshal(body, &dnsRecord); err != nil {
		log.Fatalln("Error Parsing JSON: ", err)
	}
    
    record_name = dnsRecord.Answer[0].Name
    if dnsRecord.Answer[0].Type == 1 {
        record_type = "A"
    }
    record_ttl = dnsRecord.Answer[0].TTL
    record_value = dnsRecord.Answer[0].Data
    
    return record_name, record_type, record_ttl, record_value
}

func main() {
	queryName := flag.String("n", "example.com", "The name of the record you wish to resolve")
	flag.Parse()
    
    fmt.Println(resolveDNSGoogle(*queryName))

}