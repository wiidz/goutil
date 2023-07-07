package imgHelper

type Size struct {
	Width  float64
	Height float64
}

type Position struct {
	X float64
	Y float64
}

type CoverImgInterface interface {
	GetLocalFilePath() string
	GetSize() *Size
	GetPosition() *Position
	GetRoundCorner() int
	GetCropCircle() bool
}

type LocalCoverImg struct {
	LocalFilePath string // 本地文件地址
	Size          *Size
	Position      *Position
	RoundCorner   int
	CropCircle    bool
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
func (t *LocalCoverImg) GetRoundCorner() int {
	return t.RoundCorner
}
func (t *LocalCoverImg) GetCropCircle() bool {
	return t.CropCircle
}

type NetworkCoverImg struct {
	NetworkURL    string // 网络文件地址
	LocalFilePath string // 本地文件地址
	Size          *Size
	Position      *Position
	RoundCorner   int
	CropCircle    bool
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

func (t *NetworkCoverImg) GetRoundCorner() int {
	return t.RoundCorner
}
func (t *NetworkCoverImg) GetCropCircle() bool {
	return t.CropCircle
}
