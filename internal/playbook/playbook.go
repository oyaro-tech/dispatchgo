package playbook

import (
	"errors"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type hostConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
}

type hostEntry struct {
	Name   string
	Config hostConfig
}

type job struct {
	Hosts []string `yaml:"hosts"`
	Tasks []task   `yaml:"tasks"`
}

type task struct {
	Name   string `yaml:"name"`
	Script string `yaml:"script"`
}

type Playbook struct {
	Name  string         `yaml:"playbook_name"`
	Hosts []hostEntry    `yaml:"hosts"`
	Jobs  map[string]job `yaml:"jobs"`
}

func New() *Playbook {
	return &Playbook{
		Name:  "",
		Hosts: make([]hostEntry, 0),
		Jobs:  make(map[string]job),
	}
}

func (p *Playbook) Parse(file string) error {
	if err := p.loadFile(file); err != nil {
		return err
	}

	return nil
}

func (p *Playbook) loadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var row struct {
		Name  string                  `yaml:"playbook_name"`
		Hosts []map[string]hostConfig `yaml:"hosts"`
		Jobs  map[string]job          `yaml:"jobs"`
	}

	if err = yaml.Unmarshal(data, &row); err != nil {
		return errors.New(yaml.FormatError(err, true, true))
	}

	p.Name = row.Name
	p.Jobs = row.Jobs

	for _, hostMap := range row.Hosts {
		for name, config := range hostMap {
			p.Hosts = append(p.Hosts, hostEntry{
				Name:   name,
				Config: config,
			})
		}
	}

	return nil
}
