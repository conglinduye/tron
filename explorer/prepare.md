# TRON golang grpc stub

[toc]

## 1. tron golang stub 生成
1. 准备golang proto环境
    1. install proto buffer c++
        ```shell
        mkdir -p $GOPATH/src/github.com/protocolbuffers
        cd $GOPATH/src/github.com/protocolbuffers
        git clone https://github.com/protocolbuffers/protobuf.git
        cd protobuf
        git submodule update --init --recursive
        ./autogen.sh
        ./configure
        make -j8
        make -j8 check
        sudo make install
        sudo ldconfig
        ```
    2. install golang proto buffer
        ```shell
        go get -u github.com/golang/protobuf/proto
        go get -u github.com/golang/protobuf/protoc-gen-go
        echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bashrc
        source ~/.bashrc
        ```
2. 准备golang grpc环境
    ```
    go get -u google.golang.org/grpc
    ```

3. 准备第三方包
    1. mysql:::`go get -u github.com/go-sql-driver/mysql`
    2. uuid:::`go get github.com/satori/go.uuid`
    3. log:::
        ```shell
        mkdir -p $GOPATH/src/github.com/golang
        git clone https://github.com/golang/glog.git
        ```
    4. kafka:::
        ```shell
        mkdir -p $GOPATH/src/Shopify
        cd $GOPATH/src/Shopify
        git clone https://github.com/Shopify/sarama.git
        ```
    5. redis:::`go get github.com/go-redis/redis`
        ```shell
        
        ```
    6. base58:::`go get -u github.com/btcsuite/btcutil`
    7. go-ethereum:::
        ```shell
        mkdir -p $GOPATH/src/github.com/ethereum && cd $GOPATH/src/github.com/ethereum && git clone https://github.com/ethereum/go-ethereum.git
        // build libsecp256k1
        cd $GOPATH/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1
        ./autogen.sh
        ./configure
        make
        sudu make install
        ```
    8. mongo:::`go get -u gopkg.in/mgo.v2`

4. 获取tron接口协议
    1. tron接口定义 (grpc定义: api/api.proto, 数据类型定义: core/*.proto)
        ```shell
        mkdir $GOPATH/src/github.com/tronprotocol
        cd $GOPATH/src/github.com/tronprotocol
        git clone https://github.com/tronprotocol/protocol.git
        ```
    2. googleapis（作用: 使用了google grpc-gateway将grpc转化为http接口）
        ```shell
        mkdir $GOPATH/src/github.com/googleapis
        cd $GOPATH/src/github.com/googleapis
        git clone https://github.com/googleapis/googleapis.git
        ```
5. 编译
    ```shell
    export DST_DIR=~/go/src
    cd $DST_DIR
    protoc -Igithub.com/tronprotocol/protocol -Igithub.com/googleapis/googleapis github.com/tronprotocol/protocol/api/api.proto --go_out=plugins=grpc:$DST_DIR
    protoc -Igithub.com/tronprotocol/protocol -Igithub.com/googleapis/googleapis github.com/tronprotocol/protocol/core/*.proto --go_out=plugins=grpc:$DST_DIR
    cd -
    ```
    
    执行后会生成生成golang grpc stub 文件, 位置为: 
    + api: `github.com/tronprotocol/grpc-gateway/api`
    + core: `github.com/tronprotocol/grpc-gateway/api`
    
    验证是否可用(go build 通过)
    ```shell
    cd $GOPATH/src/github.com/tronprotocol/grpc-gateway/api
    go build
    cd $GOPATH/src/github.com/tronprotocol/grpc-gateway/core
    go build
    ```
6. DB init
    tron  主网数据库       
    tron_test_net 测试网数据库      
    二者表结构大体相同，表结构文件：main/schema/tron.sql          
    local DB data init:
    ```shell
    cd main/fullnode
    go build
    nohup ./fullnode -net main -start_block 2500000 -dsn "budev:tron**1@tcp(127.0.0.1:3306)/tron" > sync-main.log 2>&1 &
    nohup ./fullnode -net test -dsn "budev:tron**1@tcp(127.0.0.1:3306)/tron_test_net" > sync.log 2>&1 &
    ```
    + -dsn 本地数据库连接信息
    + make sure redis-server started at 6379


## 2. explorerService
1.  安装gin环境
```
go get -u github.com/gin-gonic/gin
```

2. import it in your code:
```
import "github.com/gin-gonic/gin"
```

3. API DOC
http://test.tronapp.co:8000/blockchain/