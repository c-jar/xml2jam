package main

import (
  "fmt"
  "flag"
  "os"
)

var filenameIn, filenameOut string
func init(){
  flag.StringVar(&filenameIn, "in", "", "xml file name to convert.")
  flag.StringVar(&filenameOut, "out", "", "out jam music file")
  flag.Parse()
}

func main() {
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
  fmt.Println(score)
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
