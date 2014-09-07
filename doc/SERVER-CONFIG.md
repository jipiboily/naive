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


## Part 2: Docker

To install on a recent version on Ubuntu 14.04:

```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
sudo sh -c "echo deb https://get.docker.io/ubuntu docker main\
> /etc/apt/sources.list.d/docker.list"
sudo apt-get update
sudo apt-get install lxc-docker
```

Let the `git` user deal with Docker.
Source: http://askubuntu.com/questions/477551/how-can-i-use-docker-without-sudo

```
sudo groupadd docker # probably not required
sudo gpasswd -a git docker
sudo service docker restart
```

This can also be useful if your are working on a Vagrant box: `sudo gpasswd -a vagrant docker`

## Part 3: nginx

`docker pull dockerfile/nginx`
basic, temporary version:
`docker run -d -p 80:80 dockerfile/nginx`

mkdir -p /vagrant/.config/nginx/sites-enabled
mkdir -p /vagrant/.logs/nginx

Sample nginx config:
```
server {
      listen 80;

      location / {
          proxy_pass http://httpstat.us/;
          proxy_set_header  X-Real-IP  $remote_addr;
      }
  }
```
docker run -d --name nginx -p 80:80 -v /vagrant/.config/nginx/sites-enabled:/etc/nginx/sites-enabled -v /vagrant/.logs/nginx:/var/log/nginx dockerfile/nginx
