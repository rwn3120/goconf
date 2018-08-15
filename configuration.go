package conf

import (
    "log"
    "gopkg.in/yaml.v2"
    "github.com/pelletier/go-toml"
    "io/ioutil"
    "fmt"
    "os"
    "github.com/pkg/errors"
)

type Configuration interface {
    Validate() []error
}

func IsValid(configuration Configuration) bool {
    if configuration == nil {
        return false
    }
    return configuration.Validate() == nil || len(configuration.Validate()) == 0
}

var tomlMarshal = toml.Marshal
var tomlUnMarshal = toml.Unmarshal
var yamlMarshal = yaml.Marshal
var yamlUnMarshal = yaml.Unmarshal

func Validate(configurations ...Configuration) []error {
    var allErrors []error
    for _, configuration := range configurations {
        if configuration == nil {
            allErrors = append(allErrors, errors.New("configuration can't be nil"))
        } else {
            errors := configuration.Validate()
            if errors != nil {
                allErrors = append(allErrors, errors...)
            }
        }
    }
    if len(allErrors) > 0 {
        return allErrors
    }
    return nil
}

func Check(configurations ...Configuration) {
    if len(configurations) == 0 {
        panic("Nothing to check")
    }

    var allErrors []error
    for _, configuration := range configurations {
        errors := Validate(configuration)
        if errors != nil {
            allErrors = append(allErrors, errors...)
        }

    }
    Handle(allErrors...)
}

func Handle(errors ...error) {
    if len(errors) > 0 {
        var message = "Configuration Errors:"
        for i := 0; i < len(errors); i++ {
            message = message + "\n\t" + errors[i].Error()
        }
        log.Panic(message)
    }
}

func save(file string, configuration Configuration, marshal func(interface{}) ([]byte, error)) error {
    content, err := marshal(configuration)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(file, content, 0644)
}

func load(file string, configuration Configuration, unmarshal func([]byte, interface{}) (error)) error {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        return err
    }

    return unmarshal(content, configuration)
}

func LoadYaml(file string, configuration Configuration) error {
    return load(file, configuration, yamlUnMarshal)
}

func SaveYaml(file string, configuration Configuration) error {
    return save(file, configuration, yamlMarshal)
}

func LoadToml(file string, configuration Configuration) error {
    return load(file, configuration, tomlUnMarshal)
}

func SaveToml(file string, configuration Configuration) error {
    return save(file, configuration, tomlMarshal)
}

func loadAndCheck(file string, configuration Configuration, load func(string, Configuration) error) Configuration {
    if configuration == nil {
        panic("Configuration can't be nil")
    }

    err := load(file, configuration)
    if err != nil {
        panic(err)
    }
    Check(configuration)
    return configuration
}

func LoadYamlAndCheck(file string, configuration Configuration) Configuration {
    return loadAndCheck(file, configuration, LoadYaml)
}

func LoadTomlAndCheck(file string, configuration Configuration) Configuration {
    return loadAndCheck(file, configuration, LoadYaml)
}

func print(configuration Configuration, to func(Configuration) (string, error)) {
    value, err := to(configuration)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
    } else {
        fmt.Fprintf(os.Stdout, "%s\n", value)
    }
}

func PrintYaml(configuration Configuration) {
    print(configuration, ToYaml)
}

func PrintToml(configuration Configuration) {
    print(configuration, ToToml)
}

func toString(configuration Configuration, marshal func(interface{}) ([]byte, error)) (string, error) {
    content, err := marshal(configuration)
    if err != nil {
        return "", err
    }
    return string(content), err
}

func ToYaml(configuration Configuration) (string, error) {
    return toString(configuration, yamlMarshal)
}

func ToToml(configuration Configuration) (string, error) {
    return toString(configuration, tomlMarshal)
}
