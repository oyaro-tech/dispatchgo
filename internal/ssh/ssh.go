package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/oyaro-tech/dispatchgo/internal/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SSH struct {
	name   string
	host   string
	config *ssh.ClientConfig
	client *ssh.Client
}

func New() *SSH {
	return &SSH{
		name:   "",
		host:   "",
		config: &ssh.ClientConfig{},
		client: nil,
	}
}

func (s *SSH) StartClient(name, host string, port int, user, password, private_key, passphrase string) error {
	err := s.authorize(user, password, private_key, passphrase)
	if err != nil {
		return fmt.Errorf("Failed to authorize SSH client: %w", err)
	}

	s.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), s.config)
	if err != nil {
		return fmt.Errorf("Failed to dial host %s:%d: %w", host, port, err)
	}

	if host != "" {
		s.host = host
	}

	if name != "" {
		s.name = name
	}

	return nil
}

func (s *SSH) CloseClient() error {
	if s.client != nil {
		return s.client.Close()
	}

	return nil
}

func (s *SSH) RunCommand(cmd string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("SSH client not started. Call StartClient first")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create SSH session: %w", err)
	}

	defer session.Close()

	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf

	if err := session.Run(cmd); err != nil {
		if errors.Is(err, io.EOF) {
			if outBuf.Len() > 0 || errBuf.Len() > 0 {
				return outBuf.String(), fmt.Errorf(
					"Command \"%s\" completed with EOF, but output was received:\nStderr: %s\nError:%w",
					cmd,
					errBuf.String(),
					err,
				)
			}

			return "", fmt.Errorf(
				"Command \"%s\" completed with EOF. Connection likely closed unexpectedly.:\nStderr: %s\nError:%w",
				cmd,
				errBuf.String(),
				err,
			)
		}

		return "", fmt.Errorf(
			"Failed to run command: \"%s\": %w\nStderr: %s",
			cmd,
			err,
			errBuf.String(),
		)
	}

	return outBuf.String(), nil
}

func (s *SSH) authorize(user, password, private_key, passphraze string) error {
	file, err := utils.ExpandTilde(private_key)
	if err != nil {
		return fmt.Errorf(
			"Failed to expand private key path \"%s\": %w",
			private_key,
			err,
		)
	}

	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf(
				"Failed to open private key file \"%s\": %w",
				private_key,
				err,
			)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf(
				"Failed to read private key file \"%s\": %w",
				private_key,
				err,
			)
		}

		if private_key != "" {
			var signer ssh.Signer

			if passphraze != "" {
				signer, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(passphraze))
				if err != nil {
					return fmt.Errorf(
						"Failed to parse private key with passphrase \"%s\": %w",
						private_key,
						err,
					)
				}
			} else {
				signer, err = ssh.ParsePrivateKey(data)
				if err != nil {
					return fmt.Errorf(
						"Failed to parse private key \"%s\": %w",
						private_key,
						err,
					)
				}
			}

			s.config.Auth = append(s.config.Auth, ssh.PublicKeys(signer))
		}
	}

	knownHostsfile, err := utils.ExpandTilde("~/.ssh/known_hosts")
	if err != nil {
		return fmt.Errorf("Failed to expand known_hosts path: %w", err)
	}

	s.config.HostKeyCallback, err = knownhosts.New(knownHostsfile)
	if err != nil {
		return fmt.Errorf("Failed to create known_hosts callback from \"%s\": %w", knownHostsfile, err)
	}

	s.config.User = user

	if password != "" {
		s.config.Auth = append(s.config.Auth, ssh.Password(password))
	}

	return nil
}

func (ssh *SSH) Host() string {
	return ssh.host
}

func (ssh *SSH) Name() string {
	return ssh.name
}
