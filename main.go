package main

import (
	"fmt"
	color "github.com/fatih/color"
	//pretty "github.com/kr/pretty"
	goyaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

func main() {
	file1 := os.Args[1]
	file2 := os.Args[2]

	buf, err := ioutil.ReadFile(file1)
	if err != nil {
		panic(err)
	}
	m1 := make(map[interface{}]interface{})
	err = goyaml.Unmarshal([]byte(buf), &m1)
	if err != nil {
		panic(err)
	}

	buf2, err := ioutil.ReadFile(file2)
	if err != nil {
		panic(err)
	}
	m2 := make(map[interface{}]interface{})
	err = goyaml.Unmarshal([]byte(buf2), &m2)
	if err != nil {
		panic(err)
	}

	m3 := make(map[interface{}]interface{})
	addLeftParam(m3, m1)
	addRightParam(m3, m2)

	ommitSameParam(m3)
	//pretty.Printf("%# v\n\n", m3)

	printDiffYaml(m3, "")
}

type Comparer struct {
	leftValue  interface{}
	rightValue interface{}
}

type Empty struct {
}

func addLeftParam(mergedMap, leftMap map[interface{}]interface{}) {
	for key, value := range leftMap {
		if value == nil || reflect.TypeOf(value).Kind() != reflect.Map {
			mergedMap[key] = value
		} else {
			if _, ok := mergedMap[key]; !ok {
				mergedMap[key] = make(map[interface{}]interface{})
			}
			addLeftParam(mergedMap[key].(map[interface{}]interface{}), leftMap[key].(map[interface{}]interface{}))
		}
	}
}

func addRightParam(mergedMap, rightMap map[interface{}]interface{}) {
	for key, value := range rightMap {
		if value == nil || reflect.TypeOf(value).Kind() != reflect.Map {
			if _, ok := mergedMap[key]; !ok {
				mergedMap[key] = Comparer{Empty{}, value}
			} else if mergedMap[key] == nil {
				mergedMap[key] = Comparer{nil, value}
			} else {
				mergedMap[key] = Comparer{mergedMap[key], value}
			}
		} else {
			if _, ok := mergedMap[key]; !ok {
				mergedMap[key] = make(map[interface{}]interface{})
			}
			if mergedMap[key] == nil || reflect.TypeOf(mergedMap[key]).Kind() != reflect.Map {
				mergedMap[key] = Comparer{mergedMap[key], value}
			} else {
				addRightParam(mergedMap[key].(map[interface{}]interface{}), rightMap[key].(map[interface{}]interface{}))
			}
		}
	}

	setComparer(mergedMap)
}

func setComparer(mergedMap map[interface{}]interface{}) {
	for key, value := range mergedMap {
		if value == nil {
			mergedMap[key] = Comparer{nil, Empty{}}
		} else if reflect.TypeOf(value).Kind() == reflect.Map {
			setComparer(mergedMap[key].(map[interface{}]interface{}))
		} else {
			vt := reflect.TypeOf(value).Kind()
			if vt != reflect.Map && vt != reflect.Struct {
				mergedMap[key] = Comparer{value, Empty{}}
			}
		}
	}
}

func ommitSameParam(m map[interface{}]interface{}) {
	for key, value := range m {
		if reflect.TypeOf(value).Kind() != reflect.Map {
			tmp := value.(Comparer)
			if reflect.DeepEqual(tmp.leftValue, tmp.rightValue) {
				delete(m, key)
			}
		} else {
			ommitSameParam(m[key].(map[interface{}]interface{}))
			if len(m[key].(map[interface{}]interface{})) == 0 {
				delete(m, key)
			}
		}
	}
}

func printDiffYaml(m map[interface{}]interface{}, indent string) {
	for key, value := range m {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			println(indent + key.(string) + ":")
			printDiffYaml(m[key].(map[interface{}]interface{}), indent+"    ")
		} else {
			tmp := value.(Comparer)
			println(indent + key.(string) + ":")
			color.Set(color.FgRed)
			printValue(tmp.leftValue, indent+"  - ")
			color.Set(color.FgGreen)
			printValue(tmp.rightValue, indent+"  + ")
			color.Unset()
		}

	}
}

func printValue(value interface{}, indent string) {
	if value != nil {
		if reflect.TypeOf(value).Kind() < reflect.Array {
			fmt.Printf("%s%#v\n", indent, reflect.ValueOf(value).Interface())
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			fmt.Printf("%s[", indent)
			for i, v := range value.([]interface{}) {
				if i != 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%#v", reflect.ValueOf(v).Interface())
			}
			println("]")
		case reflect.String:
			fmt.Printf("%s%#v\n", indent, reflect.ValueOf(value).Interface())
		case reflect.Map:
			for key, v := range value.(map[interface{}]interface{}) {
				println(indent + key.(string) + ":")
				printValue(v, "    "+indent)
			}
		}

	} else {
		println(indent + "null")
	}
}
