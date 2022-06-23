package imgHelper

type Size struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

type CoverImg struct {
	LocalFilePath string // 本地文件地址
	NetworkURL    string // 网络文件地址
	Size          *Size
	Position      *Position
}
