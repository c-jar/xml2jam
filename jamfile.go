package main

import (
  "fmt"
)

func convertScoreToJam(score Score) (string, error) {
  var strOut string

  fmt.Println("Add Movment Title in comment ...")
  strOut += fmt.Sprintf("; Movement Title: %s\n", score.MovementTitle)

  fmt.Println("Add Identification in comment ...")
  if score.Identification.Composer != "" {
    strOut += fmt.Sprintf("; Composer: %s\n", score.Identification.Composer)
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
  for i, m := range measures {
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
    for _, n := range m.Notes {
      if n.Voice != 1 {
        fmt.Println("Ignore Voice :", n.Voice)
        continue
      }
      if n.Chord.Local != "" {
        fmt.Println("Ignore Chord")
        continue
      }
      if n.Rest.Local != ""{
        strOut += "PAUSE"
      } else {
        // TODO Accidental and armure.
        strOut += fmt.Sprintf("%s%d", applyAlter(n.Pitch.Step, n.Pitch.Alter), n.Pitch.Octave)
      }
      strOut += fmt.Sprintf(" %.2f\n",
        (float64(n.Duration) / division))
    }
  }

  return strOut, nil
}

func applyAlter(step string, alter int8) string {
  switch alter {
  case 1:
    return step + "S"
  case -1:
    step = string(byte(step[0]) - 1)
    if step < "A" {
      step = "G"
    }
    return step + "S"
  }
  return step
}
