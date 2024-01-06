package doh

import (
	// "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSOARecord = []byte(`
    {
      "Status": 0,
      "TC": false,
      "RD": true,
      "RA": true,
      "AD": true,
      "CD": false,
      "Question": [
        {
          "name": "example.com",
          "type": 6
        }
      ],
      "Answer": [
        {
          "name": "example.com",
          "type": 6,
          "TTL": 2286,
          "data": "ns.icann.org. noc.dns.icann.org. 2022091389 7200 3600 1209600 3600"
        }
      ]
    }
  `)

/* EXAMPLE HTTP SERVER
func TestGetFixedValue(t *testing.T) {
   server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       if r.URL.Path != "/fixedvalue" {
           t.Errorf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
       }
       if r.Header.Get("Accept") != "application/json" {
           t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
       }
       w.WriteHeader(http.StatusOK)
       w.Write([]byte(`{"value":"fixed"}`))
   }))
   defer server.Close()

   value, _ := GetFixedValue(server.URL)
   if value != "fixed" {
       t.Errorf("Expected 'fixed', got %s", value)
   }
}
*/

func TestDOHRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"Query must have -name- parameter."}`))
		} else if r.Header.Get("Accept") == "application/dns-json" {
			switch r.URL.Query().Get("type") {
			case "SOA":
				w.WriteHeader(http.StatusOK)
				w.Write(testSOARecord)
			default:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error": "Invalid type"}`))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Error": "Accept: application/dns-json header is required by external DNS over HTTPS Providers"}`))
		}
	}))

	defer server.Close()

	t.Log("Test Record")
	body, err := DOHRequest(fmt.Sprintf("%v/dns-query?name=", server.URL), "example.com", "SOA")
	if err != nil {
		t.Error(err)
	}
	if body == nil {
		t.Errorf("Empty reponse")
	}
	t.Log("Received", len(body), "bytes")
	assert.Equal(t, 438, len(body))
	t.Log("Testing Handling of a Bad Request")
	body, err = DOHRequest(fmt.Sprintf("%v/dns-query?ame=", server.URL), "example.com", "SOA")
	isNil := assert.Nil(t, body)
	if isNil {
		t.Log("Body returned is Nil")
	}
	errorReturned := assert.Error(t, err)
	if errorReturned {
		t.Log("Error returned: ", err)
	}
}

func TestMapRecords(t *testing.T) {
	recordMap := map[int]string{
		1:     "A",
		2:     "NS",
		5:     "CNAME",
		6:     "SOA",
		12:    "PTR",
		13:    "HINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		24:    "SIG",
		25:    "KEY",
		28:    "AAAA",
		29:    "LOC",
		33:    "SRV",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		39:    "DNAME",
		42:    "APL",
		43:    "DS",
		44:    "SSHFP",
		45:    "IPSECKEY",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		53:    "SMIMEA",
		55:    "HIP",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		62:    "CSYNC",
		63:    "ZONEMD",
		64:    "SVCB",
		65:    "HTTPS",
		108:   "EUI48",
		109:   "EUI64",
		249:   "TKEY",
		250:   "TSIG",
		256:   "URI",
		257:   "CAA",
		258:   "AVC",
		32768: "TA",
		32769: "DLV",
	}
	for key, value := range recordMap {
		t.Log("Testing mapping of", value, "to", key)
		mapRecordsEqual := assert.Equal(t, value, mapRecords(key))
		if mapRecordsEqual {
			t.Log("Mapping", value, "to", key, "is successful")
		}
	}
}

func TestDecodeResponse(t *testing.T) {
	// Test case 1: decode a valid JSON response for an A record
	body := []byte(`{
        "Answer": [
            {
                "name": "example.com",
                "type": 1,
                "TTL": 299,
                "data": "93.184.216.34"
            }
        ]
    }`)
	dnsQuery := DNSQuery{}
	dnsRecords := []DNSRecord{}
	err := decodeResponse(body, &dnsQuery, &dnsRecords)
	if err != nil {
		t.Errorf("decodeResponse failed with error: %v", err)
	}
	if len(dnsRecords) != 1 {
		t.Errorf("decodeResponse did not return expected number of DNS records")
	}
	if dnsRecords[0].Name != "example.com" || dnsRecords[0].Type != "A" || dnsRecords[0].TTL != 299 || dnsRecords[0].Data != "93.184.216.34" {
		t.Errorf("decodeResponse returned unexpected DNS record: %v", dnsRecords[0])
	}

	// Test case 2: decode an invalid JSON response
	body = []byte(`invalid json`)
	err = decodeResponse(body, &dnsQuery, &dnsRecords)
	decodeError := assert.Error(t, err)
	if decodeError {
		t.Log("Request body:", string(body), "failed as expected with error", err)
	}
}

// func TestRunQuery(t *testing.T) {
// 	err := RunQuery("example.com", "a", false, false)
// 	if err != nil {
// 		t.Errorf("RunQuery did not return any DNS records")
// 	}
// }
//
// func TestQueryExtensive(t *testing.T) {
// 	err := QueryExtensive("example.com")
// 	if err != nil {
// 		t.Errorf("QueryExtensive did not return any DNS records")
// 	}
// }
//
// func TestQueryAll(t *testing.T) {
// 	err := QueryAll("example.com")
// 	if err != nil {
// 		t.Errorf("QueryAll did not return any DNS records")
// 	}
// }
