# Server config / build

This will/should be scripted at some point, but for now, I am just going to take notes.

## Install required packages

`sudo apt-get install git`

## Part 1: gitreceive

### Step 1: install gitreceive
```bash
(
  cd /usr/local/bin \
    && sudo wget https://raw.githubusercontent.com/progrium/gitreceive/v1.0.0/gitreceive \
    && sudo chmod +x gitreceive
)

`sudo gitreceive init`
```

### Step 2: setup receiver script

Modify `/home/git/receiver` to fit our needs. See the `source/receiver` file in this repo as use it's content in place of the original content.

### Step 2.5: if you are using Vagrant or a non-standard SSH port...

Modify `~/.ssh/config` on your local machine (if you are using a differnet port, with Vagrant, as an example):
```bash
Host localhost
    Port 2222
```

### Step 3: add a git remote and push

```bash
# add ssh key to user (for Vagrant)
cat ~/.ssh/id_rsa.pub | ssh vagrant@localhost -p 2222 -i /Users/jipiboily/.vagrant.d/insecure_private_key "sudo gitreceive upload-key jp"

# add remote
git remote add paas git@localhost:example

# push!
git push paas
```

That's about it for part 1, we are now able to git push our app.
