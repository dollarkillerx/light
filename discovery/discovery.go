package discovery

// Discovery 服务发现
type Discovery interface {
	Discovery(serName string) ([]string, error)
}
