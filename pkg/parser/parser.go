package parser

import (
	"bytes"
	"log"
	"regexp"
	"strconv"

	service "github.com/geoirb/face-search/pkg/face-search"
)

type Parser struct {
	profileReg *regexp.Regexp
}

func New(layout string) (p *Parser, err error) {
	p = &Parser{}
	p.profileReg, err = regexp.Compile(layout)
	return
}

// GetProfileList from payload.
func (p *Parser) GetProfileList(payload []byte) []service.Profile {
	str, err := strconv.Unquote(string(payload))
	if err != nil {
		log.Fatalln(err)
	}
	payload = bytes.ReplaceAll([]byte(str), []byte("\\"), []byte(""))

	match := p.profileReg.FindAllSubmatch(payload, -1)
	profiles := make([]service.Profile, 0, len(match))
	for _, submatch := range match {
		if len(submatch) == 5 {
			profiles = append(
				profiles,
				service.Profile{
					FullName:    string(submatch[1]),
					Confidence:  string(submatch[2]),
					LinkProfile: string(submatch[3]),
					LinkPhoto:   string(submatch[4]),
				},
			)
		}
	}
	return profiles
}
