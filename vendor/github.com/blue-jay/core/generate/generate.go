// Package generate will parse and create files from template pairs.
// A template pair is a set of template files:
//   * default.json - json file
//   * default.gen - any type of text file
//
// Both files are template files and are parsed using the Go text/template
// package. The 'jay generate' tool loops through the first level of key pairs
// for empty strings. For every empty string, an argument is required to be
// passed (whether empty or not) to the 'jay generate' command.
//
// Let's look at generate/model/default.json:
// {
//   "config.type": "single",
//   "config.output": "model/{{.package}}/{{.package}}.go",
//   "package": "",
//   "table": ""
// }
//
// Let's break down this command into pieces:
// jay generate model/default package:automobile table:car
//
// Argument: 'model/default'
// Specifies generate/model/default.json and generate/model/default.gen are the
// template pair.
//
// Argument: 'package:automobile'
// The key, 'package', from default.json will be filled with the value:
// 'automobile'
//
// Argument: 'table:car'
// The key, 'table', from default.json will be filled with the value: 'car'
//
// The .json file is actually parsed up to 100 times (LoopLimit of 100 can be
// changed at the package level) to ensure all variables like '{{.package}}' are
// set to the correct value.
//
// In the first iteration of parsing, the 'package' key is set to 'car'.
// In the second iteration of parsing, the '{{.package}}' variables
// are set to 'car' also since the 'package' key becomes a variable.
//
// All first level keys (info, package, table) become variables after the first
// iteration of parsing so they can be used without the file. If a variable is
// misspelled and is never filled, a helpful error will be displayed.
//
// The 'output' key under 'info' is required. It should be the relative output
// file path to the project root for the generated file.
//
// The folder structure of the templates (model, controller, etc) has no effect
// on the generation, it's purely to aid with organization of the template pairs.
//
// You must store the path to the env.json file in the environment
// variable: JAYCONFIG. The file is at project root that is prepended to all
// relative file paths.
//
// Examples:
//   jay generate model/default package:car table:car
//   Generate a new model from variables in model/default.json and applies
//   those variables to model/default.gen.
//   jay generate controller/default package:car url:car model:car view:car
//   Generate a new controller from variables in controller/default.json
//   and applies those variables to controller/default.gen.
//
// Flags:
//   Argument 1 - model/default or controller/default
//   Relative path without an extension to the template pair. Any combination
//   of folders and files can be used.
//   Argument 2,3,etc - package:car
//   Key pair to set in the .json file. Required for every empty key in the
//   .json file.
package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/blue-jay/core/file"
)

var (
	// LoopLimit specified the max number of iterations of replacing variables
	// with values to prevent an infinite loop.
	LoopLimit = 100
)

// Container is an easy way to read only the generation info from the config
// file.
type Container struct {
	Generation Info `json:"Generation"`
}

// ParseJSON unmarshals bytes to structs
func (c *Container) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// Info holds the config.
type Info struct {
	// TemplateFolder is the folder path for the code generation templates
	TemplateFolder string `json:"TemplateFolder"`
}

// Run starts the template generation logic
func Run(args []string, projectFolder string, templateFolder string) error {
	// Ensure the template pair files exist
	jsonFilePath := filepath.Join(templateFolder, args[0]+".json")
	if !file.Exists(jsonFilePath) {
		return fmt.Errorf("File doesn't exist: %v", jsonFilePath)
	}

	argMap, err := argsToMap(args)
	if err != nil {
		return err
	}

	// Get the json file as a map - not parsed
	mapFile, err := jsonFileToMap(jsonFilePath)
	if err != nil {
		return err
	}

	// Generate variable map
	variableMap, err := generateVariableMap(mapFile, argMap)
	if err != nil {
		return err
	}

	// Check for config type
	configType, ok := variableMap["config.type"]
	if !ok {
		return errors.New("Key, 'config.type', is missing from the .json file")
	}

	// Handle based on config.type
	switch configType {
	case "single":
		// Template File
		genFilePath := filepath.Join(templateFolder, args[0]+".gen")

		// Generate the template
		return generateSingle(projectFolder, genFilePath, variableMap)
	case "collection":
		return generateCollection(projectFolder, templateFolder, variableMap)
	default:
		return fmt.Errorf("Value of '%v' for key 'config.type' is not supported", configType)
	}
}

