# deploycloud 

*cf plugin for using a github remote for managing cf app deployments*

[![wercker status](https://app.wercker.com/status/9a6553ba12248db71c8e452c2723e6c3/s/master "wercker status")](https://app.wercker.com/project/bykey/9a6553ba12248db71c8e452c2723e6c3)

## installation from release binaries

**On OSX using release binaries**
```
$ wget -O $GOPATH/bin/deploycloud https://github.com/xchapter7x/deploycloud/releases/download/v1.0.1/deploycloud_darwin_amd64
$ cf install-plugin $GOPATH/bin/deploycloud
```

**On Linux using release binaries**
```
$ wget -O $GOPATH/bin/deploycloud https://github.com/xchapter7x/deploycloud/releases/download/v1.0.1/deploycloud_linux_amd64
$ cf install-plugin $GOPATH/bin/deploycloud
```

**On Windows using release binaries**
```
#in your browser download https://github.com/xchapter7x/deploycloud/releases/download/v1.0.1/deploycloud_windows_amd64.exe and place it in $GOPATH/bin/deploycloud.exe

$ cf install-plugin $env:GOPATH/bin/deploycloud.exe
```

## installation using go get

**On *nix using go get**
```
$ go get github.com/xchapter7x/deploycloud
$ cf install-plugin $GOPATH/bin/deploycloud
```

**On Windows using go get**
```
$ go get github.com/xchapter7x/deploycloud
$ cf install-plugin $env:GOPATH/bin/deploycloud.exe
```

**Usage**
```
#List available apps to deploy
$ cf cloud-deploy --list --org xchapter7x --repo deploycloud --token <ghtoken> --cfuser <myuser> --cfpass <mypass>


#Show details of deployment
$ cf cloud-deploy --show myapp.development --org xchapter7x --repo deploycloud --token <ghtoken> --cfuser <myuser> --cfpass <mypass>

#Run an available application deployment
$ cf cloud-deploy --run myapp.development --org xchapter7x --repo deploycloud --token <ghtoken> --cfuser <myuser> --cfpass <mypass>
```

## Env Var powered
**by setting the following env vars you can avoid the need to pass them into the cli**
```
`GH_TOKEN` - env var name to store your github oauth token
`CF_USER` - env var name to store your cf user
`CF_PASS` - env var name to store your cf user's password
```

## Repo Structure
**By default the plugin will look for a `config.yml` file in the root of the given repo.**
**You can overwrite this default by passing the `--config` flag to the plugin with the relative path of the config in the repo**

## Config file Structure

**Sample:**
```
---
applications:
  #an application record
  myapp1:
    #list of deployments associated with this app
    deployments:
      # one of the deployement details records
      - name: development #name of deployment
        url: api.pivotal.io #url to target during push
        org: myorg # org to target during push
        space: thespace # space to target during push
        path: myapp1/development #path to desired manifest (relative to the root of this repo)
        push_cmd: cf push appname -i 2 #push command to execute (note: the above manifest will be used via a added `-f`, so dont add it here)
      
      # another deployment definition for the myapp1 application
      - name: production
        url: api.pivotal.io
        org: myorg
        space: prodspace
        path: myapp1/production
        push_cmd: cf push-zdd appname_prod -i 8
  myotherapp:
    deployments:
      - name: dev
        url: api.other.pivotal.io
        org: otherorg
        space: otherspace
        path: myotherapp/development
        push_cmd: cf push appname_dev -i 2
```
