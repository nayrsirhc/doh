package doh

import (
	"fmt"
	"testing"
    "time"
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
                body, err:= DOHRequest(value, "example.com", "a")
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
}

func TestDecodeResponse(t *testing.T) {
    data := []byte{
        123,34,83,116,97,116,117,115,34,58,48,44,34,84,67,34,58,102,97,108,115,101,44,34,82,68,
        34,58,116,114,117,101,44,34,82,65,34,58,116,114,117,101,44,34,65,68,34,58,116,114,117,
        101,44,34,67,68,34,58,102,97,108,115,101,44,34,81,117,101,115,116,105,111,110,34,58,91,
        123,34,110,97,109,101,34,58,34,101,120,97,109,112,108,101,46,99,111,109,34,44,34,116,
        121,112,101,34,58,49,53,125,93,44,34,65,110,115,119,101,114,34,58,91,123,34,110,97,109,
        101,34,58,34,101,120,97,109,112,108,101,46,99,111,109,34,44,34,116,121,112,101,34,58,49,
        53,44,34,84,84,76,34,58,56,54,51,53,49,44,34,100,97,116,97,34,58,34,48,32,46,34,125,93,125,
    }
    dnsQuery := DNSQuery{}
    dnsRecords := []DNSRecord{}
    err := decodeResponse(data, &dnsQuery, &dnsRecords)

    if len(dnsRecords) > 0 && err == nil {
        t.Log("Decoded Data Successfully!")
        t.Log("Decoded Values:", dnsRecords)
    } else {
        t.Errorf("Unable to decode response")
    }
}
