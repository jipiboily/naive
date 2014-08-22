# Server config / build

This will/should be scripted at some point, but for now, I am just going to take notes.


## gitreceive

```bash
(
  cd /usr/local/bin \
    && sudo wget https://raw.githubusercontent.com/progrium/gitreceive/master/gitreceive \
    && sudo chmod +x gitreceive
)

# git + dispatcher = gitpatcher
sudo GITUSER=gitpatcher gitreceive init
```

Modify `/home/gitpatcher/receiver` to fit our needs...(to be continued)
