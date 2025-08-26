package vmware_esxi

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/crypto/ssh"
)

func (v *VMwareESXi) Consume(c *certificate.Resource) error {
	// Execute certificate installation
	_, err := methods.InstallServerCertificate(context.Background(), v.Client, &types.InstallServerCertificate{
		This: types.ManagedObjectReference{
			Type:  "HostCertificateManager",
			Value: "ha-certificate-manager",
		},
		Cert: string(c.Certificate),
	})
	if err != nil {
		return fmt.Errorf("certificate installation failed: %v", err)
	}

	v.Logger.Info("Certificate installed successfully via vim25 API")
	return nil
}

func (v *VMwareESXi) startSSH() error {
	serviceName := "TSM-SSH"
	v.Logger.Infof("Starting service: %s\n", serviceName)

	// Execute service restart
	_, err := methods.StartService(context.Background(), v.Client, &types.StartService{
		This: types.ManagedObjectReference{
			Type:  "HostServiceSystem",
			Value: "serviceSystem",
		},
		Id: serviceName,
	})
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}

	v.Logger.Infof("Service %s started successfully\n", serviceName)
	return nil
}

func (v *VMwareESXi) RestartManagementServices() error {
	err := v.startSSH()
	if err != nil {
		return err
	}
	// Services to restart for certificate changes
	servicesToRestart := []string{"hostd", "vpxa", "rhttpproxy"}

	v.Logger.Debug("Restarting management services...")

	addr, err := url.Parse(v.URL)
	if err != nil {
		return err
	}

	for _, serviceName := range servicesToRestart {
		_, err := v.remoteRun(addr.Host, v.Username, v.Password, fmt.Sprintf("/etc/init.d/%s restart", serviceName))
		if err != nil {
			v.Logger.WithError(err).Warningf("Failed to restart %s", serviceName)
			continue
		}

		// Wait between service restarts
		v.Logger.Info("Waiting 3 seconds before next service...\n")
		time.Sleep(3 * time.Second)
	}

	v.Logger.Debug("Management services restart completed")
	return nil
}

func (v *VMwareESXi) remoteRun(addr string, user string, password string, cmd string) (string, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) > 0 {
					return []string{password}, nil
				}
				return []string{}, nil
			}),
		},
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "22"), config)
	if err != nil {
		return "", err
	}
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			v.Logger.WithError(err).Warning("failed to close SSH session")
		}
	}()
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(cmd)
	return b.String(), err
}
