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
        return "", s, errors.New("no és del formulari")
    }

    email := "N/A"
    s.nom = "N/A"
    s.cp = "N/A"
    scanner := bufio.NewScanner(m.Body)
    for scanner.Scan() {
        switch scanner.Text() {
        case "--- Nom i cognoms ---":
            scanner.Scan()
            scanner.Scan()
            s.nom = scanner.Text()
        case "--- Email ---":
            scanner.Scan()
            scanner.Scan()
            email = scanner.Text()
        case "--- Codi Postal ---":
            scanner.Scan()
            scanner.Scan()
            s.cp = scanner.Text()
        }
    }
    if email == "N/A" {
        return "", s, errors.New("correu no disponible")
    }

    s.fitxer = fname

    return email, s, nil
}

func dumpDB(db map[string]Signant, fname string) {
    f, err := os.Create(fname)
    if err != nil {
        fmt.Println("Error creant %s: %s", fname, err.Error())
    }
    defer f.Close()
    for c, s := range db {
        f.WriteString(fmt.Sprintf("%s (%s), %s\n", s.nom, c, s.cp))
    }
}

func main() {
    dir := os.Args[1]
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    db := make(map[string]Signant)
    for _, f := range files {
        correu, signant, err := procMail(fmt.Sprintf("%v/%v", os.Args[1], f.Name()))
        if err != nil {
            fmt.Printf("Error processant %v: %v\n", f.Name(), err.Error())
            continue
        }

        if _, exists := db[correu]; exists {
            fmt.Printf("Entrada duplicada: %v\n", correu)
            continue
        }

        db[correu] = signant
    }

    fmt.Printf("File count: %v\n", len(files))
    fmt.Printf("DB count: %v\n", len(db))

    dumpDB(db, "signants.txt")
}
