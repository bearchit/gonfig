package gonfig

//go:generate mockgen -source=gonfig.go -package=mocks -destination=./mocks/gonfig.go

import (
	"io/ioutil"

	"github.com/mitchellh/mapstructure"

	"gopkg.in/yaml.v2"

	"github.com/kelseyhightower/envconfig"
)

type Engine struct {
	scanners []Scanner
}

func (e *Engine) AddScanner(scanners ...Scanner) {
	e.scanners = append(e.scanners, scanners...)
}

func New(options ...func(*Engine)) *Engine {
	e := new(Engine)
	for _, option := range options {
		option(e)
	}

	return e
}

func WithScanners(scanners ...Scanner) func(*Engine) {
	return func(engine *Engine) {
		engine.AddScanner(scanners...)
	}
}

func (e Engine) Unmarshal(v interface{}) error {
	for _, scanner := range e.scanners {
		if err := scanner.Struct(v); err != nil {
			if !scanner.SkipOnFail() {
				return err
			}
		}
	}
	return nil
}

type Scanner interface {
	Struct(v interface{}) error
	SkipOnFail() bool
}

type ymlScanner struct {
	filePath   string
	skipOnFail bool
}

func (s ymlScanner) Struct(v interface{}) error {
	fb, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if s.skipOnFail {
			return nil
		}
		return err
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(fb, m); err != nil {
		return err
	}
	return mapstructure.Decode(m, v)
}

func (s ymlScanner) SkipOnFail() bool {
	return s.skipOnFail
}

func NewYMLScanner(
	filePath string,
	skipOnFail bool,
) Scanner {
	return &ymlScanner{
		filePath:   filePath,
		skipOnFail: skipOnFail,
	}
}

type envScanner struct {
	prefix     string
	skipOnFail bool
}

func (s envScanner) Struct(v interface{}) error {
	return envconfig.Process(s.prefix, v)
}

func (s envScanner) SkipOnFail() bool {
	return s.skipOnFail
}

func NewEnvScanner(
	prefix string,
	skipOnFail bool,
) Scanner {
	return &envScanner{
		prefix:     prefix,
		skipOnFail: skipOnFail,
	}
}
