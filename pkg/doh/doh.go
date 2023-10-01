package doh

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "time"
    "errors"
    "sync"
)

// DNSRecord type for storing unmarshaled JSON data
type Question []struct {
    Name string `json:"name"`
    Type int    `json:"type"`
}
type Answer []struct {
    Name string `json:"name"`
    Type int    `json:"type"`
    TTL  int    `json:"TTL"`
    Data string `json:"data"`
}

type DNSQuery struct {
    Question `json:"Question"`
    Answer    `json:"Answer"`
}

type DNSRecord struct {
    Name string `json:"name"`
    Type string `json:"type"`
    TTL  int    `json:"TTL"`
    Data string `json:"data"`
}

func mapRecords(recordId int) (recordName string) {
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

    return recordMap[recordId]
}

// DOHRequest Makes a DNS-over-HTTP request which takes different providers, eg. Google, Cloudflare
func DOHRequest(provider string, recordName string, recordType string) (body []byte, ouch error) {
    var resolveQuery string

    if recordType == "Not Specified" {
        resolveQuery = provider + recordName
    } else {
        resolveQuery = provider + recordName + "&type=" + recordType
    }

    req, err := http.NewRequest("GET", resolveQuery, nil)
    if err != nil {
        fail := fmt.Sprintf("The HTTP request failed with error %s", err)
        ouch = errors.New(fail)
        return nil, ouch
    }
    req.Header.Set("accept", "application/dns-json")
    //We Read the response body on the line below.
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fail := fmt.Sprintf("The HTTP request failed with error %s", err)
        ouch = errors.New(fail)
        return nil, ouch
    }
    body, err = io.ReadAll(resp.Body)
    if err != nil {
        fail := fmt.Sprintf("The HTTP request failed with error %s", err)
        ouch = errors.New(fail)
        return nil, ouch
    }
    // Check if body contains the words blocked or not
    if strings.Contains(string(body), "blocked") {
        fail := fmt.Sprintf("The reqwuest was blocked\n%v", string(body))
        ouch = errors.New(fail)
        return nil, ouch
    }
    return body, nil
}

func decodeResponse(body []byte, dnsQuery *DNSQuery, dnsRecords *[]DNSRecord) (err error) {

    if err := json.Unmarshal(body, &dnsQuery); err != nil {
        newErr := fmt.Sprintf("Failed to decode: %v\n", err)
        ouch := errors.New(newErr)
        return ouch
    }

    if len(dnsQuery.Answer) > 0 {
        for _, record := range dnsQuery.Answer {
            recordType := mapRecords(record.Type)
            *dnsRecords = append(*dnsRecords, DNSRecord{
                record.Name,
                recordType,
                record.TTL,
                record.Data,
            })
        }
    }
    return nil
}

func RunQuery(queryName, queryType string, extensive bool, json bool) (error) {
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
        go func(key chan []byte, value string) error {
            defer close(key)
            body, err:= DOHRequest(value, queryName, queryType)
            if err != nil {
                time.Sleep(3 * time.Second)
                return err
            }
            key <- body
            return nil
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
        log.Fatalln("Request timed out")
    }
    dnsQuery := DNSQuery{}
    dnsRecords := []DNSRecord{}
    err := decodeResponse(body, &dnsQuery, &dnsRecords)
    if err != nil {
        return err
    }
    if json {
        fmt.Println(string(body))
    } else {
        if extensive && len(dnsRecords) > 0 {
            fmt.Printf("\n%s:\n\n", queryType)
        }
        for i := range dnsRecords {
            fmt.Printf("%s\t%s\t%d\t%s\n",
            strings.ToLower(dnsRecords[i].Name),
            strings.ToUpper(dnsRecords[i].Type),
            dnsRecords[i].TTL,
            dnsRecords[i].Data)
        }
    }
    return nil
}

func QueryExtensive(queryName string) error {

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
    var wg sync.WaitGroup
    for _, record := range dnsRecords {
        wg.Add(1)
        go func(queryName string, record string) error {
            defer wg.Done()
            err := RunQuery(queryName, record, true, false)
            if err != nil {
                return err
            }
            return nil
        }(queryName, record)
    }
    wg.Wait()
    return nil
}

func QueryAll(queryName string) error {

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
    var wg sync.WaitGroup
    for _, record := range dnsRecords {
        wg.Add(1)
        go func(queryName string, record string) error {
            defer wg.Done()
            err := RunQuery(queryName, record, false, false)
            if err != nil {
                return err
            }
            return nil
        }(queryName, record)
    }
    wg.Wait()
    return nil
}
