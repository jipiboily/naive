# Server config / build

This will/should be scripted at some point, but for now, I am just going to take notes.

## Install required packages

`sudo apt-get install git`

## gitreceive

```bash
(
  cd /usr/local/bin \
    && sudo wget https://raw.githubusercontent.com/progrium/gitreceive/v1.0.0/gitreceive \
    && sudo chmod +x gitreceive
)

`sudo gitreceive init`
```

Modify `/home/git/receiver` to fit our needs...

`sudo vi /home/git/receiver`

```bash
#!/bin/bash
mkdir -p /home/git/tmp && cat | tar -x -C /home/git/tmp
```

NOTE: I think we might need to create `/home/git/tmp` before while it should not...

Modify `~/.ssh/config` on your local machine (if you are using a differnet port, with Vagrant, as an example):
```bash
Host localhost
    Port 2222
```

```bash
# add ssh key to user (for Vagrant)
cat ~/.ssh/id_rsa.pub | ssh vagrant@localhost -p 2222 -i /Users/jipiboily/.vagrant.d/insecure_private_key "sudo gitreceive upload-key jp"

# add remote
git remote add paas git@localhost:example
```