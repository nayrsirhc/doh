package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func resolveDNSGoogle(recordName string, recordType string) (record_name []string, record_type []string, record_ttl []int, record_value []string) {

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

	var dnsRecord DNSRecord

	if err := json.Unmarshal(body, &dnsRecord); err != nil {
		switch recordType {
		case "A", "NS", "CNAME", "SOA", "PTR", "HINFO", "MX":
			log.Fatalln("Error Parsing JSON: ", err)
		case "TXT", "RP", "AFSDB", "SIG", "KEY", "AAAA", "LOC":
			log.Fatalln("Error Parsing JSON: ", err)
		case "SRV", "NAPTR", "KX", "CERT", "DNAME", "APL", "DS":
			log.Fatalln("Error Parsing JSON: ", err)
		case "SSHFP", "IPSECKEY", "RRSIG", "NSEC", "DNSKEY", "DHCID":
			log.Fatalln("Error Parsing JSON: ", err)
		case "NSEC3", "NSEC3PARAM", "TLSA", "SMIMEA", "HIP", "CDS":
			log.Fatalln("Error Parsing JSON: ", err)
		case "CDNSKEY", "OPENPGPKEY", "CSYNC", "ZONEMD", "SVCB", "HTTPS":
			log.Fatalln("Error Parsing JSON: ", err)
		case "EUI48", "EUI64", "TKEY", "TSIG", "URI", "CAA", "TA", "DLV":
			log.Fatalln("Error Parsing JSON: ", err)
		default:
			log.Fatalln("Unrecognized DNS Record Type")
		}

	}

	if len(dnsRecord.Answer) > 0 {
		for _,record := range dnsRecord.Answer {
			record_name = append(record_name, record.Name)
			record_ttl = append(record_ttl, record.TTL)
			record_value = append(record_value, record.Data)
			var value string
			switch record.Type {
			case 1:
				value = "A"
			case 2:
				value = "NS"
			case 5:
				value = "CNAME"
			case 6:
				value = "SOA"
			case 12:
				value = "PTR"
			case 13:
				value = "HINFO"
			case 15:
				value = "MX"
			case 16:
				value = "TXT"
			case 17:
				value = "RP"
			case 18:
				value = "AFSDB"
			case 24:
				value = "SIG"
			case 25:
				value = "KEY"
			case 28:
				value = "AAAA"
			case 29:
				value = "LOC"
			case 33:
				value = "SRV"
			case 35:
				value = "NAPTR"
			case 36:
				value = "KX"
			case 37:
				value = "CERT"
			case 39:
				value = "DNAME"
			case 42:
				value = "APL"
			case 43:
				value = "DS"
			case 44:
				value = "SSHFP"
			case 45:
				value = "IPSECKEY"
			case 46:
				value = "RRSIG"
			case 47:
				value = "NSEC"
			case 48:
				value = "DNSKEY"
			case 49:
				value = "DHCID"
			case 50:
				value = "NSEC3"
			case 51:
				value = "NSEC3PARAM"
			case 52:
				value = "TLSA"
			case 53:
				value = "SMIMEA"
			case 55:
				value = "HIP"
			case 59:
				value = "CDS"
			case 60:
				value = "CDNSKEY"
			case 61:
				value = "OPENPGPKEY"
			case 62:
				value = "CSYNC"
			case 63:
				value = "ZONEMD"
			case 64:
				value = "SVCB"
			case 65:
				value = "HTTPS"
			case 108:
				value = "EUI48"
			case 109:
				value = "EUI64"
			case 249:
				value = "TKEY"
			case 250:
				value = "TSIG"
			case 256:
				value = "URI"
			case 257:
				value = "CAA"
			case 32768:
				value = "TA"
			case 32769:
				value = "DLV"
			}
			record_type = append(record_type, value)
		}
	} else {
		log.Fatalln("Record Response is empty, please check query")
	}

	return record_name, record_type, record_ttl, record_value
}

func main() {

	queryName := flag.String("n", "example.com", "The name of the record you wish to resolve")
	queryType := flag.String("t", "Not Specified", "DNS Record Type")
	flag.Parse()

	names,types,ttls,values := resolveDNSGoogle(*queryName, *queryType)

	for i := range names {
		fmt.Println(strings.ToLower(names[i]), strings.ToUpper(types[i]), ttls[i], values[i])
	}
	

}
