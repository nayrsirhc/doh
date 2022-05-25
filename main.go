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
		Type int    `json: "type"`
	} `json: "Question"`
	Answer []struct {
		Name string `json: "name"`
		Type int    `json: "type"`
		TTL  int    `json: "TTL"`
		Data string `json: "data"`
	} `json: "Answer"`
}

func resolveDNSGoogle(recordName string, recordType string) (record_name string, record_type string, record_ttl int, record_value string) {

	var resolveQuery string

	if recordType == "Not Specified" {
	    resolveQuery = "https://dns.google/resolve?name=" + recordName
	} else {
		resolveQuery = "https://dns.google/resolve?name=" + recordName + "&type=" + recordType
	}

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
	switch dnsRecord.Answer[0].Type {
	case 1:
		record_type = "A"
	case 2:
		record_type = "NS"
	case 5:
		record_type = "CNAME"
	case 6:
		record_type = "SOA"
	case 12:
		record_type = "PTR"
	case 13:
		record_type = "HINFO"
	case 15:
		record_type = "MX"
	case 16:
		record_type = "TXT"
	case 17:
		record_type = "RP"
	case 18:
		record_type = "AFSDB"
	case 24:
		record_type = "SIG"
	case 25:
		record_type = "KEY"
	case 28:
		record_type = "AAAA"
	case 29:
		record_type = "LOC"
	case 33:
		record_type = "SRV"
	case 35:
		record_type = "NAPTR"
	case 36:
		record_type = "KX"
	case 37:
		record_type = "CERT"
	case 39:
		record_type = "DNAME"
	case 42:
		record_type = "APL"
	case 43:
		record_type = "DS"
	case 44:
		record_type = "SSHFP"
	case 45:
		record_type = "IPSECKEY"
	case 46:
		record_type = "RRSIG"
	case 47:
		record_type = "NSEC"
	case 48:
		record_type = "DNSKEY"
	case 49:
		record_type = "DHCID"
	case 50:
		record_type = "NSEC3"
	case 51:
		record_type = "NSEC3PARAM"
	case 52:
		record_type = "TLSA"
	case 53:
		record_type = "SMIMEA"
	case 55:
		record_type = "HIP"
	case 59:
		record_type = "CDS"
	case 60:
		record_type = "CDNSKEY"
	case 61:
		record_type = "OPENPGPKEY"
	case 62:
		record_type = "CSYNC"
	case 63:
		record_type = "ZONEMD"
	case 64:
		record_type = "SVCB"
	case 65:
		record_type = "HTTPS"
	case 108:
		record_type = "EUI48"
	case 109:
		record_type = "EUI64"
	case 249:
		record_type = "TKEY"
	case 250:
		record_type = "TSIG"
	case 256:
		record_type = "URI"
	case 257:
		record_type = "CAA"
	case 32768:
		record_type = "TA"
	case 32769:
		record_type = "DLV"
	}
	record_ttl = dnsRecord.Answer[0].TTL
	record_value = dnsRecord.Answer[0].Data

	return record_name, record_type, record_ttl, record_value
}

func main() {

	queryName := flag.String("n", "example.com", "The name of the record you wish to resolve")
	queryType := flag.String("t", "Not Specified", "DNS Record Type")
	flag.Parse()

	fmt.Println(resolveDNSGoogle(*queryName, *queryType))

}
