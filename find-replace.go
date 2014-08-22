package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        line = strings.Replace(line, "function(", "function (", -1)
        line = strings.Replace(line, "if(", "if (", -1)
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}

func handleFiles(path string, wg *sync.WaitGroup) {
    defer wg.Done()
    lines, err := readLines(path)

    if err != nil {
        log.Printf("Reading File: %s failed. Err: %s", path, err)
        return
    }
    if err := writeLines(lines, path); err != nil {
        log.Printf("Writing file: %s failed. Err: %s", path, err)
        return
    }
}

func walkpath(path string, fileOrDir os.FileInfo, err error) error {
    var wg sync.WaitGroup

    if !fileOrDir.IsDir() {
        wg.Add(1)
        go handleFiles(path, &wg)
    } else if path == ".git" {
        return filepath.SkipDir
    }

    wg.Wait()
    return nil
}

func main() {
    path := flag.String("path", ".", "Specify a path")
    flag.Parse()

    filepath.Walk(*path, walkpath)
}
