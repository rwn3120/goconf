package conf

import (
    "testing"
    "errors"
    "time"
)

type TestConfiguration1 struct {
    Number   int64
    Float    float64
    Text     string
    Duration time.Duration
}

func (tc *TestConfiguration1) Validate() []error {
    return nil
}

type TestConfiguration2 struct {
}

func (tc *TestConfiguration2) Validate() []error {
    return []error{}
}

type TestConfiguration3 struct {
}

func (tc *TestConfiguration3) Validate() []error {
    return []error{errors.New("err1")}
}

func TestValid1(t *testing.T) {
    println(t.Name(), "... running")
    c := &TestConfiguration1{}
    if !IsValid(c) {
        t.Error("Not valid")
    }
}

func TestValid2(t *testing.T) {
    println(t.Name(), "... running")
    c := &TestConfiguration2{}
    if !IsValid(c) {
        t.Error("Not valid")
    }
}

func TestNonValid3(t *testing.T) {
    println(t.Name(), "... running")
    c := &TestConfiguration3{}
    if IsValid(c) {
        t.Error("Valid")
    }
}

func TestPrintYaml(t *testing.T) {
    println(t.Name(), "... running")
    c1 := &TestConfiguration1{
        Number:   123,
        Float:    0.123,
        Text:     "abc",
        Duration: time.Minute,
    }
    PrintYaml(c1)
}

func TestPrintToml(t *testing.T) {
    println(t.Name(), "... running")
    c1 := &TestConfiguration1{
        Number:   123,
        Float:    0.123,
        Text:     "abc",
        Duration: time.Minute,
    }
    PrintToml(c1)
}

func TestSaveLoad(t *testing.T) {
    println(t.Name(), "... running")
    c1 := &TestConfiguration1{
        Number:   123,
        Float:    0.123,
        Text:     "abc",
        Duration: time.Minute,
    }
    
    if err := SaveYaml("/tmp/configuration.test", c1); err != nil {
        t.Error(err.Error())
    }

    c2 := &TestConfiguration1{}
    if err := LoadYaml("/tmp/configuration.test", c2); err != nil {
        t.Error(err.Error())
    }
    if c2.Number != 123 {
        t.Error("unexpected number:", c2.Number)
    }
    if c2.Float != 0.123 {
        t.Error("unexpected number:", c2.Float)
    }
    if c2.Text != "abc" {
        t.Error("unexpected text:", c2.Text)
    }
    if c2.Duration != time.Minute {
        t.Error("unexpected duration:", c2.Duration)
    }
}
