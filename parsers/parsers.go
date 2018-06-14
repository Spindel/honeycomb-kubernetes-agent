package parsers

import (
	"fmt"

	"github.com/honeycombio/honeycomb-kubernetes-agent/config"
)

type Parser interface {
	Parse(line string) (map[string]interface{}, error)
}

type ParserFactory interface {
	Init(options map[string]interface{}) error
	New() Parser
}

func NewParserFactory(config *config.ParserConfig) (ParserFactory, error) {
	var factory ParserFactory
	switch config.Name {
	case "json":
		factory = &JSONParserFactory{}
	case "nop":
		factory = &NoOpParserFactory{}
	case "nginx", "envoy", "nginx-ingress":
		factory = &NginxParserFactory{
			// Default log format depends on the parser name specified in
			// configuration
			parserName: config.Name,
		}
	case "glog":
		factory = &GlogParserFactory{}
	case "redis":
		factory = &RedisParserFactory{}
	case "keyval":
		factory = &KeyvalParserFactory{}
	case "audit":
		factory = &AuditParserFactory{}
	default:
		return nil, fmt.Errorf("Unknown parser type %s", config.Name)
	}
	err := factory.Init(config.Options)
	if err != nil {
		return nil, err
	}
	return factory, nil
}
