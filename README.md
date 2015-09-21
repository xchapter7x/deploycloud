# deploycloud 

*cf plugin for using a github remote for managing cf app deployments*

[![wercker status](https://app.wercker.com/status/9a6553ba12248db71c8e452c2723e6c3/s/master "wercker status")](https://app.wercker.com/project/bykey/9a6553ba12248db71c8e452c2723e6c3)

## installation

**On *nix**
```
$ go get github.com/xchapter7x/deploycloud
$ cf install-plugin $GOPATH/bin/deploycloud
```

**On Windows**
```
$ go get github.com/xchapter7x/deploycloud
$ cf install-plugin $env:GOPATH/bin/deploycloud.exe
```

**Usage**
```
#List available apps to deploy
$ cf cloud-deploy --list --org xchapter7x --repo deploycloud --token <ghtoken> --cfuser <myuser> --cfpass <mypass>
```
