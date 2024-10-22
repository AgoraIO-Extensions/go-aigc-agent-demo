
## 运行环境
- 支持的Linux版本：
    - Ubuntu 18.04 LTS及以上版本
- go版本：
    - go 1.20及以上版本

## 运行方式
- 完善 ./config/chat-robot.toml配置文件中未填写的必要信息
- 如果是 Ubuntu 20.04/22.04/24.04 ，则在项目根路径下执行：
```bash
./run.sh
```
- 如果是 Ubuntu 18.04 ，则在项目根路径下执行：
```bash
./run_ubuntu1804.sh
```

## 代码介绍

### 代码结构介绍
- **main.go**: 主入口文件
- **business/**: 业务逻辑模块
    - **aigcCtx**: 信息链控制模块
      - **sentence**:用于记录每一个sentence的一些元数据
      - **ctx.go、interrupt.go**: 信息链控制，打断逻辑相关
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
  - **httputil**:
    - **client.go**: 使用自定义的client，主要是解决首次建连的耗时问题
    - **request.go**: 基于自定义的client发起post请求
  - **logger**: 日志包


### 概念/功能介绍
#### sentence
- sentence在代码中的意思是一个完整的句子，一段音频、一个文本句子都可以叫sentence
- 同一时间，stt下游的模块只会处理一个sentence（这个由「打断功能」来保障）
#### sid、sgid
sid是指对一个sentence的唯一标识。如果对几个sentence划分到一个group，那么这个group的唯一标识就是sgid。打印日志时会加上这两个标识。
#### 打断功能
- 含义：后来的sentence会执行打断函数，从而停止先来的sentence的执行。
- 触发条件：目前有两个可选触发条件：1.识别到新的sentence音频到来时立即打断。2.stt识别到新的sentence文本第一个字后立即打断。
#### 音频分句
- 含义：filter输出的音频chunk带有(MuteToSpeak/Speaking/SpeakToMute等)标记，依据这些标记可以判断filter对原始音频的分句，例如一段完整音频被分解为「sAudio1,sAudio2,sAudio3...」
- 意义：音频分句后，stt可以并发消费多段音频
#### 音频分组
- 含义：音频分句的「sAudio1,sAudio2,sAudio3...」会依据一定的规则被分组，例如 「sAudio1,sAudio2」被分为一组，「sAudio3」被分为一组。这个就叫音频分组。
- 性质：分组后的音频的sgid是相同的；属于同一音频组的音频在经过stt转话为文本后也属于同一组；同一组的文本会拼接到一起作为一个整体传输给llm等下游模块
- 意义：原始说话人声可能存在停顿，停顿就会触发音频分句，例如分为「sAudio1,sAudio2」：
  - 停顿很短：那么希望是「sAudio1,sAudio2」的处理应该按一句话来处理，所以需要分到同一组。
  - 停顿太长：那么就应该是说话者已经不关心 sAudio1，所以就不会把 sAudio1和sAudio2分到一组。
#### 文本拼接
-含义：在「音频分组」中已经提到了“属于同一音频组的音频在经过stt转话为文本后也属于同一组；同一组的文本会拼接到一起作为一个整体传输给llm等下游模块”，这个就是文本拼接。
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
#### http预连接
pkg/httputil/client.go中的NewClient函数如果返回的是基于自定义transport的client，那么该client的首次http/https请求将省去建立连接的耗时。
此外，该client还支持默认client所支持的连接池功能，以及http1.1升级http2等功能

## 常见错误
#### 阿里stt/tts触发限流报错：
如果阿里的stt/tts的出现了类似「Gateway:TOO_MANY_REQUESTS:Too many requests!」这样的错误信息，那么说明被限流了，可能的原因： 创建连接实例过于频繁或者是请求频率太高。
解决的方法：要么降低响应的建连频率/请求并发数，要么调大客户账号在阿里云配置的对应服务的并发数（参考链接：https://nls-portal.console.aliyun.com/servicebuy）。

## 发版日志
v2.0：支持阿里/微软   
v2.1: 优化stt/tts初始化逻辑   
v2.2: 微软语音服务兼容 ubuntu 18.04、20.04、22.04、24.04   
v2.3: 增加vad参数，控制人声识别敏感度   
v2.4: 1.打断时机可配置，支持在filter或者是stt模块进行打断，以适应不同的客户场景需求。2. sentence分组现在支持两种策略：dependOnRTCSend/dependOnTime   
v2.5: 1. engine下的数据模块划分更清晰；2. sid、sgid等sentence元信息隐藏到了ctx中，会在打日志时自动补充到日志tag中；3. 支持两种打断策略（时间间隔是否达到一个阈值 和 sentence是否直行道往rtc发送的阶段）4. 往rtc发音频的算法改为了临牌桶算法，上层业务将不需要关心发送   
v2.6: 1. 优化了ctx的设计；2. rtc发送音频阶段增加了动态等待策略，可以解决部分噪音场景的错误打断问题，同时也更好的兼容了用户说话短暂停顿后继续说话的场景