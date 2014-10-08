# Step 1
First step, create a server. I use [Digital Ocean](https://www.digitalocean.com/?refcode=9b3537dd733f) (this is a referral link) and installed Ubuntu 14.04.

# Step 2: first, we'll setup some variables that we will re-use during the setup:
  On local machine AND server, set this: `export PROJECT_NAME="my-vocal-blorgh"`

  On your local machine: `export SERVER_IP=0.0.0.0`


Log on the server and install:
- update your packages...
  ```bash
  sudo apt-get update
  sudo apt-get upgrade
  ```
- git: `sudo apt-get install git`
- docker
  ```bash
  sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
  sudo sh -c "echo deb https://get.docker.io/ubuntu docker main\
  > /etc/apt/sources.list.d/docker.list"
  sudo apt-get update
  sudo apt-get install lxc-docker
  ```


- gitreceive + make `git` user able to deal with docker
    # Install gitreceive
    ```bash
    (
      cd /usr/local/bin \
        && sudo wget https://raw.githubusercontent.com/progrium/gitreceive/v1.0.0/gitreceive \
        && sudo chmod +x gitreceive
    )
    ```

    sudo gitreceive init

    # The `git`user will deal with docker...so...
    # sudo groupadd docker # probably not required
    sudo gpasswd -a git docker
    sudo service docker restart

- naive
    mkdir /naive/templates -p

    # FIGURE OUT HOW THE DISTRIBUTION WILL BE DONE!!!
    $(cd /naive/templates && wget https://raw.githubusercontent.com/jipiboily/naive/master/templates/nginx-site.conf)

    # wget receiver binary ... TODO !!!
    on server...
      mkdir /naive/src -p
    on local...
      scp -r source/* root@$SERVER_IP:/naive/src
    on server...
      sudo apt-get install golang
      $(cd /naive/src && go build receiver.go && cp receiver /home/git/receiver)

    chown git:git -R /naive

- nginx container
    docker pull dockerfile/nginx
    mkdir -p /naive/nginx/config/sites-enabled
    mkdir -p /naive/nginx/logs/
    chown git:git -R /naive

    # add a test config in sites-enabled
      server {
          listen 80;



          location / {
              proxy_pass http://perdu.com;
              proxy_set_header  X-Real-IP  $remote_addr;
          }
      }


    # then start it
    docker run -d --name nginx -p 80:80 -v /naive/nginx/config/sites-enabled:/etc/nginx/sites-enabled -v /naive/nginx/logs:/var/log/nginx dockerfile/nginx

    # remove the test

- install postgresql (optional)
  # For this one, I don't use a container but the host directly, mostly for lazyness of finding the right way of dealing with a database with state in a container.
  sudo apt-get update
  sudo apt-get install postgresql postgresql-contrib

  # Create a new role and DB with same name
  sudo -i -u postgres
  export PROJECT_NAME='my-vocal-blorgh'
  createdb $PROJECT_NAME
  createuser --login $PROJECT_NAME
  psql
  grant all on DATABASE "MYDB" to "MYUSER";
  ALTER USER "my-vocal-blorgh" WITH PASSWORD 'CHANGETHIS'; // Add your own secure password, more secure than that :)

  add this to `/etc/postgresql/9.3/main/pg_hba.conf`:
    `host all all 0.0.0.0/0 md5`
  and this to `/etc/postgresql/9.3/main/postgresql.conf`:
    `listen_addresses = '*'`


- set environment variables (optional)
  You need to edit the `app.env` file. One environment variable per line, in the format of "KEY=VALUE", no comment, no empty line or anything else. Current way of dealing with those is too basic/stupid to support anything else.
  `mkdir -p /naive/apps/$PROJECT_NAME && vi /naive/apps/$PROJECT_NAME/app.env`

  Note that if you want to get the Docker host IP, you can use: `DOCKER_HOST_IP=$(netstat -nr | grep '^0\.0\.0\.0' | awk '{print $2}')` which seems to work although I am not sure this is the best way to achieve this.

- push!
  # first upload key...
  cat ~/.ssh/id_rsa.pub | ssh root@$SERVER_IP "sudo gitreceive upload-key naive"

  git remote add prd git@$SERVER_IP:my-vocal-blorgh
  git push prd master
