package geecache

// PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 的 Get 方法用于在对应 group 中查找缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)	// 回调函数
}