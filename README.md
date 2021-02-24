cd into $GOPATH/src/desktopClient/cmd/desktopClient and build like this with injected variables:
go build -o SingleProxyDesktopClient_1_1_2 -ldflags "-X desktopClient/config.Key=ac62247a-f1af-4b9f-87c3-8153fbe182c9 -X desktopClient/config.InjectedRemoteSshPort=10060"

to cross compile for other OSes:
env GOOS=windows GOARCH=amd64 go build -o SingleProxyDesktopClient_1_1_2 -ldflags "-X desktopClient/internal.Key=ac62247a-f1af-4b9f-87c3-8153fbe182c9 -X desktopClient/internal.InjectedRemoteSshPort=10060"

logging is configured in config.go to be written to /var/log/desktopClient.log

to check service status:
systemctl status SingleProxyDesktopClient_1_1_2

to remove service:
sudo systemctl stop SingleProxyDesktopClient_1_1_2
sudo systemctl disable SingleProxyDesktopClient_1_1_2
sudo rm /etc/systemd/system/SingleProxyDesktopClient_1_1_2.service
sudo systemctl daemon-reload
sudo systemctl reset-failed


sudo systemctl stop SingleProxyDesktopClient_1_1_2 && sudo systemctl disable SingleProxyDesktopClient_1_1_2 && sudo rm /etc/systemd/system/SingleProxyDesktopClient_1_1_2.service && sudo systemctl daemon-reload && sudo systemctl reset-failed