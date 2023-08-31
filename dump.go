package gic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type dump struct {
	mu           sync.Mutex
	dir          string
	overrideRoot struct {
		root     string
		override string
	}
	initCount uint
	data      *data
	got       []*component // keeps components was got during component creation. needs to dump dependencies
}

type withDump struct{ d *dump }

func (w withDump) applyGlobalContainerOption(c *container) {
	c.dump = w.d
}

func WithDump(opts ...dumpOption) withDump {
	d := &dump{}
	for _, opt := range opts {
		opt.applyDumpOption(d)
	}

	if d.dir == "" {
		panic("empty dump directory. use ic.WithDumpDir")
	}

	d.data = &data{
		Files:      map[string][]string{},
		StageImpls: map[string][]*compCoordsJSON{},
	}

	return withDump{d: d}
}

type dumpOption interface{ applyDumpOption(*dump) }
type withDumpDir struct{ dir string }
type withOverrideRoot struct {
	root     string
	override string
}

func (w withDumpDir) applyDumpOption(d *dump) { d.dir = w.dir }

func (w withOverrideRoot) applyDumpOption(d *dump) {
	d.overrideRoot = struct {
		root     string
		override string
	}{root: w.root, override: w.override}
}

// WithDumpDir specifies directory where to put dump
func WithDumpDir(dir string) withDumpDir {
	return withDumpDir{dir: dir}
}

// WithOverrideRoot overrides source root
// Rare case. If binary was compiled on one machine
// and then run on another you may override first
// machines root with second machine root.
// Root is used to lookup source code files to show
// components source code definitions on dump page
func WithOverrideRoot(root, override string) withOverrideRoot {
	return withOverrideRoot{root: root, override: override}
}

type data struct {
	// Fills when gic.Init called
	Initialized []*compJSON                  `json:"initialized"`
	Files       map[string][]string          `json:"files"`
	Stages      []*stageJSON                 `json:"stages"`
	StageImpls  map[string][]*compCoordsJSON `json:"stage_impls"`
}

type compJSON struct {
	Order      uint              `json:"order"`
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	File       string            `json:"file"`
	LineStart  int               `json:"line_start"`
	LineEnd    int               `json:"line_end"`
	DirectDeps []*compCoordsJSON `json:"direct_deps"`
}

type compCoordsJSON struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type stageJSON struct {
	ID       string `json:"id"`
	Order    string `json:"order"`
	Parallel bool   `json:"parallel"`
}

func dumpComponent(c *container, comp *component) {
	if c.dump == nil {
		return
	}

	c.dump.mu.Lock()
	defer c.dump.mu.Unlock()

	rt := reflect.TypeOf(comp.c)

	cJSON := &compJSON{
		Order: c.dump.initCount,
		Type:  fmt.Sprintf("%s", rt),
		ID:    string(comp.id),
	}

	var file string
	if comp.caller != nil {
		file = strings.TrimPrefix(comp.caller.file, c.dump.overrideRoot.root)
		file = fmt.Sprintf("%s%s", c.dump.overrideRoot.override, file)

		dumpFile(c, file)

		cJSON.File = file
		cJSON.LineStart, cJSON.LineEnd = findStartEnd(c, file, comp.caller.line)
	}

	for _, d := range c.dump.got {
		cJSON.DirectDeps = append(
			cJSON.DirectDeps,
			&compCoordsJSON{Type: fmt.Sprintf("%s", reflect.TypeOf(d.c)), ID: string(d.id)},
		)
	}

	c.dump.data.Initialized = append(c.dump.data.Initialized, cJSON)
	c.dump.initCount += 1
}

func dumpStageImpl(c *container, stg *stage, comp *component) {
	if c.dump == nil {
		return
	}

	c.dump.mu.Lock()
	defer c.dump.mu.Unlock()

	c.dump.data.StageImpls[string(stg.id)] = append(c.dump.data.StageImpls[string(stg.id)], &compCoordsJSON{
		Type: fmt.Sprintf("%s", reflect.TypeOf(comp.c)),
		ID:   string(comp.id),
	})
}

func dumpFile(c *container, path string) {
	if c.dump == nil {
		return
	}

	if _, ok := c.dump.data.Files[path]; ok {
		return
	}

	bb, err := os.ReadFile(path)
	if err != nil {
		c.LogWarnf("dumping file: %s", err)
		return
	}

	c.dump.data.Files[path] = strings.Split(string(bb), "\n")
}

func writeDump(c *container) {
	if c.dump == nil {
		return
	}

	c.dump.mu.Lock()
	defer c.dump.mu.Unlock()

	for _, s := range c.stages {
		order := ""
		switch s.order {
		case NoOrder:
			order = ""
		case InitOrder:
			order = "InitOrder"
		case ReverseInitOrder:
			order = "ReverseInitOrder"
		}

		c.dump.data.Stages = append(c.dump.data.Stages, &stageJSON{
			ID:       string(s.id),
			Order:    order,
			Parallel: !s.disableParallel,
		})
	}

	if err := writeResources(c); err != nil {
		c.LogWarnf("%s", err)
		return
	}

	d, err := json.Marshal(c.dump.data)
	if err != nil {
		c.LogWarnf("error marshaling dump: %s", err)
		return
	}

	jsData := "var data = " + string(d)
	if err = os.WriteFile(filepath.Join(c.dump.dir, "data.js"), []byte(jsData), os.ModePerm); err != nil {
		c.LogWarnf("%s", err)
		return
	}
}

func writeResources(c *container) error {
	for _, resource := range dumpResources {
		resource.file = filepath.Join(c.dump.dir, resource.file)
		if err := os.MkdirAll(filepath.Dir(resource.file), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(resource.file, []byte(resource.content), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func findStartEnd(c *container, file string, line int) (int, int) {
	firstLine := line - 1
	fileLines, ok := c.dump.data.Files[file]
	if !ok {
		return firstLine, firstLine + 1
	}

	if len(fileLines) < line-1 {
		return firstLine, firstLine + 1
	}

	braces := 0
	for i := line - 1; i < len(fileLines); i++ {
		for _, r := range fileLines[i] {
			if r == '(' {
				braces += 1
			}
			if r == ')' {
				braces -= 1
			}
		}
		if braces == 0 {
			return firstLine, i + 1
		}
	}

	return firstLine, firstLine + 1
}
