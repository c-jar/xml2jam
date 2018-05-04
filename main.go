package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type ChannelFileVoieChord struct {
  filename string
  voice int
  chord int
}

var filenameIn, filenameOut string
var listVoices, allVoices, listChords, allChords bool
var voice int

func init() {
	flag.StringVar(&filenameIn, "in", "", "xml or mxl file name to convert.")
	flag.StringVar(&filenameOut, "out", "", "out jam music file.")
	flag.BoolVar(&listVoices, "list-voices", false, "list voices on select part.")
	flag.IntVar(&voice, "voice", 1, "Select voice.")
	flag.BoolVar(&allVoices, "all-voices", false, "create output file for each voices.")
  flag.BoolVar(&listChords, "list-chords", false, "give max chord on selected voice.")
  flag.BoolVar(&allChords, "all-chords", false, "create output file for each chord.")
	flag.Parse()
}

func verifyParam() bool {
	if filenameIn == "" {
		fmt.Println("You must give an input file. Use -in.")
		flag.Usage()
		return false
	}
	return true
}

func exportJAM(score Score, filename string, voice, chord int) {
  outFile, err := os.Create(filename)
  if err != nil {
    panic(err)
    return
  }
  defer outFile.Close()

  out, err := convertScoreToJam(score, voice, chord)

  if err != nil {
    panic(err)
  }
  _, err = outFile.WriteString(out)
  if err != nil {
    panic(err)
  }

  outFile.Sync()
  fmt.Println(filename, "created successfuly !")
}

func createFilename(score Score, defaultVoice, defaultChord int, ch chan ChannelFileVoieChord) {
  ret := ChannelFileVoieChord{filenameOut, defaultVoice, defaultChord}
  if allVoices || allChords {
    p := score.Parts[0]
    voices := p.GetVoices()
    splitFilename := strings.Split(filenameOut, ".")
    for _, v := range voices {
      if !allVoices && v != defaultVoice {
        continue
      }
      ret.voice = v
      filename := fmt.Sprintf("%s.voice_%d", splitFilename[0], v)
      if allChords {
        cmax := p.GetChords(v)
        for i := 0; i < cmax; i++ {
          filename := fmt.Sprintf("%s.chord_%d", filename, i)
          for i, f := range splitFilename {
            if i != 0 {
              filename += "." + f
            }
          }
          ret.chord = i
          ret.filename = filename
          ch <- ret
        }
      } else {
        for i, f := range splitFilename {
          if i != 0 {
            filename += "." + f
          }
        }
        ret.filename = filename
        ch <- ret
      }

    }
  } else {
    ch <- ret
  }
  close(ch)
}

func main() {
	if !verifyParam() {
		return
	}

	splitFilename := strings.Split(filenameIn, ".")
	if len(splitFilename) < 2 {
		fmt.Println("Please make an extension.")
		return
	}

	if splitFilename[len(splitFilename)-1] == "mxl" {
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

	score, err := readMusicXML(xmlFile)
	fmt.Println(score)

	if listVoices {
		p := score.Parts[0]
		voices := p.GetVoices()
		fmt.Println("Voices :")
		for _, v := range voices {
			fmt.Println(" -", v)
		}
	}

  if listChords {
    fmt.Println("Max chord is ", score.Parts[0].GetChords(voice))
  }

	if filenameOut != "" {
    ch := make(chan ChannelFileVoieChord)
    go createFilename(score, voice, 0, ch)
    for c := range ch {
      exportJAM(score, c.filename, c.voice, c.chord)
    }
	}

}
