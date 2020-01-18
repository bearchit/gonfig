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
			if scanner.BreakOnError() {
				return err
			}
		}
	}
	return nil
}

type Scanner interface {
	Struct(v interface{}) error
	BreakOnError() bool
}

type ymlScanner struct {
	filePath     string
	breakOnError bool
}

func (s ymlScanner) Struct(v interface{}) error {
	fb, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(fb, m); err != nil {
		return err
	}
	return mapstructure.Decode(m, v)
}

func (s ymlScanner) BreakOnError() bool {
	return s.breakOnError
}

func NewYMLScanner(
	filePath string,
	breakOnError bool,
) Scanner {
	return &ymlScanner{
		filePath:     filePath,
		breakOnError: breakOnError,
	}
}

type envScanner struct {
	prefix       string
	breakOnError bool
}

func (s envScanner) Struct(v interface{}) error {
	return envconfig.Process(s.prefix, v)
}

func (s envScanner) BreakOnError() bool {
	return s.breakOnError
}

func NewEnvScanner(
	prefix string,
	breakOnError bool,
) Scanner {
	return &envScanner{
		prefix:       prefix,
		breakOnError: breakOnError,
	}
}
