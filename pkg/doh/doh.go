package doh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// DNSRecord type for storing unmarshaled JSON data
type DNSRecord struct {
	Question []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"Question"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

func LogError(err error) {
	if err != nil {
		time.Sleep(1 * time.Second)
		log.Fatalln(err)
	}
}
// DOHRequest Makes a DNS-over-HTTP request which takes different providers, eg. Google, Cloudflare
func DOHRequest(provider string, recordName string, recordType string) (body []byte) {
	var resolveQuery string

	if recordType == "Not Specified" {
		resolveQuery = provider + recordName
	} else {
		resolveQuery = provider + recordName + "&type=" + recordType
	}

	req, err := http.NewRequest("GET", resolveQuery, nil)
	if err != nil {
		LogError(err)
	}
	req.Header.Set("accept", "application/dns-json")
	//We Read the response body on the line below.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogError(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError(err)
	}

	return body
}

func valdateRecordType(recordType string) (rRecordType string) {
	if recordType != "Not Specified" {
		recordType = strings.ToUpper(recordType)
		switch recordType {
		case "A", "NS", "CNAME", "SOA", "PTR", "HINFO", "MX":
			rRecordType = recordType
		case "TXT", "RP", "AFSDB", "SIG", "KEY", "AAAA", "LOC":
			rRecordType = recordType
		case "SRV", "NAPTR", "KX", "CERT", "DNAME", "APL", "DS":
			rRecordType = recordType
		case "NSEC3", "NSEC3PARAM", "TLSA", "SMIMEA", "HIP", "CDS":
			rRecordType = recordType
		case "CDNSKEY", "OPENPGPKEY", "CSYNC", "ZONEMD", "SVCB", "HTTPS":
			rRecordType = recordType
		case "EUI48", "EUI64", "TKEY", "TSIG", "URI", "CAA", "TA", "DLV":
			rRecordType = recordType
		default:
			log.Fatalln("Unrecognized DNS Record Type")
		}
	} else {
		rRecordType = recordType
	}
	return rRecordType
}

func resolveGoogle(recordName string, recordType string, c chan []byte) {
	body := DOHRequest("https://dns.google/resolve?name=", recordName, recordType)
	c <- body
	close(c)
}

func resolveCloudflare(recordName string, recordType string, c chan []byte) {
	body := DOHRequest("https://1.1.1.1/dns-query?name=", recordName, recordType)
	c <- body
	close(c)
}

func resolveQuad9(recordName string, recordType string, c chan []byte) {
	body := DOHRequest("https://dns.quad9.net:5053/dns-query?name=", recordName, recordType)
	c <- body
	close(c)
}

func decodeResponse(body []byte) (recordName []string, recordType []string, recordTTL []int, recordValue []string) {

	var dnsRecord DNSRecord

	if err := json.Unmarshal(body, &dnsRecord); err != nil {
		log.Fatalln("Failed to decode: ", err)
	}

	if len(dnsRecord.Answer) > 0 {
		for _, record := range dnsRecord.Answer {
			recordName = append(recordName, record.Name)
			recordTTL = append(recordTTL, record.TTL)
			recordValue = append(recordValue, record.Data)
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
			recordType = append(recordType, value)
		}
	}

	return recordName, recordType, recordTTL, recordValue
}

func RunQuery(queryName, queryType string, extensive bool) {
	valdateRecordType(queryType)
	timer1 := time.NewTimer(4 * time.Second)
	google := make(chan []byte)
	cloudflare := make(chan []byte)
	quad9 := make(chan []byte)
	go resolveGoogle(queryName, queryType, google)
	go resolveCloudflare(queryName, queryType, cloudflare)
	go resolveQuad9(queryName, queryType, quad9)

	var body []byte

	select {
	case x := <-google:
		body = x
	case y := <-cloudflare:
		body = y
	case z := <-quad9:
		body = z
	case <-timer1.C:
		log.Fatalln("Request timed out")
	}

	names, types, ttls, values := decodeResponse(body)

	if extensive && len(names) > 0 {
		fmt.Printf("\n%s:\n\n", queryType)
	}

	for i := range names {
		fmt.Printf("%s\t%s\t%d\t%s\n",
			strings.ToLower(names[i]),
			strings.ToUpper(types[i]),
			ttls[i],
			values[i])
	}
}

func QueryExtensive(queryName string) {

		dnsRecords := []string{
			"SOA",
			"NS",
			"A",
			"AAAA",
			"CNAME",
			"MX",
			"SRV",
			"TXT",
			"PTR",
			"HINFO",
			"RP",
			"AFSDB",
			"SIG",
			"KEY",
			"LOC",
			"NAPTR",
			"KX",
			"CERT",
			"DNAME",
			"APL",
			"DS",
			"NSEC3",
			"NSEC3PARAM",
			"TLSA",
			"SMIMEA",
			"HIP",
			"CDS",
			"CDNSKEY",
			"OPENPGPKEY",
			"CSYNC",
			"ZONEMD",
			"SVCB",
			"HTTPS",
			"EUI48",
			"EUI64",
			"TKEY",
			"TSIG",
			"URI",
			"CAA",
			"TA",
			"DLV",
		}
		for _, record := range dnsRecords {
			RunQuery(queryName, record, true)
		}
}

func QueryAll(queryName string) {

		dnsRecords := []string{
			"SOA",
			"NS",
			"A",
			"AAAA",
			"CNAME",
			"MX",
			"SRV",
			"TXT",
		}
		for _, record := range dnsRecords {
			RunQuery(queryName, record, false)
		}
}