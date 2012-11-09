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

func doit(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<html><body>")
    pos := readWords("data/LoughranMcDonald_Positive.csv")
    neg := readWords("data/LoughranMcDonald_Negative.csv")

    b := readUrl("http://www.nytimes.com/services/xml/rss/nyt/HomePage.xml")

    v := new(RSS)
    xml.Unmarshal(b, &v)

    reg, _ := regexp.Compile("[<>,\\.]")
    pos_count := 0
    neg_count := 0
    allwords := make([]string, 0, 1000)
    for _, details := range v.Channel[0].Item {
        s := getBodyTextFromGlass(details.Link)
       
        s2 := reg.ReplaceAll([]byte(s), []byte(" "))
        words := strings.Split(string(s2), " ")

        for _, word := range words {
            if pos[strings.ToUpper(word)] == true {
                pos_count++
                allwords = append(allwords, word)
            }
            if neg[strings.ToUpper(word)] == true {
                neg_count++
                allwords = append(allwords, word)
            }
        }
    }
    
    fmt.Fprintf(w, "<span style='font-size: %dpx'>☺</span>", pos_count)
    fmt.Fprintf(w, "<span style='font-size: %dpx'>☹</span>", neg_count)

    /*
    for wordy := range allwords {
        fmt.Fprintf(w, "%s<br>\n", wordy)
    }
    */


    fmt.Fprintf(w, "</body></html>")

}

func main() {
    http.HandleFunc("/", doit)
    http.ListenAndServe(":8080", nil)
}
