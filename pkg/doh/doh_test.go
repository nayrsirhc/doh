package doh

import (
	"testing"
    "time"
)

func TestValidateRecord(t *testing.T) {
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
        recordType := valdateRecordType(record)
        t.Log("Validating Record Type", record)
        if recordType != record {
            t.Errorf("Failed to validate record string")
        }
        t.Log("Success")
    }
}

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
        body := DOHRequest("https://dns.google/resolve?name=", "exmaple.com", record)
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

    names, types, ttls, values, err := decodeResponse(data)

    if len(names) > 0 && err == nil {
        t.Log("Decoded Data Successfully!")
        t.Log("Decoded Values:",names,types,ttls,values)
    } else {
        t.Errorf("Unable to decode response")
    }
}

func TestResolveGoogle(t *testing.T) {

    timer1 := time.NewTimer(1 * time.Second)
    google := make(chan []byte)
    go resolveGoogle("google.co.za", "a", google)
    select {
    case x := <-google:
        t.Log("Received", len(x), "bytes")
    case <-timer1.C:
        t.Errorf("Response took longer than a second")
    }
}

func TestResolveCloudflare(t *testing.T) {

    timer1 := time.NewTimer(1 * time.Second)
    cloudflare := make(chan []byte)
    go resolveCloudflare("cloudflare.net", "a", cloudflare)
    select {
    case x := <-cloudflare:
        t.Log("Received", len(x), "bytes")
    case <-timer1.C:
        t.Errorf("Response took longer than a second")
    }
}

func TestResolveQuad9(t *testing.T) {

    timer1 := time.NewTimer(1 * time.Second)
    quad9 := make(chan []byte)
    go resolveQuad9("quad9.net", "a", quad9)
    select {
    case x := <-quad9:
        t.Log("Received", len(x), "bytes")
    case <-timer1.C:
        t.Errorf("Response took longer than a second")
    }
}
