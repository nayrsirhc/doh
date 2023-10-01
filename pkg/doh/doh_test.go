package doh

import (
    // "errors"
    "testing"
    "time"
    "fmt"

    "github.com/stretchr/testify/assert"
)

func TestDOHRequest(t *testing.T) {

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

    for _,record := range dnsRecords {
        t.Log("Resolving", record, "Record")
        timer1 := time.NewTimer(4 * time.Second)
        google := make(chan []byte)
        cloudflare := make(chan []byte)
        quad9 := make(chan []byte)
        providers := map[chan []byte]string{
            cloudflare: "https://1.1.1.1/dns-query?name=",
            google: "https://dns.google/resolve?name=",
            quad9: "https://dns.quad9.net:5053/dns-query?name=",
        }
        for key, value := range providers {
            go func(key chan []byte, value string) {
                defer close(key)
                body, err:= DOHRequest(value, "example.com", record)
                if err != nil {
                    time.Sleep(3 * time.Second)
                    t.Errorf(fmt.Sprintf("Failed to decode: %v\n", err))
                }
                key <- body
            }(key, value)
        }

        var body []byte

        select {
        case x := <-google:
            body = x
        case y := <-cloudflare:
            body = y
        case z := <-quad9:
            body = z
        case <-timer1.C:
            t.Errorf("Request timed out")
        }
        t.Log("Received", len(body), "bytes")
        if body == nil {
            t.Errorf("Empty reponse")
        }
    }
    body, err := DOHRequest("htts://1.1.1.1/dns-query!name=", "i???", "?????")
    assert.Nil(t, body)
    assert.Error(t, err)
}

func TestMapRecords(t *testing.T) {
    recordMap := map[int]string{
        1: "A",
        2: "NS",
        5: "CNAME",
        6: "SOA",
        12: "PTR",
        13: "HINFO",
        15: "MX",
        16: "TXT",
        17: "RP",
        18: "AFSDB",
        24: "SIG",
        25: "KEY",
        28: "AAAA",
        29: "LOC",
        33: "SRV",
        35: "NAPTR",
        36: "KX",
        37: "CERT",
        39: "DNAME",
        42: "APL",
        43: "DS",
        44: "SSHFP",
        45: "IPSECKEY",
        46: "RRSIG",
        47: "NSEC",
        48: "DNSKEY",
        49: "DHCID",
        50: "NSEC3",
        51: "NSEC3PARAM",
        52: "TLSA",
        53: "SMIMEA",
        55: "HIP",
        59: "CDS",
        60: "CDNSKEY",
        61: "OPENPGPKEY",
        62: "CSYNC",
        63: "ZONEMD",
        64: "SVCB",
        65: "HTTPS",
        108: "EUI48",
        109: "EUI64",
        249: "TKEY",
        250: "TSIG",
        256: "URI",
        257: "CAA",
        258: "AVC",
        32768: "TA",
        32769: "DLV",
    }
    for key, value := range recordMap {
        assert.Equal(t, value, mapRecords(key))
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
    assert.Error(t, err)
}

func TestRunQuery(t *testing.T) {
    err := RunQuery("example.com", "a", false, false)
    if err != nil {
        t.Errorf("RunQuery did not return any DNS records")
    }
}

func TestQueryExtensive(t *testing.T) {
    err := QueryExtensive("example.com")
    if err != nil {
        t.Errorf("QueryExtensive did not return any DNS records")
    }
}

func TestQueryAll(t *testing.T) {
    err := QueryAll("example.com")
    if err != nil {
        t.Errorf("QueryAll did not return any DNS records")
    }
}
