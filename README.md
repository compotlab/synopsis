Synopsis - Composer Package Repository Generator
===

 To run application need install [golang](https://golang.org/doc/install) :
---
    `go get github.com/compotlab/synopsis`
    `cd $GOPATH/src/github.com/compotlab/synopsis/`
    `go build`
    `./synopsis`

To show info about building packages `http://localhost:9091/` (port change in app.yaml config file).
To show admin panel go to `http://localhost:9091/admin`
Build `config.json` file look like [satis](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md)

Exampale `config.json` file:
---
```
{
  "name": "Private Package Repository",
  "homepage": "http://localhost:9091",
  "archive": {
    "directory": "dist",
    "skip-dev": false
  },
  "repositories": [
    {
      "type": "vcs",
      "url": "git@github.com:compotlab/synopsis.git"
    }
  ]
}
```

Tips:
---
- If you have some crushing, change `thread` count, because your server can't run a lot of sys call processes.

License
===
Licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/compotlab/synopsis/blob/master/LICENSE) for the full license text.