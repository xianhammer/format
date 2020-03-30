package oxmsg

import (
	"net"
	"strings"
)

type Block struct {
	lines []string
}

func (b *Block) Body() string {
	return TrimMultipleSpaces(strings.Join(b.lines, " "))
}

func (b *Block) Lines() []string {
	return b.lines
}

func (b *Block) Map() (m map[string][]string) {
	m = make(map[string][]string)
	for _, line := range b.lines {
		s := strings.SplitN(line, ":", 2)
		m[s[0]] = append(m[s[0]], s[1])
	}
	return
}

func (b *Block) ScanMail(callback func(key, mail string)) {
	for _, line := range b.lines {
		s := strings.SplitN(line, ":", 2)
		line = s[1]
		for i, l := 0, len(line); i < l; i++ {
			if line[i] != '@' {
				continue
			}

			start := i - 1
			c := line[start]
			for ; start >= 0 && ((c-'0') < 10 || (c&0xDF)-'A' < 27) && c != '.' && c != '-' && c != '_'; start-- {
				c = line[start]
			}
			start += 2 // Skip last found
			if start >= i {
				continue
			}

			dot := i + 1
			for ; dot < l && line[dot] != '.'; dot++ {
			}

			if dot-i < 2 {
				continue
			}

			end := dot + 1
			for ; end < l && (line[end]&0xDF)-'A' < 27; end++ {
			}

			if end-dot < 2 {
				continue
			}

			callback(s[0], string(line[start:end]))
		}
	}
}

func (b *Block) ScanIP(callback func(key string, ip net.IP)) {
	for _, line := range b.lines {
		s := strings.SplitN(line, ":", 2)
		line = s[1]
		for i, l := 0, len(line); i < l; i++ {
			c := line[i] //& 0xDF // Make character uppercase - has no effect on digits
			if c-'0' >= 10 && (c&0xDF)-'A' >= 6 {
				continue
			}

			start, part, separators := i, 0, 0
			for ; i < l; i++ {
				c := line[i]
				if c-'0' < 10 || (c&0xDF)-'A' < 6 {
					part++
					if part > 4 {
						break
					}
				} else if c == '.' || c == ':' {
					separators++
					part = 0
				} else {
					break
				}
			}

			if separators >= 3 {
				if ip := net.ParseIP(string(line[start:i])); ip != nil {
					callback(s[0], ip)
				}
			}
		}
	}
}

func (b *Block) add(line string, continued bool) {
	strings.TrimSpace(line)
	if continued {
		b.lines[len(b.lines)-1] += line
	} else {
		b.lines = append(b.lines, line)
	}
}
