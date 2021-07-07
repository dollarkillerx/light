# Design Light RPC

## 分层抽象
- 传输层
- 协议层
- 服务层

## 服务端设计
### 基础服务定义 与内部注册
我们以一个 helloWorld 服务来演示

**1.1 定义服务端**
```go
type server struct {}

type HelloWorldRequest struct {
	Name string
}

type HelloWorldResponse struct {
	Msg string
}

func (s *server) HelloWorld(ctx *light.Context, req *HelloWorldRequest, resp *HelloWorldResponse) error {
	resp.Msg = fmt.Sprintf("hello world by: %s", req.Name)
	return nil
}
```
**1.2 服务注册**
服务管理:
```go
type Server struct {
	serviceMap map[string]*service
}
```
具体每个服务:
```go
type service struct {
	name       string                 // server name
	refVal     reflect.Value          // server reflect value
	refType    reflect.Type           // server reflect type
	methodType map[string]*methodType // server method
}
```
具体每个方法:
```go
type methodType struct {
	method       reflect.Method
	RequestType  reflect.Type
	ResponseType reflect.Type
}
```

绑定具体方法
```go
// constructionMethods Get specific method
func constructionMethods(typ reflect.Type) (map[string]*methodType, error) {
	methods := make(map[string]*methodType)
	for idx := 0; idx < typ.NumMethod(); idx++ {
		method := typ.Method(idx)
		mType := method.Type
		mName := method.Name

		if !utils.IsPublic(mName) {
			return nil, errors.New("Registered non-public service")
		}

		// 默认是4个
		if mType.NumIn() != 4 { // func(*server.MethodTest, *light.Context, *server.MethodTestReq, *server.MethodTestResp) error
			continue
		}

		// 检验它第一个参数是否是ctx
		ctxType := mType.In(1)
		if !(ctxType.Elem() == typeOfContext) {
			continue
		}

		// request 参数检查
		requestType := mType.In(2)
		if requestType.Kind() != reflect.Ptr {
			continue
		}

		if !utils.IsPublicOrBuiltinType(requestType) {
			continue
		}

		// response 参数检查
		responseType := mType.In(3)
		if responseType.Kind() != reflect.Ptr {
			continue
		}

		if !utils.IsPublicOrBuiltinType(responseType) {
			continue
		}

		// 校验返回参数
		if mType.NumOut() != 1 {
			continue
		}

		returnType := mType.Out(1)
		if returnType != typeOfError {
			continue
		}

		methods[mName] = &methodType{
			method:       method,
			RequestType:  requestType,
			ResponseType: responseType,
		}
	}

	if len(methods) == 0 {
		return nil, errors.New("No service is available, or provide service is not open")
	}

	return methods, nil
}
```
构造服务:
```go
func newService(server interface{}, serverName string, useName bool) (*service, error) {
	ser := &service{
		refVal:  reflect.ValueOf(server),
		refType: reflect.TypeOf(server),
	}

	sName := reflect.Indirect(ser.refVal).Type().Name()
	if !utils.IsPublic(sName) {
		return nil, pkg.ErrNonPublic
	}

	if useName {
		if serverName == "" {
			return nil, errors.New("Server Name is null")
		}

		sName = serverName
	}

	ser.name = sName
	methods, err := constructionMethods(ser.refType)
	if err != nil {
		return nil, err
	}
	ser.methodType = methods

	return ser, nil
}
```
调用服务方法
```go
// call 方法调用
func (s *service) call(ctx *light.Context, mType *methodType, request, response reflect.Value) (err error) {
	// recover 捕获堆栈消息
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[painc service internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mType.method.Name, request.Interface(), buf)
			log.Println(err)
		}
	}()

	fn := mType.method.Func
	returnValue := fn.Call([]reflect.Value{s.refVal, reflect.ValueOf(ctx), request, response})
	errInterface := returnValue[0].Interface()
	if errInterface != nil {
		return errInterface.(error)
	}

	return nil
}
```

服务注册到manager
```go
func (s *Server) Register(server interface{}) error {
	return s.register(server, "", false)
}

func (s *Server) RegisterName(server interface{}, serverName string) error {
	return s.register(server, serverName, true)
}

func (s *Server) register(server interface{}, serverName string, useName bool) error {
	ser, err := newService(server, serverName, useName)
	if err != nil {
		return err
	}

	s.serviceMap[ser.name] = ser
	return nil
}
```

### 1.3编码解码
codes 目录下创建 序列化工具

实现以下interface
```go
type Serialization interface {
    Encode(i interface{}) ([]byte, error)
    Decode(data []byte, i interface{}) error
}
```
编写manager进行管理
```go
type serializationManager struct {
    codes map[SerializationType]Serialization
}

var Manager = &serializationManager{
    codes: map[SerializationType]Serialization{},
}

type SerializationType byte

const (
    CodeJson SerializationType = iota
    CodeMsgPack
)

func (m *serializationManager) register(key SerializationType, code Serialization) {
    m.codes[key] = code
}

func (m *serializationManager) Get(key SerializationType) (Serialization, bool) {
    code, ex := m.codes[key]
    return code, ex
}
```

具体实现: 放到 serialization_plugin 下

###  1.4 compressor 压缩
基础逻辑同上serialization

### 1.5 协议设计
``` 
	/**
	crc32	:	total	:	offset	: magicNumberSize: magicNumber: serverNameSize : serverMethodSize:  respType : compressorType: serializationType : serverName : serverMethod :  payload
	4 		:	4 		: 	4 	    :     4          :     xxxx   :       4        :         4        :     1    :        1      :          1        : xxx        :      xxx     :  xxx
	*/
```
具体编码解码: protocol/protocaol.go
