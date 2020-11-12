package api

import (
	"fmt"
	"log"
	"strings"
)

// VarSet - set project variable
func (p *Project) VarSet(key string, value string) error {
	log.Printf("Set var '%s.'", key)
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return fmt.Errorf("invalid variable key")
	}
	if p.Variables[keySplit[0]] == nil {
		p.Variables[keySplit[0]] = make(map[string]string)
	}
	p.Variables[keySplit[0]][keySplit[1]] = value
	return nil
}

// VarGet - retrieve project variable
func (p *Project) VarGet(key string) (string, error) {
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return "", fmt.Errorf("invalid variable key")
	}
	if p.Variables[keySplit[0]] == nil {
		return "", nil
	}
	return p.Variables[keySplit[0]][keySplit[1]], nil
}

// VarDelete - delete a project variable
func (p *Project) VarDelete(key string) error {
	log.Printf("Delete var '%s.'", key)
	keySplit := strings.Split(key, ":")
	if len(keySplit) != 2 {
		return fmt.Errorf("invalid variable key")
	}
	if p.Variables[keySplit[0]] == nil {
		return nil
	}
	delete(p.Variables[keySplit[0]], keySplit[1])
	return nil
}
