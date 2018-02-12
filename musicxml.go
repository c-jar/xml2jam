package main

import (
  "fmt"
  "encoding/xml"
  "os"
  "io/ioutil"
)

// Score represents score, root element.
type Score struct {
  MovementTitle string `xml:"movement-title"`
  Identification IdentificationType `xml:"identification"`
  PartList PartListType `xml:"part-list"`
  Parts []PartType `xml:"part"`
}

func (score Score) String() string {
  return fmt.Sprintf("Title: %s,\nIdentification: %s,\nPart-list : %s\nParts: %s",
    score.MovementTitle,
    score.Identification,
    score.PartList,
    score.Parts)
}

// IdentificationType represents score's informations.
type IdentificationType struct {
  Composer string `xml:"creator"`
	Rights   string `xml:"rights"`
	Source   string `xml:"source"`
	Title    string `xml:"movement-title"`
}

func (id IdentificationType) String() string {
  return fmt.Sprintf("Composer : %s, Rights: %s, Source: %s, Title :%s",
    id.Composer,
    id.Rights,
    id.Source,
    id.Title)
}

// PartListType represents list of part.
type PartListType struct {
  ScoreParts []ScorePartType `xml:"score-part"`
}

func (p PartListType) String() string {
  str := "[\n"
  for _, s := range p.ScoreParts {
    str += fmt.Sprintf("\t{%s}\n", s)
  }
  return str + "]"
}

// ScorePartType represents score's part.
type ScorePartType struct {
  ID string `xml:"id,attr"`
  PartName string `xml:"part-name"`
}

func (s ScorePartType) String() string {
    return fmt.Sprintf("Id: %s, name: %s", s.ID, s.PartName)
}

// PartType represents content part.
type PartType struct {
  ID string `xml:"id,attr"`
  Measures []MeasureType `xml:"measure"`
}

func (p PartType)String() string {
    return fmt.Sprintf("{\n\tId: %s\n\tMeasures: %s\n}",
      p.ID,
      p.Measures)
}

// MeasureType represents a measure in part.
type MeasureType struct {
  Number int `xml:"number,attr"`
  Attributes AttributesType `xml:"attributes"`
  Sound SoundType `xml:"sound"`
  Notes []NoteType `xml:"note"`
}

func (m MeasureType) String() string {
  str :=  fmt.Sprintf("{Number: %d,", m.Number)
  if m.Number == 1 {
    str += fmt.Sprintf("Attributes: {%s}, ", m.Attributes)
  }
  return str + fmt.Sprintf("Tempo: %d,\nNotes:\n%s\n}", m.Sound.Tempo, m.Notes)
}

// SoundType contains tempo.
type SoundType struct {
  Tempo int `xml:"tempo,attr"`
}

// AttributesType represents part's attributes.
type AttributesType struct {
  Division int `xml:"divisions"`
  KeyFifths int `xml:"key>fifths"`
  KeyMode string `xml:"key>mode"`
  TimeBeats int `xml:"time>beats"`
  TimeBeatType int `xml:"time>beat-type"`
  ClefSign string `xml:"clef>sign"`
  ClefLine int `xml:"clef>line"`
}

func (a AttributesType) String() string {
  return fmt.Sprintf(
    "Division: %d, Key fifths: %d, Key mode: %s, Beats: %d %d, Clef: %s %d",
    a.Division,
    a.KeyFifths,
    a.KeyMode,
    a.TimeBeats,
    a.TimeBeatType,
    a.ClefSign,
    a.ClefLine)
}

// NoteType represents a note in a measure
type NoteType struct {
	Pitch    PitchType    `xml:"pitch"`
	Duration int      `xml:"duration"`
	Voice    int      `xml:"voice"`
	Type     string   `xml:"type"`
	Rest     xml.Name `xml:"rest"`
	Chord    xml.Name `xml:"chord"`
}

func (n NoteType) String() string {
  return fmt.Sprintf("Pitch: {%s}, Duration: %d, Voice: %d, Type: %s, Rest: %s, Chord: %s\n",
    n.Pitch, n.Duration, n.Voice, n.Type, n.Rest, n.Chord)
}

// PitchType represents the pitch of a note
type PitchType struct {
	Alter      int8   `xml:"alter"`
	Step       string `xml:"step"`
	Octave     int    `xml:"octave"`
}

func (p PitchType) String() string {
  return fmt.Sprintf("Alter: %d, Step: %s, Octave: %d",
    p.Alter, p.Step, p.Octave)
}

// readMusicXML read xmlFile and create Score object.
func readMusicXML(xmlFile *os.File) (Score, error) {
  b, err := ioutil.ReadAll(xmlFile)
  if err != nil {
    panic(err)
  }
	var s Score
  xml.Unmarshal(b, &s)
  return s, nil
}
