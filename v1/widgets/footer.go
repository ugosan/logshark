package widgets

import (
	"image"

	. "github.com/gizak/termui/v3"
	"github.com/ugosan/logshark/v1/logging"
)

var (
	logs = logging.GetManager()
)

type Footer struct {
	Block
	Text      string
	TextStyle Style
	WrapText  bool
}

func NewFooter() *Footer {
	return &Footer{
		Block:     *NewBlock(),
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
	}
}

func (self *Footer) Draw(buf *Buffer) {
	// no padding
	self.Block.Inner = image.Rect(
		self.Block.Min.X,
		self.Block.Min.Y,
		self.Block.Max.X,
		self.Block.Max.Y,
	)
	
	self.Block.Draw(buf)

	cells := ParseStyles(self.Text, self.TextStyle)
	
	runes := []rune(" ")

	cellsLength := len(cells)
	for i := cellsLength; i < self.Block.Max.X; i++ {
		cells = append(cells, Cell{runes[0], self.TextStyle})
	}

	for _, cx := range BuildCellWithXArray(cells) {
		x, cell := cx.X, cx.Cell
		buf.SetCell(cell, image.Pt(x, 0).Add(self.Inner.Min))
	}
}
