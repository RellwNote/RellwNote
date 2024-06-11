package theme

import (
	"errors"
	"os"
	"rellwnote/core/config"
	"rellwnote/core/files"
	"strings"
)

type Theme struct {
	Name      string
	Code      string
	ColorSets []string
}

func LoadAll() []Theme {
	var res []Theme
	for _, v := range strings.Split(config.Themes, ",") {
		v = strings.TrimSpace(v)
		t, err := Load(v)
		if err != nil {
			continue
		}
		res = append(res, t)
	}

	return res
}

func Load(name string) (t Theme, err error) {
	paths := []string{
		files.LibraryPath(config.ThemeDir, name+".css"),
		files.ProgramPath(config.ThemeDir, name+".css"),
	}
	for _, path := range paths {
		if files.IsFile(path) {
			read, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			t.Name = name
			t.Code = string(read)
			t.ColorSets = parseColorSets(t.Code)
			return t, nil
		}
	}
	return t, errors.New("theme " + name + " not found")
}

func parseColorSets(code string) []string {
	var res []string
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) <= 4 {
			continue
		}

		if strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") {

		}
		content := strings.TrimSpace(line[2 : len(line)-2])
		if !strings.HasPrefix(content, "ColorSets ") {
			continue
		}
		colors := strings.Split(content[10:], ",")
		for _, c := range colors {
			c = strings.TrimSpace(c)
			if len(c) > 0 {
				res = append(res, c)
			}
		}
	}
	return res
}
