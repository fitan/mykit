package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func Gen(path string) (err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &m)
	if err != nil {
		return
	}
	spew.Dump(m)

	depth(m, []string{})
	fmt.Println(typeM)
	return nil
}

var typeM = make(map[string]string)

func depth(m map[interface{}]interface{}, prePath []string) {
	for k, v := range m {
		ks := k.(string)
		tmpPath := append(prePath, ks)
		prePathStr := strings.Join(tmpPath, "")
		switch vt := v.(type) {
		case string:
			typeM[prePathStr] = "string"
		case float64:
			typeM[prePathStr] = "float64"
		case int:
			typeM[prePathStr] = "int"
		case map[interface{}]interface{}:
			depth(vt, tmpPath)
		}
	}
}
