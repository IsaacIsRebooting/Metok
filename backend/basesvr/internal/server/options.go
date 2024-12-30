package server

type Params struct {
	addr string
	// dbShardingTableConfig map[string]conf.DomainShardingTableConfig
}

type Option func(*Params)

func WithAddr(addr string) Option {
	return func(p *Params) {
		p.addr = addr
	}
}
