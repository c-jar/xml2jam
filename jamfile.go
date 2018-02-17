package main

import (
  "fmt"
)

func convertScoreToJam(score Score) (string, error) {
  var strOut string

  fmt.Println("Add Title in comment ...")
  strOut += fmt.Sprintf("; Title: %s\n", score.GetTitle())

  fmt.Println("Add Identification in comment ...")
  for _, c := range score.Identification.Creators {
    strOut += fmt.Sprintf("; %s: %s\n", c.Type, c.Name)
  }
  if score.Identification.Rights != "" {
    strOut += fmt.Sprintf("; Rights: %s\n", score.Identification.Rights)
  }
  if score.Identification.Source != "" {
    strOut += fmt.Sprintf("; Source: %s\n", score.Identification.Source)
  }

  var selectScorePart ScorePartType
  if len(score.PartList.ScoreParts) > 1 {
    // TODO implement selection by user
    fmt.Println("There are", len(score.PartList.ScoreParts), "parts, select first.")
    selectScorePart = score.PartList.ScoreParts[0]
  }else {
    selectScorePart = score.PartList.ScoreParts[0]
  }
  fmt.Println("Use Part: {", selectScorePart, "}")

  var selectPart *PartType
  selectPart = nil
  for _, p := range score.Parts {
    if p.ID == selectScorePart.ID {
      selectPart = &p
    }
  }
  if selectPart == nil {
    return "", fmt.Errorf("Part with id == %s not found", selectScorePart.ID)
  }

  measures := selectPart.Measures
  var division, beatType float64
  repeatMeasure := -1
  repeatNumber := uint8(1)
  ignoreMeasure := false
  for i := 0; i < len(measures); i++ {
    m := measures[i]
    if i == 0 {
      tempo := m.Sound.Tempo
      if tempo == 0 {
        //return "", fmt.Errorf("tempo == 0")
        fmt.Println("Tempo == 0 => use 120")
        tempo = 120
      }
      fmt.Println("Add TEMPO", tempo)
      strOut += fmt.Sprintf("TEMPO %d\n", tempo)
      division = float64(m.Attributes.Division)
      if division == 0 {
        return "", fmt.Errorf("division == 0")
      }
      fmt.Println("Division is", division)
      beatType = float64(m.Attributes.TimeBeatType)
      if division == 0 {
        return "", fmt.Errorf("division == 0")
      }
      fmt.Println("Time beat type is", beatType)
    }
    fmt.Println("Read measure ", i)
    strOut += fmt.Sprintf("; Measure %d\n", i)

    // Bar line left
    for b, barline := range m.Barlines {
      if barline.Location != "left"{
        continue
      }
      if barline.Repeat.Direction == "forward" {
        repeatMeasure = i
        if barline.Repeat.Pass == 0 {
          fmt.Println("Barline repeat forward")
          strOut += "; Barline repeat forward\n"
          m.Barlines[b].Repeat.Pass ++
        }
      }
      if barline.Ending.Type == "start"{
        if repeatNumber != barline.Ending.Number {
          ignoreMeasure = true
        }
      }
    }

    // Notes
    if ! ignoreMeasure {
      for _, n := range m.Notes {
        if n.Voice != 1 {
          //fmt.Println("Ignore Voice :", n.Voice)
          continue
        }
        if n.Chord.Local != "" {
          //fmt.Println("Ignore Chord")
          continue
        }
        if n.Rest.Local != ""{
          strOut += "PAUSE"
        } else {
          strOut += fmt.Sprintf("%s%d", applyAlter(n.Pitch.Step, n.Pitch.Alter), n.Pitch.Octave)
        }
        strOut += fmt.Sprintf(" %.2f\n",
          (float64(n.Duration) / division))
      }
    }

    // Bar line right
    for b, barline := range m.Barlines {
      if barline.Location != "right"{
        continue
      }
      if barline.Repeat.Direction == "backward" && ! ignoreMeasure {
        if repeatMeasure == -1 {
          return "", fmt.Errorf("Read barline backward but I'm not read bar line forward")
        }
        if barline.Repeat.Pass == 0 {
          fmt.Println("Barline repeat backward to", repeatMeasure, "(", repeatNumber, ")")
          strOut += fmt.Sprintf("; Barline repeat backward to %d\n", repeatMeasure)
          i = repeatMeasure - 1
          repeatNumber ++
          repeatMeasure = -1
        }
        m.Barlines[b].Repeat.Pass ++
      }
      if barline.Ending.Type == "stop" {
        ignoreMeasure = false
      }
    }
  }

  return strOut, nil
}

func applyAlter(step string, alter int8) string {
  switch alter {
  case 1:
    switch step {
    case "E":
      return "F"
    case "B":
      return "C"
    }
    return step + "S"
  case -1:
    switch step {
    case "C":
      return "B"
    case "F":
      return "E"
    }
    step = string(byte(step[0]) - 1)
    if step < "A" {
      step = "G"
    }
    return step + "S"
  }
  return step
}
