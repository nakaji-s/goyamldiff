# Go Yaml Diff
## Install
$ go get github.com/nakaji-s/goyamldiff 

## Usage
$ goyamldiff file1.yaml file2.yaml
```diff
nuooo:
  - 5
  + a:
      + 2
nuwaa:
  - b:
      - 3
  + 6
mail_from:
  - "hoge@fuga.com"
  + "hoge@fuga.com2"
var:
    var1:
      - 2
      + 3
    static_ips:
      - ["10.244.0.26"]
      + ["10.244.0.26", "10.244.0.50", "10.244.0.54"]
    var3:
      + 5
nulllll:
  - null
guooo:
  + 1234
```

## License
MIT
