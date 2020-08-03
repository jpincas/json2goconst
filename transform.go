package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type JsonMap map[string]json.RawMessage

type structRep struct {
	name, defs string
}

func (sp structRep) String() string {
	return sp.defs
}

func transform(jsonData []byte) (string, error) {
	simpleStringMap := JsonMap{}
	if err := json.Unmarshal(jsonData, &simpleStringMap); err != nil {
		return "", err
	}

	structReps, err := simpleStringMap.walkFromRoot()
	if err != nil {
		return "", err
	}

	// Sort alphabetically by node key to aid testing
	sort.Slice(structReps, func(i, j int) bool {
		return structReps[i].name < structReps[j].name
	})

	var structDefs []string
	for _, sr := range structReps {
		structDefs = append(structDefs, sr.String())
	}

	return strings.Join(structDefs, "\n\n"), nil
}

func (jm JsonMap) walkFromRoot() ([]structRep, error) {
	return jm.walk([]string{}, []structRep{})
}

func createConstant(ss []string) string {
	var titleCased []string

	for _, s := range ss {
		s = strings.Title(s)
		s = strings.ReplaceAll(s, "-", "_")

		titleCased = append(titleCased, s)
	}

	return strings.Join(titleCased, "_")
}

func createJsPath(ss []string) string {
	return strings.Join(ss, ".")
}

func (jm JsonMap) walk(thisLevelName []string, children []structRep) ([]structRep, error) {
	var defLines []string

	// Sort alphabetically by node key to aid testing
	type Node struct {
		Key   string
		Value json.RawMessage
	}

	var nodes []Node
	for k, v := range jm {
		nodes = append(nodes, Node{k, v})
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Key < nodes[j].Key
	})

	for _, node := range nodes {
		// try decoding the contents as another struct
		nextLevel := JsonMap{}
		if err := json.Unmarshal(node.Value, &nextLevel); err == nil {
			children, err = nextLevel.walk(append(thisLevelName, node.Key), children)
			if err != nil {
				return []structRep{}, err
			}

		} else {
			var s string
			// an error here means the json is invalid for our
			// purposes
			if err := json.Unmarshal(node.Value, &s); err != nil {
				return []structRep{}, err
			}

			defLines = append(defLines, fmt.Sprintf(`	%s = "%s"`,
				createConstant(append(thisLevelName, node.Key)),
				createJsPath(append(thisLevelName, node.Key)),
			))
		}
	}

	if len(defLines) > 0 {
		thisStruct := structRep{
			name: createConstant(thisLevelName),
			defs: strings.Join(defLines, "\n"),
		}

		return append(children, thisStruct), nil
	}

	return children, nil

}
