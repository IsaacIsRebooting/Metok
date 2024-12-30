package gofer

type GroupOption func(*Group)

func UserErrorGroup() GroupOption {
	return func(g *Group) {
		g.isErrorGroup = true
	}
}

func WithUsableGroup(num int) GroupOption {
	return func(g *Group) {
		g.numG = num
	}
}

func WithWaitQueue(size int) GroupOption {
	return func(g *Group) {
		g.queueSize = size
	}
}
