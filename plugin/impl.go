package plugin

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

// Settings for the plugin.
type Settings struct {
	Address string

	Username string
	Password string
	Key      string

	Source string
	Target string

	Templating bool

	auth []ssh.AuthMethod
}

const defaultScpPort = 22

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	// Validation of the settings.
	if len(p.settings.Address) == 0 {
		return fmt.Errorf("no scp address provided")
	}

	if len(p.settings.Username) == 0 {
		return fmt.Errorf("no scp username provided")
	}

	if len(p.settings.Password) == 0 && len(p.settings.Key) == 0 {
		return fmt.Errorf("no scp password or key provided")
	}

	if len(p.settings.Source) == 0 {
		return fmt.Errorf("no scp source file provided")
	}

	if strings.Index(p.settings.Address, ":") == -1 {
		p.settings.Address += ":" + strconv.Itoa(defaultScpPort)
	}

	if len(p.settings.Password) != 0 {
		p.settings.auth = append(p.settings.auth, ssh.Password(p.settings.Password))
	}

	if len(p.settings.Key) != 0 {
		signer, err := ssh.ParsePrivateKey([]byte(p.settings.Key))
		if err != nil {
			return fmt.Errorf("error while parsing private key: %w", err)
		}

		p.settings.auth = append(p.settings.auth, ssh.PublicKeys(signer))
	}

	if len(p.settings.Target) == 0 {
		p.settings.Target = path.Base(p.settings.Source)
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	fmt.Printf("Connecting to %s\n", p.settings.Address)

	client, err := ssh.Dial("tcp", p.settings.Address,
		&ssh.ClientConfig{
			User:            p.settings.Username,
			Auth:            p.settings.auth,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
	if err != nil {
		return fmt.Errorf("error while connecting to %s: %w", p.settings.Address, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error while opening new session: %w", err)
	}
	defer session.Close()

	fmt.Printf("Copying file: %s\n", p.settings.Source)

	f, err := os.Open(p.settings.Source)
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error while stating file: %w", err)
	}

	var contents io.Reader
	var size int64

	if p.settings.Templating {
		tmpl, err := template.ParseFiles(p.settings.Source)
		if err != nil {
			return fmt.Errorf("error while parsing template file: %w", err)
		}

		var tpl bytes.Buffer

		m := make(map[string]interface{})
		for _, e := range os.Environ() {
			if p := strings.SplitN(e, "=", 2); len(p) == 2 {
				m[p[0]] = p[1]
			}
		}

		err = tmpl.Execute(&tpl, m)
		if err != nil {
			return fmt.Errorf("error while executing template: %w", err)
		}

		var v = tpl.String()
		contents = strings.NewReader(v)
		size = int64(len(v))
	} else {
		contents = f
		size = s.Size()
	}

	err = scp.Copy(size, s.Mode().Perm(), path.Base(p.settings.Source), contents, p.settings.Target, session)
	if err != nil {
		return fmt.Errorf("error while copying file: %w", err)
	}

	return nil
}
