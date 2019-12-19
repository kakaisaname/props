# props

[![Build Status](https://travis-ci.org/kakaisaname/props.svg?branch=master)](<https://travis-ci.org/kakaisaname/props>)
[![GoDoc Documentation](http://godoc.org/github.com/kakaisaname/props?status.png)](<https://godoc.org/github.com/kakaisaname/props>)
[![Sourcegraph](https://sourcegraph.com/github.com/kakaisaname/props/-/badge.svg)](https://sourcegraph.com/github.com/kakaisaname/props?badge)
[![Coverage Status](https://coveralls.io/repos/github/kakaisaname/props/badge.svg?branch=master)](https://coveralls.io/github/kakaisaname/props?branch=master)
[![GitHub release](https://img.shields.io/github/release/kakaisaname/props.svg)](https://github.com/kakaisaname/props/releases)

 

 统一的配置工具库，将各种配置源抽象或转换为类似properties格式的key/value，并提供统一的API来访问这些key/value。支持 properties 文件、ini 文件、zookeeper k/v、zookeeper k/props、consul k/v、consul k/props等配置源，并且支持通过 Unmarshal从配置中抽出struct；支持上下文环境变量的eval，${}形式；支持多种配置源组合使用。


## 特性
### 支持的配置源：

- properties格式文件
- ini格式文件
- yaml格式文件
- [Apollo](<https://github.com/ctripcorp/apollo>) k/v,k/props,k/ini,k/yaml
- [Nacos](<http://nacos.io>) k/props[properties],k/yaml,k/ini,k/ini_props
- zookeeper k/v
- zookeeper k/props[properties],k/yaml,k/ini,k/ini_props
- consul k/v
- consul k/props[properties],k/yaml,k/ini,k/ini_props
- etcd API V2 k/v
- etcd API V2 k/props
- etcd API V3 k/v
- etcd API V3 k/props

### key/value支持的数据类型：

- key只支持string
- value 5种数据类型的支持：
	- string
	- int
	- float64
	- bool
	- time.Duration：
	    - 比如 "300ms", "-1.5h" or "2h45m". 
	    - 合法的时间单位： "ns", "us" (or "µs"), "ms", "s", "m", "h".

### 其他特性

- Unmarshal支持
- 上下文变量eval支持，`${}`形式
- 支持多配置源组合
- 默认添加了系统环境变量，优先级最低

## Install

> go get -u github.com/kakaisaname/props

**或者通过go mod：**


> go mod tidy 
> go mod vendor





## 配置源和配置形式使用方法：

### properties格式文件

格式：`[key][=|:][value] \n`
每行为key/value键值对 ,用`=`或`：`分割，key可以是除了`=`和`:`、以及空白字符的任何字符

例子：

`server.port=8080`

或者

`server.port: 8080`


### 通过kvs.ReadPropertyFile读取文件

```golang

p, err := kvs.ReadPropertyFile("config.properties")
if err != nil {
	panic(err)
}
stringValue, err := p.Get("prefix.key1")
fmt.Println(stringValue, err)
//如果不存在，则返回默认值
stringDefaultValue := p.GetDefault("prefix.key1", "default value")
fmt.Println(stringDefaultValue)
boolValue, err := p.GetBool("prefix.key2")
fmt.Println(boolValue)
boolDefaultValue := p.GetBoolDefault("prefix.key2", false)
fmt.Println(boolDefaultValue)
intValue, err := p.GetInt("prefix.key3")
fmt.Println(intValue)
intDefaultValue := p.GetIntDefault("prefix.key3", 1)
fmt.Println(intDefaultValue)
floatValue, err := p.GetFloat64("prefix.key4")
fmt.Println(floatValue)
floatDefaultValue := p.GetFloat64Default("prefix.key4", 1.2)
fmt.Println(floatDefaultValue)

```

#### 通过kvs.NewProperties()从io.Reader中读取

```
 p := kvs.NewProperties()
 p.Load(strings.NewReader("some data"))
 p.Load(bytes.NewReader([]byte("some data")))
```

#### 通过kvs.NewPropertiesConfigSource()

```
file := "/path/to/config.properties"
p := kvs.NewPropertiesConfigSource(file)
p = kvs.NewPropertiesConfigSourceByFile("name", file)
//通过map构造内存型
m := make(map[string]string)
m["key"]="value"
p = kvs.NewPropertiesConfigSourceByMap("name", m)

```
#### Properties ConfigSource

```golang

var cs kvs.ConfigSource
//
cs = kvs.NewPropertiesConfigSource("config.properties")
cs = kvs.NewPropertiesConfigSourceByFile("config", "config.properties")


stringValue, err := cs.Get("prefix.key1")
//如果不存在，则返回默认值
stringDefaultValue := cs.GetDefault("prefix.key1", "default value")
boolValue, err := cs.GetBool("prefix.key2")
boolDefaultValue := cs.GetBoolDefault("prefix.key2", false)
intValue, err := cs.GetInt("prefix.key3")
intDefaultValue := cs.GetIntDefault("prefix.key3", 1)
floatValue, err := cs.GetFloat64("prefix.key4")
floatDefaultValue := cs.GetFloat64Default("prefix.key4", 1.2)

```




### ini格式文件。

格式：参考 [wiki百科：INI_file](<https://en.wikipedia.org/wiki/INI_file>)
​	
```ini
[section]
[key1][=|:][value1] 
[key1][=|:][value1]
...
```

不支持sub section

例子：

```ini
[server]
port: 8080
read.timeout=6000ms

[client]
connection.timeout=6s
query.timeout=6s
```

#### 使用方法：

```golang
file := "/path/to/config.ini"
p := ini.NewIniFileConfigSource(file)
p = ini.NewIniFileConfigSourceByFile("name", file)
```

### Nacos

只支持key/properties配置形式。

例如有如下配置：

http://127.0.0.1:8848/nacos/v1/cs/configs?dataId=test.id&group=testGroup&tenant=testTenant

```properties
key-0.x0=value-00
key-0.x1=value-01
key-0.x2=value-02
```

 

基本用法：

```go
 address := "127.0.0.1:8848"
 c := NewNacosPropsConfigSource(address)
 c.DataId = "test.id"
 c.Tenant = "testTenant"
 c.Group = "testGroup"
 v:=c.GetDefault("key-0.x0","defaultval") //value-00

```



### zookeeper 

支持key/value和key/properties配置形式，key/properties配置和ini类似，将key作为section name。

key/value形式，将path去除root path部分并替换`/`为`.`作为key。

key/properties形式，在root path下读取所有子节点，将子节点名称作为section name，value为子properties格式内存，通过子节点名称和子properties中的key组合成新的key作为key。

#### by zookeeper key/value

##### 基本例子

```golang
root := "/config/kv/app1/dev"
var conn *zk.Conn
p := zk.NewZookeeperConfigSource("zookeeper-kv", root, conn)
```

##### CompositeConfigSource多context例子

```golang
var cs kvs.ConfigSource
urls := []string{"172.16.1.248:2181"}
contexts := []string{"/configs/apps","/configs/users"}
cs = zk.NewZookeeperCompositeConfigSource(contexts, urls, time.Second*3)

```

#### 用properties来配置： key/properties

value值为properties格式内容, 整体设计类似ini格式,例如：

##### key: 
/config/kv/app1/dev/datasource
##### value:

```properties
url=tcp(127.0.0.1:3306)/Test?charset=utf8
username=root
password=root

```



```golang
root := "/config/kv/app1/dev"
var conn *zk.Conn
p := zk.NewZookeeperIniConfigSource("zookeeper-props", root, conn)

```


### consul 多层key/value形式


#### by consul key/value

```golang
例如：

config101/test/demo1/server/port=8080

获取的属性和值是：

server.port=8080

address := "127.0.0.1:8500"
root := "config101/test/demo1"
c := consul.NewConsulKeyValueConfigSource("consul", address, root)
stringValue, err := cs.Get("prefix.key1")
stringDefaultValue := cs.GetDefault("prefix.key1", "default value")

```
#### 用properties来配置： key/properties

value值为properties格式内容, 整体设计类似ini格式,配置样式如下图：

![](<docs/consul_key_kvs.png>)

```golang
root := "config/app1/dev"
address := "127.0.0.1:8500"
p := consul.NewConsulIniConfigSourceByName("consul-props", address, root)
```

### 支持Unmarshal

支持的数据类型：
- int,int8,int16,int32,int64
- uint,uint8,uint16,uint32,uint64
- string
- bool
- float32,float64
- map
- time.Duration
- 嵌套struct
- map：key只支持string，value支持以上除struct的基本类型


在struct中规定命名为`_prefix `、类型为`string `、并且指定了`prefix`tag, 使用feild `_prefix `的`prefix`tag作为前缀，将struct feild名称转换后组合成完整的key，并从ConfigSource中获取数据并注入struct实例，feild类型只支持ConfigSource所支持的数据类型（string、int、float、bool、time.Duration）。

```golang


type Port struct {
    Port    int  `val:"8080"`
    Enabled bool `val:"true"`
}
type ServerProperties struct {
    _prefix string        `prefix:"http.server"`
    Port    Port
    Timeout int           `val:"1"`
    Enabled bool
    Foo     int           `val:"1"`
    Time    time.Duration `val:"1s"`
    Float   float32       `val:"0.000001"`
    Params  map[string]string
    Times      map[string]time.Duration
}

func main() {
   
    p := kvs.NewMapProperties()
    p.Set("http.server.port.port", "8080")
    p.Set("http.server.params.k1", "v1")
    p.Set("http.server.params.k2", "v2")
    p.Set("http.server.Times.m1", "1s")
    p.Set("http.server.Times.m2", "1h")
    p.Set("http.server.Times.m3", "1us")
    p.Set("http.server.port.enabled", "false")
    p.Set("http.server.timeout", "1234")
    p.Set("http.server.enabled", "true")
    p.Set("http.server.time", "10s")
    p.Set("http.server.float", "23.45")
    p.Set("http.server.foo", "23")
    s := &ServerProperties{
        Foo:   1234,
        Float: 1234.5,
    }
    p.Unmarshal(s)
    fmt.Println(s)

}


```


### 上下文变量表达式（或者占位符）的支持

支持在props上下文中替换占位符：`${}` 

```
p := kvs.NewEmptyMapConfigSource("map2")
p.Set("orign.key1", "v1")
p.Set("orign.key2", "v2")
p.Set("orign.key3", "2")
p.Set("ph.key1", "${orign.key1}")
p.Set("ph.key2", "${orign.key1}:${orign.key2}")
p.Set("ph.key3", "${orign.key3}")
conf := kvs.NewDefaultCompositeConfigSource(p)
phv1, err := conf.GetInt("ph.key1")//v1
phv2, err := conf.Get("ph.key2")//v1:v1
phv3, err := conf.GetInt("ph.key3")//2

```


### 多种配置源组合使用

优先级以追加相反的顺序,最后添加优先级最高。

```golang

kv1 := []string{"go.app.key1", "value1", "value1-2"}
kv2 := []string{"go.app.key2", "value2", "value2-2"}

p1 := kvs.NewEmptyMapConfigSource("map1")
p1.Set(kv1[0], kv1[1])
p1.Set(kv2[0], kv2[1])
p2 := kvs.NewEmptyMapConfigSource("map2")
p2.Set(kv1[0], kv1[2])
p2.Set(kv2[0], kv2[2])
conf.Add(p1)
conf.Add(p2)

//value1==value1-2
value1, err := conf.Get(kv1[0])
fmt.Println(value1)
//value2=value2-2
value2, err := conf.Get(kv2[0])
fmt.Println(value2)

	
```

