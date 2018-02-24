package config_test

import (
	"testing"

	"github.com/mobentum/config"
	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	cfg, err := config.ParseJSONFile("resources/config/default.conf")
	if err != nil {
		t.Error(err)
	}

	debug, _ := cfg.Bool("debug")
	assert.Equal(t, true, debug)
	assert.Equal(t, true, cfg.MustBool("debug", false))

	age, _ := cfg.Int("age")
	assert.Equal(t, 26, age)
	assert.Equal(t, 24, cfg.MustInt("age1", 24))

	name, _ := cfg.String("name")
	assert.Equal(t, "John", name)
	assert.Equal(t, "lambda", cfg.MustString("name1", "lambda"))

	height, _ := cfg.Float("height")
	assert.Equal(t, 5.10, height)
	assert.Equal(t, 5.11, cfg.MustFloat("height1", 5.11))

	hobbies, _ := cfg.List("hobbies")
	assert.Equal(t, []interface{}{"skateboard", "snowboard", "go", "music"}, hobbies)
	assert.Equal(t, []interface{}{"tennis", "videogames"}, cfg.MustList("hobbies2", []interface{}{"tennis", "videogames"}))

	clothes, _ := cfg.Map("clothes.pants")
	assert.Equal(t, map[string]interface{}{"waist": 32.0, "height": 32.0}, clothes)
	assert.Equal(t, map[string]interface{}{"size": "large"}, cfg.MustMap("clothes.pants1", map[string]interface{}{"size": "large"}))

	//Nested
	val, _ := cfg.String("nested.1.2.3.0.b")
	assert.Equal(t, "c", val)
}

func Test_ConfigExtend(t *testing.T) {
	dcfg, err := config.ParseJSONFile("resources/config/default.conf")
	if err != nil {
		t.Error(err)
	}
	pcfg, err := config.ParseJSONFile("resources/config/production.conf")
	if err != nil {
		t.Error(err)
	}
	ecfg, err := dcfg.Extend(pcfg)
	if err != nil {
		t.Error(err)
	}

	debug, _ := ecfg.Bool("debug")
	assert.Equal(t, false, debug)
	assert.Equal(t, false, ecfg.MustBool("debug", true))

	env, _ := ecfg.String("env")
	assert.Equal(t, "production", env)
	assert.Equal(t, "default", ecfg.MustString("env1", "default"))
}
