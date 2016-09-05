package main

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"
	"encoding/json"
	"os"
)

func main() {
	target_file := "JSComapct.js"
	search_pattern := ""

	args_len := len(os.Args)

	if args_len > 1 {
		target_file = os.Args[1]
		if args_len > 2 {
			search_pattern = os.Args[2]
			last_char := search_pattern[len(search_pattern) - 1]
			if last_char != '/' && last_char != '\\' {
				search_pattern += "\\"
			}
		}
	}

	// find all the javascript files
	files, _ := filepath.Glob(search_pattern + "*js")
	fmt.Println(files)

	if len(files) <= 0 {
		fmt.Println("No JavaScript files were found!")
		return
	}

	// compact all the javascript files into one
	var buf []byte
	var err error
	var data = ""

	for _, f := range files  {
		buf, err = ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err)
		}
		data += string(buf) + "\n"
	}

	// create the request
	form := url.Values{}
	form.Add("js_code", data)
	form.Add("compilation_level", "SIMPLE_OPTIMIZATIONS")
	form.Add("output_format", "json")
	form.Add("output_info", "compiled_code")
	form.Add("output_info", "errors")

	req, err := http.NewRequest("POST", "https://closure-compiler.appspot.com/compile", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// do the request
	hc := http.Client{}
	resp, err := hc.Do(req)

	// read the response
	response_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var j interface{}
	json.Unmarshal(response_data, &j)

	result := j.(map[string]interface{})

	// check for errors
	if result["errors"] != nil {
		code_errors := result["errors"].([]interface{})
		if len(code_errors) > 0 {
			fmt.Println("Errors in code:")
			for i, e := range code_errors {
				e_map := e.(map[string]interface{})
				fmt.Println(i + 1, e_map["error"], "=>", e_map["line"])
			}
			return
		}
	}

	// save the compiled code to file
	err = ioutil.WriteFile(target_file, []byte(result["compiledCode"].(string)), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
