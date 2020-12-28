
package widgets

import (
	"image"
	. "github.com/gizak/termui/v3"
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
	
	// fills the remaining space so it takes the whole row
	for i := len(self.Text); i < self.Block.Max.X; i++ {
    self.Text += " "
	}

	self.Block.Draw(buf)
	
	cells := ParseStyles(self.Text, self.TextStyle)
	if self.WrapText {
		cells = WrapCells(cells, uint(self.Inner.Dx()))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = TrimCells(row, self.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
	}
}