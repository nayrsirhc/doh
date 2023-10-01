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
            *dnsRecords = append(*dnsRecords, DNSRecord{
                record.Name,
                value,
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
