package imgHelper

type Size struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

type CoverImgInterface interface {
	GetLocalFilePath() string
	GetSize() *Size
	GetPosition() *Position
}

type LocalCoverImg struct {
	LocalFilePath string // 本地文件地址
	Size          *Size
	Position      *Position
}

func (t *LocalCoverImg) GetLocalFilePath() string {
	return t.LocalFilePath
}
func (t *LocalCoverImg) GetSize() *Size {
	return t.Size
}
func (t *LocalCoverImg) GetPosition() *Position {
	return t.Position
}

type NetworkCoverImg struct {
	NetworkURL    string // 网络文件地址
	LocalFilePath string // 本地文件地址
	Size          *Size
	Position      *Position
}

func (t *NetworkCoverImg) GetLocalFilePath() string {
	return t.LocalFilePath
}
func (t *NetworkCoverImg) GetSize() *Size {
	return t.Size
}
func (t *NetworkCoverImg) GetPosition() *Position {
	return t.Position
}
