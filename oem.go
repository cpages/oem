package main

import (
    "bufio"
    "fmt"
    "errors"
    "io/ioutil"
    "log"
    "net/mail"
    "os"
)

type Signant struct {
    nom string
    cp string
    fitxer string
}

func procMail(fname string) (string, Signant, error) {
    var s Signant
    file, err := os.Open(fname)
    if err != nil {
        return "", s, err
    }
    defer file.Close()

    m, err := mail.ReadMessage(file)
    if err != nil {
        return "", s, err
    }
    from := m.Header.Get("From")
    if from != "Olot es mou <form@olotesmou.cat>" {
        return "", s, errors.New("not from form")
    }

    scanner := bufio.NewScanner(m.Body)
    for scanner.Scan() {
        if scanner.Text() == "--- Email ---" {
            scanner.Scan()
            scanner.Scan()
            fmt.Println(scanner.Text())
        }
    }

    return fname, s, nil
}

func main() {
    dir := os.Args[1]
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    db := make(map[string]Signant)
    for _, f := range files[:10] {
        //fmt.Println(f.Name())
        correu, signant, err := procMail(fmt.Sprintf("%v/%v", os.Args[1], f.Name()))
        if err != nil {
            fmt.Printf("Error processant %v: %v\n", f.Name(), err.Error())
            continue
        }
        db[correu] = signant
    }

    fmt.Printf("File count: %v\n", len(files))
    fmt.Printf("DB count: %v\n", len(db))
}
