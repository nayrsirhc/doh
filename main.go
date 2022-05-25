package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

	recordType = strings.ToUpper(recordType)

	switch recordType {
	case "A", "NS", "CNAME", "SOA", "PTR", "HINFO", "MX":
		record_type = recordType
	case "TXT", "RP", "AFSDB", "SIG", "KEY", "AAAA", "LOC":
		record_type = recordType
	case "SRV", "NAPTR", "KX", "CERT", "DNAME", "APL", "DS":
		record_type = recordType
	case "SSHFP", "IPSECKEY", "RRSIG", "NSEC", "DNSKEY", "DHCID":
		record_type = recordType
	case "NSEC3", "NSEC3PARAM", "TLSA", "SMIMEA", "HIP", "CDS":
		record_type = recordType
	case "CDNSKEY", "OPENPGPKEY", "CSYNC", "ZONEMD", "SVCB", "HTTPS":
		record_type = recordType
	case "EUI48", "EUI64", "TKEY", "TSIG", "URI", "CAA", "TA", "DLV":
		record_type = recordType
	default:
		log.Fatalln("Unrecognized DNS Record type")
		os.Exit(1)
	}

	var dnsRecord DNSRecord

	if err := json.Unmarshal(body, &dnsRecord); err != nil {
		log.Fatalln("Error Parsing JSON: ", err)
	}

	record_name = dnsRecord.Answer[0].Name
	record_ttl = dnsRecord.Answer[0].TTL
	record_value = dnsRecord.Answer[0].Data

	return strings.ToLower(record_name), record_type, record_ttl, record_value
}

func main() {

	queryName := flag.String("n", "example.com", "The name of the record you wish to resolve")
	queryType := flag.String("t", "Not Specified", "DNS Record Type")
	flag.Parse()

	fmt.Println(resolveDNSGoogle(*queryName, *queryType))

}
