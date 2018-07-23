package main

import ("log"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
)

type Configuration interface {
    Validate() *[]string
}

func IsValid(configuration Configuration) bool {
    if configuration == nil {
        return false
    }
    return configuration.Validate() == nil
}

func Validate(configurations ...Configuration) *[]string {
    var allErrors []string
    for _, configuration := range configurations {
        if configuration == nil {
            err := "Configuration can't be nil"
            allErrors = append(allErrors, err)
        } else {
            errors := configuration.Validate()
            if errors != nil {
                allErrors = append(allErrors, *errors...)
            }
        }
    }
    if len(allErrors) > 0 {
        return &allErrors
    }
    return nil
}

func Check(configurations ...Configuration) {
    if len(configurations) == 0 {
        panic("Nothing to check")
    }

    var allErrors []string
    for _, configuration := range configurations {
        errors := Validate(configuration)
        if errors != nil {
            allErrors = append(allErrors, *errors...)
        }

    }
    Handle(allErrors...)
}

func Handle(errors ...string) {
    if len(errors) > 0 {
        var message = "Configuration Errors:"
        for i := 0; i < len(errors); i++ {
            message = message + "\n\t" + errors[i]
        }
        log.Panic(message)
    }
}

func Load(file string, configuration Configuration) error {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        return err
    }

    return yaml.Unmarshal(content, configuration)
}

func LoadAndCheck(file string, configuration Configuration) Configuration {
    if configuration == nil {
        panic("Configuration can't be nil")
    }

    err := Load(file, configuration)
    if err != nil {
        panic(err)
    }
    Check(configuration)
    return configuration
}

func Save(file string, configuration Configuration) error {
    content, err := yaml.Marshal(configuration)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(file, content, 0644)
}

func Print(configuration Configuration) {
    content, err := yaml.Marshal(configuration)
    if err != nil {
        fmt.Printf("Configuration is not valid: %s", err.Error())
        return
    }
    fmt.Printf("%s\n", string(content[:]))
}