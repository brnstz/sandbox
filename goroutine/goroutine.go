package main

import (
    "fmt"
    "os"
    "http"
    "io/ioutil"
    //"bytes"
)

/*
func byteString(b []byte) (string) {
    pos := bytes.IndexByte(b, 0)

    if pos == -1 {
        pos = len(b)
    }

    return string(b[0:pos])
}
*/

func printNum(i int, c *http.Client, ch chan int) {
    url := fmt.Sprintf("http://www.reddit.com?count=%d", i)
    resp, _ := c.Get(url)
    defer resp.Body.Close()

    readBytes, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("Status: %s, Bytes: %d\n", resp.Status, len(readBytes))
    ch <- 1

   //fmt.Print(byteString(readBytes))
}

func main() {
    client := new(http.Client)
    ch := make(chan int)

    for i := 0; i <= 1000; i += 25 {
        //fmt.Printf("%d", i)
        go printNum(i, client, ch)
    }
    for j := 40; j > 0; j-- {
        <-ch
        fmt.Println(j)
    }

    os.Stdout.Sync()
}
