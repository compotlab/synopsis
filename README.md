Synopsis - Composer Package Repository Generator
==================

 To run application need install [golang](https://golang.org/doc/install) :
---------
 1. `> go get github.com/compotlab/synopsis`
 2. `> cd $GOPATH/src/github.com/compotlab/synopsis/`
 3. `> cp sample.app.yaml app.yml && cp sample.config.json config.json`
 4. `> go build`
 5. `> ./synopsis`

To show info about building packages `http://localhost:9091/` (port change in app.yaml config file).
To show admin panel go to `http://localhost:9091/admin`
Build config file look like [satis](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md)

Tips:
---------
- If you have some crushing, change `thread` count, because your server can't run so many concurrent processes.

License
====
Licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/compotlab/synopsis/blob/master/LICENSE) for the full license text.