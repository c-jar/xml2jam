package main

import (
  "fmt"
  "archive/zip"
  "strings"
  "path"
  "os"
  "io"
)


func getXML(src, dirDst string) (string, error) {
  r, err := zip.OpenReader(src)
  if err != nil {
    return "", fmt.Errorf("Can't open %s", src)
  }
  defer r.Close()

  filepath := ""
  for _, f := range r.File {
    if strings.Contains(f.Name, "META-INF") {
      continue
    }
    splitString := strings.Split(f.Name, ".")
    if len(splitString) < 2 {
      continue
    }
    if splitString[len(splitString) - 1] != "xml" {
      continue
    }

    filepath = path.Join(dirDst, f.Name)
    file, err := os.OpenFile(filepath,
      os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
    if err != nil {
      return "", err
    }
    defer file.Close()

    rc, err := f.Open()
    if err != nil {
      return "", err
    }
    defer rc.Close()

    _, err = io.Copy(file, rc)
    if err != nil {
      return "", err
    }
    fmt.Println(filepath, "created")
  }

  return filepath, nil
}
