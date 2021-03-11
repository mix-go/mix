# Go-Loading



Go-Loading, or more likely `loading` is your custom loader in your every CLI application.



![CircleCI](https://img.shields.io/circleci/project/github/RedSparr0w/node-csgo-parser.svg)

### Installation:
Get it by using: `go get -u github.com/whaangbuu/loading`

### How does it work?
It adds a loader in your CLI Application.

#### See clip below:
[![asciicast](https://asciinema.org/a/3ft6E5zcUUsiyOQPqkErw9NM7.png)](https://asciinema.org/a/3ft6E5zcUUsiyOQPqkErw9NM7)

##### NOTE:
[`go-imgur-cli`](https://github.com/whaangbuu/go-imgur-cli) is my sample CLI application that uploads images to [imgur](https://imgur.com)

#### How to use?
Simply create a new `loading` instance.

```sh
loading := loading.StartNew("<some_title>")
```


#### Changing colors of the loader:


For now, the allowed colors are: `red, green, blue, and white`.

```sh
loading.SetColor("green")
```

#### Changing the default loaders:
```sh
loading.SetLoaders([]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"})
```

#### Stopping the loader:
```sh
loading.Stop()
```

#### Stopping the loader with a specific time.
```sh
time.Sleep(3 * time.Second)
loading.Stop()
```
This stops the loader after 3 seconds.


#### Stopping the loader while waiting for a task.
```sh
// some productive task here.
loading.Stop()
```
See `/examples` directory for more concise guide.


#### Testing:
See `loading_test.go`


#### TODO:
- [ ] Improve tests
- [ ] Write more test
