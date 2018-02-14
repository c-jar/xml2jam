package main

import (
  "fmt"
  "flag"
  "os"
  "strings"
)

var filenameIn, filenameOut string
func init(){
  flag.StringVar(&filenameIn, "in", "", "xml or mxl file name to convert.")
  flag.StringVar(&filenameOut, "out", "", "out jam music file")
  flag.Parse()
}

func main() {

  splitFilename := strings.Split(filenameIn, ".")
  if len(splitFilename) < 2 {
    fmt.Println("Please make an extension.")
    return
  }

  if splitFilename[len(splitFilename) - 1] == "mxl" {
    f, err := getXML(filenameIn, os.TempDir())
    if err != nil {
      panic(err)
    }
    filenameIn = f
  }

  xmlFile, err := os.Open(filenameIn)
	if err != nil {
		panic(err)
		return
  }
  defer xmlFile.Close()

  outFile, err := os.Create(filenameOut)
  if err != nil {
		panic(err)
		return
  }
  defer outFile.Close()

  score, err := readMusicXML(xmlFile)
  // fmt.Println(score)
  out, err := convertScoreToJam(score)

  if err != nil{
    panic(err)
  }
  _, err = outFile.WriteString(out)
  if err != nil{
    panic(err)
  }

  outFile.Sync()
  fmt.Println(filenameOut, "created successfuly !")
}
