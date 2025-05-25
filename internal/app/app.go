package app

import (
	"github.com/oyaro-tech/dispatchgo/internal/flags"
	"github.com/oyaro-tech/dispatchgo/internal/playbook"
	"github.com/oyaro-tech/dispatchgo/internal/ssh"
	"log"
	"slices"
	"sync"
)

type app struct {
	flags    *flags.Flags
	playbook *playbook.Playbook
}

func New() *app {
	return &app{
		flags:    flags.New(),
		playbook: playbook.New(),
	}
}

func (a *app) Run() error {
	a.flags.Parse()

	if err := a.playbook.Parse(a.flags.Playbook()); err != nil {
		return err
	}

	var clients []*ssh.SSH

	log.Print("[@] Connecting to hosts")

	for _, host := range a.playbook.Hosts {
		s := ssh.New()
		if err := s.StartClient(
			host.Name,
			host.Config.Host,
			host.Config.Port,
			host.Config.User,
			host.Config.Password,
			host.Config.PrivateKey,
			host.Config.Passphrase,
		); err != nil {
			return err
		}
		clients = append(clients, s)
	}

	semaphore := make(chan struct{}, a.flags.Threads())

	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(c *ssh.SSH) {
			defer wg.Done()
			defer func() { <-semaphore }()

			clientName := c.Name()

		out:
			for name, details := range a.playbook.Jobs {
				log.Printf("[%s] %s job in progress", clientName, name)

				var jobMatch bool

				if slices.Contains(details.Hosts, c.Name()) {
					jobMatch = true
				}

				if jobMatch {
					for _, task := range details.Tasks {
						log.Printf("[%s] %s", clientName, task.Name)

						output, err := c.RunCommand(task.Script)
						if err != nil {
							log.Printf("[%s] Error running command: %s", clientName, err.Error())
							break out
						}

						if a.flags.Debug() && len(output) > 0 {
							log.Printf("[%s] %s", clientName, output)
						}
					}
				}
			}
		}(client)
	}

	wg.Wait()

	var closeWg sync.WaitGroup

	for _, client := range clients {
		closeWg.Add(1)
		semaphore <- struct{}{}

		go func(c *ssh.SSH) {
			defer closeWg.Done()
			defer func() { <-semaphore }()

			_ = c.CloseClient()
		}(client)
	}

	closeWg.Wait()

	return nil
}

func (a *app) IsDebug() bool {
	return a.flags.Debug()
}
