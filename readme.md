
## 运行环境
- 支持的Linux版本：
    - Ubuntu 18.04 LTS及以上版本
    - CentOS 7.0及以上版本
- go版本：
    - go 1.20及以上版本
    - 未在go 1.19及以下版本上进行测试

## 运行方式
- 完善 ./config/alibaba.toml配置文件中未填写的必要信息
- 项目根路径下执行：
```bash
./run.sh
```

## 代码介绍

### 代码结构介绍
- **main.go**: 主入口文件
- **business/**: 业务逻辑模块
    - **engine/**: 数据流转的主干逻辑
    - **exit/**: 进程退出逻辑
    - **filter/**: 输入音频处理、过滤逻辑
        - 对原始输入音频进行vad处理，主要是判断是否为人声
        - 将处理后的音频输送到队列
    - **stt/**: stt业务模块
      - **stt.go**: 当前模块对外暴露的抽象层
      - **ali**: 访问阿里stt的业务实现层（ps:在同级目录下可以扩展其他vendor的stt）
          - **conn.go**: 访问阿里的stt用的是阿里的go-sdk，这里把每个sdk返回的client实例定义为一个业务层的连接对象
          - **connpool.go**: 连接池逻辑，可省去每次请求stt建立连接的耗时
    - **llm/**: llm业务模块：
      - **llm.go**: 当前模块对外暴露的抽象层
      - **qwen.go**: 访问阿里千问大模型的业务实现层（ps:在同级目录下可以扩展其他vendor的大模型）
      - **common/** llm下各个vendor访问的公共代码块
        - **clause/**: 当llm流式返回时，这里包含了对返回流的分句的策略（例如：按照标点分句或者不处理）
        - **dialogctx/**: 当llm需要记录上下文信息时，这个模块负责对上下文信息的管理
    - **tts/**: tts业务模块
        - **tts.go**: 当前模块对外暴露的抽象层
        - **httpsender.go**: 基于http访问方式的tts的抽象子层（ps:阿里cosy访问方式则是基于websocket访问，在同级目录下可以扩展其他访问方式的抽象子层）
        - **vendor/**: 访问各vendor的tts的业务实现层
            - **ali/**: 访问阿里tts的业务实现逻辑（ps:在同级目录下可以扩展其他vendor业务实现逻辑）
            - **sentence.go**: 主要是整合同一个sentence下的各个segment的音频
    - **rtc/**: 使用rtc推拉流、发送datastream等的业务逻辑层
    - **sentencelifecycle/**: sentence生命周期管理（属于业务层的全局的依赖模块）
    - **workerid**: 唯一表示一个agent进程，用于打印日志（属于业务层的全局的依赖模块）
- **clients/**: http客户端请求响应模块，与业务无关，只是单纯的http接口访问逻辑
  - **alitts/**: 请求阿里tts的接口
    - **client.go**: 初始化访问阿里tts域名的client，配置http连接池参数（降低并发请求情况下的同步建立概率）、发起预热连接（提前建立长连接，避免首次请求时同步建连）
    - **streamask.go**: 访问流式接口的实现（ps:在同级目录下可以扩展访问其他接口的实现，例如非流式访问）
  - **qwen/**: 请求阿里千问大模型的接口
    - **类比alitts的注解**
- **config**: 初始化配置模块
- **pkg/**: 独立的公共依赖库
  - **agora-go-sdk**: 声网rtc、vad等的go-sdk
  - **alibaba**: 与阿里有关的公共依赖
    - **speech**: 与阿里语音（识别/合成）有关的公共依赖
      - **token.go**: 初始化访问阿里语音服务的token
  - **httputil**: http请求的一些公共访问函数
  - **logger**: 日志包


### 概念/功能介绍
#### sentence
- sentence在代码中的意思是一个完整的句子，一段音频、一个文本句子都可以叫sentence
- 同一时间，程序中只会处理一个sentence（这个由「打断功能」来保障）
#### sid、sgid
sid是指对一个sentence的唯一标识。如果对几个sentence划分到一个group，那么这个group的唯一标识就是sgid
#### 打断功能
- 含义：打断是指将运行中的sentence处理逻辑给停止（实际是指stt下游的逻辑，在stt这一层时可以多个sentence生命走起并存的）
- 触发条件：只要是filter模块识别除了新的sentence，就会立即调用打断函数，触发打断
#### 断句、并句功能
为了描述这两个功能的含义，举一个例子：
假设此时filter模块返回了两段sentence音频分别为：sentence_i 和 sentence_i+1，因此它们是时间上连续的，并且 sentence_i+1 在时间上是靠后的。
如果，在接收到 sentence_i+1 的head chunk时，sentence_i的处理逻辑已经走到了往rtc发送回复音频阶段了，那么此时会触发「断句」，否则，会触发「并句」逻辑。
- 断句：就是 sentence_i+1 到来时，立即打断 sentence_i 的处理逻辑，然后进入 sentence_i+1 的处理逻辑
- 并句：就是 sentence_i+1 到来时，立即打断 sentence_i 的处理逻辑，然后将 sentence_i 的stt识别文本和 sentence_i+1 的stt识别文本合并到一起，
然后一起作为一个sentence的逻辑来处里。此时，sid是 sentence_i+1 的sid，sgid是 sentence_i 的sgid（根据递推关系，如果sentence_i也是和sentence_i-1并句后的结果，那么sgid就是sentence_i-1的sid）
- 断句和并句的一个共同点：都会触发「打断功能」
#### stt采用多实例模式
一般stt的demo代码都是提供了一个实例（一个连接），然后将音频传入到stt中，由stt来负责语音分句，并由回调函数来返回识别到的句子，但是当前项目并不是。
- 区别：
  - 当前项目将分句功能交给了filter模块（识别一句话音频的head、tail等）
  - 当前项目会将每个sentence音频交给一个独立的stt实例（连接）来处理
- 意义：
  - 并发性：sentence_i 和 sentence_i+1可以并发地被stt服务处理，而不需要等待。
  - 可用性：因为单个stt实例是存在空闲超时时间的（阿里的是10s），如果10s内没有音频输入到stt，那么实例的会被sdk销毁，此时发送音频就会报错
#### stt连接池
（当前项目中stt使用的是阿里的sdk，因此把每个client实例称为一个connection）
- 含义：自定义的连接池，用于定时创建空闲连接，这里没有设计回收连接的逻辑
- 意义：确保每次需要访问stt服务的时候都可以直接从pool中取出连接，而不需要临时地区创建连接进入引入耗时
#### segment
- 一个sentence可以划分为多个segment。例如llm流式返回的结果会被聚合成segment交给tts处理，segment会通过sid关联到sentence
#### tts并发合成语音
- 含义：因为一个sentence在llm返回的结果以多个segment的形式交给tts，所以tts会独立地、并发地处理这些segment，然后异步地将各个segment的返回的音频结果合并起来
- 优点：降低了 多个segment在 语音合成时因为同步等待所带来的耗时
- 注意：并发度够用就行，调高了会触发限流。经验值：2
