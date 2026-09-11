package main

import (
	"bufio"
	"os/exec"
	"sort"
	"strings"
)

// Package es un paquete de los repos pacman.
type Package struct {
	Name      string
	Repo      string
	Version   string
	Desc      string
	Installed bool
}

// score puntúa la relevancia contra una query en minúsculas.
// Mayor es mejor; 0 significa que no matchea.
func (p *Package) score(q string) int {
	if q == "" {
		return 1
	}
	name := strings.ToLower(p.Name)
	switch {
	case name == q:
		return 1000
	case strings.HasPrefix(name, q):
		// Los nombres cortos ganan: "vim" antes que "vim-spell-de".
		return 500 - len(name)
	case strings.Contains(name, q):
		return 200 - len(name)
	case strings.Contains(strings.ToLower(p.Desc), q):
		return 50
	}
	return 0
}

// LoadPackages lee todos los paquetes disponibles via `pacman -Ss .`.
// En el uConsole son ~13k paquetes aarch64 y tarda <1s.
//
// Formato de salida:
//
//	core/acl 2.4.0-1 [installed]
//	    Access control list utilities, libraries and headers
func LoadPackages() ([]*Package, error) {
	out, err := exec.Command("pacman", "-Ss", ".").Output()
	if err != nil {
		return nil, err
	}

	var pkgs []*Package
	var cur *Package

	sc := bufio.NewScanner(strings.NewReader(string(out)))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		// Línea indentada => descripción del paquete anterior.
		if line[0] == ' ' || line[0] == '\t' {
			if cur != nil {
				cur.Desc = strings.TrimSpace(line)
			}
			continue
		}

		p := parseHeader(line)
		if p == nil {
			continue
		}
		pkgs = append(pkgs, p)
		cur = p
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	logf("loaded %d packages", len(pkgs))
	return pkgs, nil
}

// parseHeader parsea "repo/nombre version [installed]".
func parseHeader(line string) *Package {
	slash := strings.IndexByte(line, '/')
	if slash <= 0 {
		return nil
	}
	repo := line[:slash]
	rest := line[slash+1:]

	fields := strings.Fields(rest)
	if len(fields) < 2 {
		return nil
	}

	return &Package{
		Name:      fields[0],
		Repo:      repo,
		Version:   fields[1],
		Installed: strings.Contains(rest, "[installed"),
	}
}

// Search filtra y ordena por relevancia. limit acota el resultado
// para no construir miles de filas GTK.
func Search(pkgs []*Package, query string, limit int) []*Package {
	q := strings.ToLower(strings.TrimSpace(query))

	type scored struct {
		pkg *Package
		s   int
	}
	var hits []scored
	for _, p := range pkgs {
		if s := p.score(q); s > 0 {
			hits = append(hits, scored{p, s})
		}
	}

	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].s != hits[j].s {
			return hits[i].s > hits[j].s
		}
		return hits[i].pkg.Name < hits[j].pkg.Name
	})

	if limit > 0 && len(hits) > limit {
		hits = hits[:limit]
	}

	out := make([]*Package, len(hits))
	for i, h := range hits {
		out[i] = h.pkg
	}
	return out
}
