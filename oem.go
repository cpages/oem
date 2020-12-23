package main

import (
    "bufio"
    "fmt"
    "errors"
    "io/ioutil"
    "log"
    "net/mail"
    "os"
    "sort"
    "time"
)

type Signant struct {
    nom string
    cp string
    fitxer string
    data time.Time
}

func procMail(path string, finfo os.FileInfo) (string, Signant, error) {
    fname := fmt.Sprintf("%v/%v", path, finfo.Name())
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

    s.fitxer = finfo.Name()
    s.data = finfo.ModTime()

    return email, s, nil
}

type Pair struct {
    key string
    value Signant
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].value.data.Before(p[j].value.data) }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func dumpDB(db map[string]Signant, fname string) {
    f, err := os.Create(fname)
    if err != nil {
        fmt.Println("Error creant %s: %s", fname, err.Error())
    }
    defer f.Close()

    pl := make(PairList, len(db))
    i := 0
    for k, v := range db {
      pl[i] = Pair{k, v}
      i++
    }
    sort.Sort(pl)

    for _, p := range pl {
        f.WriteString(fmt.Sprintf("%s (%s), %s\n", p.value.nom, p.key, p.value.cp))
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
        correu, signant, err := procMail(os.Args[1], f)
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
