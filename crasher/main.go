package main

import (
    "github.com/sbinet/go-imap/go1/imap"
    "crypto/tls"
    "flag"
    "fmt"
    //"reflect"
    //"os"
)

func login(server string, user string, pw string) (c *imap.Client, err error) {
    c, err = imap.DialTLS(server, new(tls.Config))

    if (err != nil) {
        return nil, err
    }

    _, err = c.Login(user, pw)

    if (err != nil) {
        return nil, err
    }

    return c, nil
}

func fetchIds(c *imap.Client) (ids []uint32) {
    seq_set, err := imap.NewSeqSet("4500:*")

    if (err != nil) {
        panic(err)
    }

    cmd, err := imap.Wait(c.Fetch(seq_set, "FLAGS"))

    if (err != nil) {
        panic(err)
    }

    for _, row := range(cmd.Data) {

        v, ok := row.Fields[0].(uint32)
        if (ok) {
            ids = append(ids, v)
        }
    }

    return ids
}

func fetchMessage(c *imap.Client, id uint32) {
    seq_set, err := imap.NewSeqSet(fmt.Sprintf("%d", id))

    if err != nil {
        panic(err)
    }

    cmd, err := imap.Wait(c.Fetch(seq_set, "ALL"))

    if err != nil {
        panic(err)
    }

    res, err := cmd.Result(0)

    if err != nil {
        panic(err)
    }

    //fmt.Println(reflect.TypeOf(cmd.Data[0].Fields[2]))
    fmt.Println(cmd.Data[0].Fields[2])

    fmt.Println(res.MessageInfo())

}

func main() {
    var (
        server string
        user string
        pw string
    )

    flag.StringVar(&server, "server", "", "IMAP server hostname")
    flag.StringVar(&user, "user", "", "IMAP username")
    flag.StringVar(&pw, "pw", "", "IMAP pw")
    flag.Parse()

    c, err := login(server, user, pw)

    if (err != nil) {
        panic(err)
    }

    _, err = imap.Wait(c.Select("INBOX", true))

    if (err != nil) {
        panic(err)
    }

    ids := fetchIds(c)

    for _, id := range(ids) {
        fetchMessage(c, id)
    }

    /*
    res, err := cmd.Result(0)
    fmt.Println(res.String())
    */
     /*
    res, err := cmd.Result(0)
    fmt.Println(res.String())

    for row := range(res)
    */

}
