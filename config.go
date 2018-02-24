// Copyright (c) 2018 Mobentum Labs, LLC.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type (

	//Config interface to provices access methods
	Config interface {
		String(string) (string, error)
		Bool(string) (bool, error)
		Int(string) (int, error)
		Float(string) (float64, error)
		Map(string) (map[string]interface{}, error)
		List(string) ([]interface{}, error)

		MustString(string, ...string) string
		MustBool(string, ...bool) bool
		MustInt(string, ...int) int
		MustFloat(string, ...float64) float64
		MustMap(string, ...map[string]interface{}) map[string]interface{}
		MustList(string, ...[]interface{}) []interface{}

		Extend(Config) (Config, error)
	}

	//ConfigImpl struct to hold configuration data
	ConfigImpl struct {
		root map[string]interface{}
	}
)

// Get returns a value for the dotted path.
func (c *ConfigImpl) Get(path string) (interface{}, error) {
	return fetchValue(c.root, path)
}

//Extend shallow merge the with other config data
func (c *ConfigImpl) Extend(cfg Config) (Config, error) {
	if cfg != nil {
		for k, v := range cfg.(*ConfigImpl).root {
			c.root[k] = v
		}
	}
	return c, nil
}

//String returns a string value for the dotted path
func (c *ConfigImpl) String(path string) (string, error) {
	x, err := c.Get(path)
	if err != nil {
		return "", err
	}
	switch x.(type) {
	case string:
		return x.(string), nil
	}
	return "", fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustString(path string, defaults ...string) string {
	s, err := c.String(path)
	if err == nil {
		return s
	}
	for _, v := range defaults {
		return v
	}
	return ""
}

//Bool returns the bool value for the dotted path
func (c *ConfigImpl) Bool(path string) (bool, error) {
	x, err := c.Get(path)
	if err != nil {
		return false, err
	}
	switch x.(type) {
	case bool:
		return x.(bool), nil
	}
	return false, fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustBool(path string, defaults ...bool) bool {
	b, err := c.Bool(path)
	if err == nil {
		return b
	}
	for _, v := range defaults {
		return v
	}
	return false
}

//Int returns the int value for the dotted path
func (c *ConfigImpl) Int(path string) (int, error) {
	x, err := c.Get(path)
	if err != nil {
		return -1, err
	}
	switch x.(type) {
	case float64:
		return int(x.(float64)), nil
	}
	return -1, fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustInt(path string, defaults ...int) int {
	i, err := c.Int(path)
	if err == nil {
		return i
	}
	for _, v := range defaults {
		return v
	}
	return -1
}

//Float returns the float value for the dotted path
func (c *ConfigImpl) Float(path string) (float64, error) {
	x, err := c.Get(path)
	if err != nil {
		return -1, err
	}
	switch x.(type) {
	case float64:
		return x.(float64), nil
	}
	return -1, fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustFloat(path string, defaults ...float64) float64 {
	i, err := c.Float(path)
	if err == nil {
		return i
	}
	for _, def := range defaults {
		return def
	}
	return -1
}

func (c *ConfigImpl) Map(path string) (map[string]interface{}, error) {
	x, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	switch x.(type) {
	case map[string]interface{}:
		return x.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustMap(path string, defaults ...map[string]interface{}) map[string]interface{} {
	val, err := c.Map(path)
	if err == nil {
		return val
	}
	for _, def := range defaults {
		return def
	}
	return map[string]interface{}{}
}

func (c *ConfigImpl) List(path string) ([]interface{}, error) {
	x, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	switch x.(type) {
	case []interface{}:
		return x.([]interface{}), nil
	}
	return nil, fmt.Errorf("config: Unknown type at %q", path)
}

func (c *ConfigImpl) MustList(path string, defaults ...[]interface{}) []interface{} {
	val, err := c.List(path)
	if err == nil {
		return val
	}
	for _, def := range defaults {
		return def
	}
	return make([]interface{}, 0)
}

//Fetch

func fetchValue(cfg interface{}, path string) (interface{}, error) {
	parts := strings.Split(strings.TrimSpace(path), ".")
	for pos, part := range parts {
		if len(strings.TrimSpace(part)) == 0 {
			continue
		}
		curPath := strings.Join(parts[0:pos+1], ".")
		switch c := cfg.(type) {
		case []interface{}:
			if ix, error := strconv.ParseInt(part, 10, 0); error == nil {
				if int(ix) < len(c) {
					cfg = c[ix]
				} else {
					return nil, fmt.Errorf("config: Index out of bound at %q", curPath)
				}
			} else {
				return nil, fmt.Errorf("config: Unknown type at %q", curPath)
			}
		case map[string]interface{}:
			if value, ok := c[part]; ok {
				cfg = value
			} else {
				return nil, fmt.Errorf("config: Unknown path at %q", curPath)
			}
		default:
			return nil, fmt.Errorf("config: Unknown type at %q", curPath)
		}
	}
	return cfg, nil
}

//JSON

func parseJSON(data []byte) (Config, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &ConfigImpl{root: out}, nil
}

func ParseJSON(data string) (Config, error) {
	return parseJSON([]byte(data))
}

func ParseJSONFile(path string) (Config, error) {
	cb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseJSON(cb)
}
