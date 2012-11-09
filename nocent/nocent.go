package main

import (
    "fmt"
    "os"
    "encoding/csv"
    "encoding/json"
    "encoding/xml"
    "net/http"
    "strings"
    "io/ioutil"
    "regexp"
)

type RSS struct {
    XMLName xml.Name `xml:"rss"` 
    Channel []Channel `xml:"channel"`
}

type Channel struct {
    Item []Item `xml:"item"`
}

type Item struct {
    Link string `xml:"guid"`
    Description string `xml:"description"`
}

type Glass struct {
    Assets []Asset `json:"assets"`
}
type Asset struct {
    Body string `json:"body"`
}


func readWords(filename string) map [string] bool {
    retmap := map [string] bool {}

    fh, _ := os.Open(filename)
    defer fh.Close()

    csv_r := csv.NewReader(fh)
    recs, _ := csv_r.ReadAll()

    for _, row := range recs {
        retmap[row[0]] = true
    }

    return retmap
}

func getBodyTextFromGlass(url string) string {
    bytes := readUrl(fmt.Sprintf("http://glass-output.prd.use1.nytimes.com/glass/outputmanager/V1/readingList.json?url=%s", url))
    v := new(Glass)
    json.Unmarshal(bytes, &v)

    return v.Assets[0].Body
}

func readUrl(url string) []byte {
    http := new(http.Client)
    resp, _ := http.Get(url)

    bytes, _ := ioutil.ReadAll(resp.Body)

    return bytes
}

func main() {
    pos := readWords("data/LoughranMcDonald_Positive.csv")
    neg := readWords("data/LoughranMcDonald_Negative.csv")

    b := readUrl("http://www.nytimes.com/services/xml/rss/nyt/HomePage.xml")

    v := new(RSS)
    xml.Unmarshal(b, &v)

    reg, _ := regexp.Compile("[<>,\\.]")
    pos_count := 0
    neg_count := 0
    for _, details := range v.Channel[0].Item {
        s := getBodyTextFromGlass(details.Link)
       
        s2 := reg.ReplaceAll([]byte(s), []byte(" "))
        words := strings.Split(string(s2), " ")

        for _, word := range words {
            if pos[strings.ToUpper(word)] == true {
                fmt.Printf("POS: %s\n", word)
                pos_count++
            }
            if neg[strings.ToUpper(word)] == true {
                fmt.Printf("NEG: %s\n", word)
                neg_count++
            }
        }
    }
    fmt.Println(pos_count)
    fmt.Println(neg_count)

}
