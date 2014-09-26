package project_config

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
)

// Nested struct representing the JSON content
type NaiveConf struct {
    Domain struct {
      Main string
      Redirects []string
    }
}

// Reading and parsing the JSON
func (conf *NaiveConf) FromJSON(file string) error {
    // Reading JSON file
    jsonFile, err := ioutil.ReadFile(file)
    if err != nil { panic(err) }

    // Umarshalling JSON into nested structs
    return json.Unmarshal(jsonFile, &conf)
}

func Parse(projectPath string) NaiveConf {
    jsonPath := fmt.Sprintf("%s/naive.conf", projectPath)

    //Setting up a struct where will place our data that we extract
    JSONStruct := NaiveConf{}

    //Extracting the JSON data into a golang struct
    err := JSONStruct.FromJSON(jsonPath)
    if err != nil { panic(err) }

    // Outputing some stuff for now
    fmt.Println("Main domain:")
    fmt.Println(JSONStruct.Domain.Main)
    fmt.Println("Domains to redirect:")
    for _,element := range JSONStruct.Domain.Redirects {
      fmt.Println(element)
    }

    // returning the config struct
    return JSONStruct
}
