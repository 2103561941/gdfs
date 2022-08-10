package directorytree

var (
	tree DirectoryTree
)

// 目录树结构
type DirectoryTree struct {
	root  *node
}

func Setup() error {
	tree = DirectoryTree{}
	return nil
}

// namenode 中目录结构文件需要保存的信息
type node struct {
	path string // 当前路径、文件名
	file filemeta
	childNodes []*node
}

// 文件元数据，用于展示文件的相关信息
type filemeta struct {
	fileType int
	fileName string
	fileMD5  string
	fileSize float64 // 单位 MB
	data []dataNode // 保存文件存放的datanode节点
}

// namenode 中记录的 dataNode 的相关信息
type dataNode struct {
	key string // 唯一标识符
}
