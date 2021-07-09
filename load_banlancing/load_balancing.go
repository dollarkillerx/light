package load_banlancing

type LoadBalancing interface {
	InitBalancing(sers []string)
	GetService() string
}