func generateCollection(projectFolder string, templateFolder string, variableMap map[string]interface{}) error {
	// Check for required key
	collectionRaw, ok := variableMap["config.collection"]
	if !ok {
		return errors.New("Key, 'config.collection', is missing from the .json file")
	}

	collection, ok := collectionRaw.([]interface{})
	if !ok {
		return errors.New("Key, 'config.collection', is not in the correct format")
	}

	// Loop through the collections
	for i, v := range collection {
		vMap, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("Values for key, 'config.collection', are not in the correct format")
		}

		for name, varArray := range vMap {
			argMap, ok := varArray.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Item at index '%v' for key, 'config.collection', is not in the correct format", i)
			}

			// Template File
			genFilePath := filepath.Join(templateFolder, name+".gen")
			jsonFilePath := filepath.Join(templateFolder, name+".json")

			// Get the json file as a map - not parsed
			mapFile, err := jsonFileToMap(jsonFilePath)
			if err != nil {
				return err
			}

			// Generate variable map
			variableMap, err := generateVariableMap(mapFile, argMap)
			if err != nil {
				return err
			}

			// Check for config type
			configType, ok := variableMap["config.type"]
			if !ok {
				return errors.New("Key, 'config.type', is missing from the .json file")
			}

			// Handle based on config.type
			switch configType {
			case "single":
				err = generateSingle(projectFolder, genFilePath, variableMap)
				if err != nil {
					return err
				}
			case "collection":
				err = generateCollection(projectFolder, templateFolder, variableMap)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("Value of '%v' for key 'config.type' is not supported", configType)
			}
		}
	}

	return nil
}

