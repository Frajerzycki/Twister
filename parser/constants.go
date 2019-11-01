package parser

import "regexp"

const KeyBase int = 16

var formatArgumentRegexp *regexp.Regexp = regexp.MustCompile("(-t)([io])([kd])")
