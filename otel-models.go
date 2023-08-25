package main

type OtelSemanticModel struct {
	Groups []struct {
		ID         string `yaml:"id"`
		Prefix     string `yaml:"prefix"`
		Type       string `yaml:"type"`
		Brief      string `yaml:"brief"`
		SpanKind   string `yaml:"span_kind,omitempty"`
		Attributes []struct {
			ID               string      `yaml:"id"`
			Tag              string      `yaml:"tag,omitempty"`
			Brief            string      `yaml:"brief,omitempty"`
			RequirementLevel interface{} `yaml:"requirement_level,omitempty"`
			Type             interface{} `yaml:"type,omitempty"`
			Examples         interface{} `yaml:"examples,omitempty"`
			Note             string      `yaml:"note,omitempty"`
			Ref              string      `yaml:"ref,omitempty"`
		} `yaml:"attributes,omitempty"`
		Constraints []struct {
			AnyOf []string `yaml:"any_of,omitempty"`
		} `yaml:"constraints,omitempty"`
		Extends string `yaml:"extends,omitempty"`
	} `yaml:"groups"`
}
