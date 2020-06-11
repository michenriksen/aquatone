package parsers

import (
	"bufio"
	"io"

	"mvdan.cc/xurls/v2"
)

type RegexParser struct{}

func NewRegexParser() *RegexParser {
	return &RegexParser{}
}

func (p *RegexParser) Parse(r io.Reader) ([]string, error) {
	var targets []string
	targetsFilter := make(map[string]struct{})

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		for _, target := range xurls.Relaxed().FindAllString(scanner.Text(), -1) {
			if _, found := targetsFilter[target]; found {
				continue
			}
			targets = append(targets, target)
			targetsFilter[target] = struct{}{}
		}
	}
	return targets, nil
}