func generateSingle(projectFolder string, genFilePath string, variableMap map[string]interface{}) error {
	// Check for required key
	if _, ok := variableMap["config.output"]; !ok {
		return errors.New("Key, 'config.output', is missing from the .json file")
	}

	// Output file
	outputRelativeFile := fmt.Sprintf("%v", variableMap["config.output"])
	outputFile := filepath.Join(projectFolder, outputRelativeFile)

	// Check if the file exists
	if file.Exists(outputFile) {
		return fmt.Errorf("Cannot generate because file already exists: %v", outputFile)
	}

	// Check if the folder exists
	dir := filepath.Dir(outputFile)
	if !file.Exists(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// If config.parse = false
	if val, ok := variableMap["config.parse"]; ok {
		if b, _ := strconv.ParseBool(fmt.Sprintf("%v", val)); !b {
			// Don't parse template, just copy to new file
			err := toFile(genFilePath, outputFile)
			if err != nil {
				return err
			}
			fmt.Println("Code generated:", outputFile)
			return nil
		}
	}

	// Parse template and write to file
	err := fromMapToFile(genFilePath, variableMap, outputFile)
	if err != nil {
		return err
	}

	fmt.Println("Code generated:", outputFile)

	return nil
}

func cloneMap(originalMap map[string]interface{}) map[string]interface{} {
	copyMap := make(map[string]interface{})

	for k, v := range originalMap {
		copyMap[k] = v
	}

	return copyMap
}

// jsonFileToMap converts json file to an interface map.
func jsonFileToMap(file string) (map[string]interface{}, error) {
	// Read the config file
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Convert json to interface
	var d map[string]interface{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// fromMapToFile will output a file by parsing a template and applying
// variables from an interface map.
func fromMapToFile(templateFile string, d map[string]interface{}, outputFile string) error {
	base := filepath.Base(templateFile)

	// Parse the template
	t, err := template.New(base).Funcs(funcmap()).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	// Create the output file
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Fills template with variables and writes to file
	//err = t.Execute(f, d)
	err = t.ExecuteTemplate(f, base, d)
	if err != nil {
		return err
	}

	return nil
}

// toFile will output a file by without parsing.
func toFile(templateFile string, outputFile string) error {

	data, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func argsToMap(args []string) (map[string]interface{}, error) {
	// Fill a new map with variables
	argMap := make(map[string]interface{})
	for _, a := range args[1:] {
		arr := strings.Split(a, ":")

		if len(arr) < 2 {
			return nil, fmt.Errorf("Arg is in wrong format: %v", a)
		}

		argMap[arr[0]] = strings.Join(arr[1:], ":")
	}

	return argMap, nil
}

func fillEmptyVariables(m, argMap map[string]interface{}) (map[string]interface{}, error) {
	// Loop through the map to fill empty variables
	for s, v := range m {
		switch t := v.(type) {
		case string:
			if len(t) == 0 {
				if val, ok := argMap[s]; ok {
					m[s] = val
				} else {
					return nil, fmt.Errorf("Variable missing: %v", s)
				}
			} else { // Else delete any values that are not empty
				delete(m, s)
			}
		default:
			delete(m, s)
		}
	}

	return m, nil
}

func funcmap() template.FuncMap {
	// Function map
	f := make(template.FuncMap)

	// Add the ability to use strings.title in the template
	f["title"] = func(s string) string {
		return strings.Title(s)
	}

	return f
}

func parseTemplate(m map[string]interface{}, mapFileBytes []byte) (map[string]interface{}, bool, error) {
	// Create the buffer
	buf := new(bytes.Buffer)

	// Parse the template
	t, err := template.New("").Funcs(funcmap()).Parse(string(mapFileBytes))
	if err != nil {
		return nil, true, err
	}

	// Fills template with variables
	err = t.Execute(buf, m)
	if err != nil {
		return nil, true, err
	}

	parsedTemplate := buf.Bytes()

	// Convert the json text back to a map
	err = json.Unmarshal(parsedTemplate, &m)
	if err != nil {
		return nil, true, err
	}

	breakNow := false

	// If the parsed template is completely filled, then stop the run, else
	// keep running
	if !strings.Contains(string(parsedTemplate), "<no value>") {
		breakNow = true
	}

	return m, breakNow, nil
}

// generateVariableMap returns the relative file output path and the map of
// variables.
func generateVariableMap(mapFile map[string]interface{}, argMap map[string]interface{}) (map[string]interface{}, error) {
	var err error

	// Clone map
	m := cloneMap(mapFile)

	// Fill empty variables of map
	m, err = fillEmptyVariables(m, argMap)
	if err != nil {
		return nil, err
	}

	// Look through the map of the file and update it with the variables
	for s := range mapFile {
		if passedVal, ok := m[s]; ok {
			mapFile[s] = passedVal
		}
	}

	// Look through the map of the variables and overwrite it with any left over arguments
	// This allows you to overwrite config.output
	for s, v := range argMap {
		if _, ok := m[s]; !ok {
			mapFile[s] = v
		}
	}

	// Counter to prevent infinite loops
	counter := 0

	for true {
		// Convert the mapFile to bytes
		mapFileBytes, err := json.Marshal(mapFile)
		if err != nil {
			return nil, err
		}

		// Parse template to determine if all the variables are passed and
		// then break
		var breakNow bool
		m, breakNow, err = parseTemplate(m, mapFileBytes)
		if breakNow || err != nil {
			break
		}

		var invalidKeys []string

		// Loop through the map to find empty variables
		for s, v := range m {
			switch t := v.(type) {
			case string:
				if strings.Contains(t, "<no value>") {
					invalidKeys = append(invalidKeys, s)
					delete(m, s)
				}
			default:
				// This if statement outputs a helpful error if in a nested
				// map
				if strings.Contains(fmt.Sprintf("%v", v), "<no value>") {
					invalidKeys = append(invalidKeys, fmt.Sprintf("%v %v", s, v))
				}
				delete(m, s)
			}
		}

		counter++

		if counter > LoopLimit {
			return nil, fmt.Errorf("Check these keys for variable mistakes: %v", invalidKeys)
		}
	}

	return m, err
}
