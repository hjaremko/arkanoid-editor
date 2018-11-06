package main

import (
    // "fmt"
    ui "github.com/gizak/termui"
    "io/ioutil"
    "strconv"
)

// Block

type Block struct {
    TermBlock *ui.Block
    Color     ui.Attribute
    Breakable int
    Focused   bool
}

func (b *Block) SetPosition(x, y int) {
    b.TermBlock.X = x
    b.TermBlock.Y = y
}

func (b *Block) Render() {
    if b.Focused {
        b.TermBlock.BorderFg = ui.AttrReverse | b.Color
    } else {
        b.TermBlock.BorderFg = b.Color
    }

    ui.Render(b.TermBlock)
}

func (b *Block) Dump() string {
    output := strconv.Itoa(b.TermBlock.Y) + " "
    output += strconv.Itoa(b.TermBlock.X) + " "
    output += strconv.Itoa(b.TermBlock.Width) + " "
    output += strconv.Itoa(b.TermBlock.Height) + " "
    output += strconv.Itoa(int(b.Color) - 1)
    output += " 1\n"

    return output
}

func (b *Block) Intersects(x int, y int) bool {
    return (y >= b.TermBlock.Y && y < (b.TermBlock.Y+b.TermBlock.Height) && x >= b.TermBlock.X && x < (b.TermBlock.X+b.TermBlock.Width))
}

func (b *Block) IntersectsBlock(other Block) bool {
    aX1 := b.TermBlock.X
    aX2 := b.TermBlock.X + b.TermBlock.Width
    aY1 := b.TermBlock.Y
    aY2 := b.TermBlock.Y + b.TermBlock.Height

    bX1 := other.TermBlock.X
    bX2 := other.TermBlock.X + other.TermBlock.Width
    bY1 := other.TermBlock.Y
    bY2 := other.TermBlock.Y + other.TermBlock.Height

    xOverlap := ValueInRange(aX1, bX1, bX2) || ValueInRange(bX1, aX1, aX2)
    yOverlap := ValueInRange(aY1, bY1, bY2) || ValueInRange(bY1, aY1, aY2)

    return xOverlap && yOverlap
}

//~Block

func ValueInRange(value, min, max int) bool {
    return ((value >= min) && (value < max))
}

func renderBlocks(blocks []Block) {
    for _, block := range blocks {
        block.Render()
    }
}

func dump(blocks []Block) {
    var blockAttrs string

    for _, block := range blocks {

        blockAttrs += block.Dump()
    }

    err := ioutil.WriteFile("level1.txt", []byte(blockAttrs), 0644)

    if err != nil {
        panic(err)
    }
}

//Button

type Button struct {
    Block  *ui.Par
    Color  ui.Attribute
    Active bool
}

func (b *Button) Toggle() {
    b.Active = !b.Active
}

func (b *Button) Render() {
    if b.Active {
        b.Block.BorderFg = b.Color
    } else {
        b.Block.BorderFg = ui.ColorWhite
    }

    ui.Render(b.Block)
}

func (b *Button) Init(x, y, height, width int, color ui.Attribute, colorName string) {
    if b.Block == nil {
        b.Block = ui.NewPar(colorName)
        b.Block.X = x
        b.Block.Y = y
        b.Block.Height = height
        b.Block.Width = width
        b.Block.BorderFg = color
        b.Color = color
    }
}

//~ColorButton

func renderButtons(blocks []Button) {
    for _, block := range blocks {
        block.Render()
    }
}

func activateButton(buttons []Button, index int) {
    for i, _ := range buttons {
        if i == index {
            buttons[i].Active = true
        } else {
            buttons[i].Active = false
        }
    }
}

func FindIntersecting(blocks []Block, mouseX, mouseY int) (index int, found bool) {
    for i, block := range blocks {
        if block.Intersects(mouseX, mouseY) {
            return i, true
        }
    }

    return -1, false
}

func FindIntersectingBlock(blocks []Block, newBlock Block) (index int, found bool) {
    for i, block := range blocks {
        if block.IntersectsBlock(newBlock) {
            return i, true
        }
    }

    return -1, false
}

func main() {
    if err := ui.Init(); err != nil {
        panic(err)
    }

    defer ui.Close()

    newBlockButton := ui.NewPar(" New Block ")
    newBlockButton.Height = 3
    newBlockButton.Width = 15
    newBlockButton.X = ui.TermWidth() - 22
    newBlockButton.Y = 10
    newBlockButton.BorderFg = ui.ColorRed

    moveButton := ui.NewPar(" Move or Del")
    moveButton.Height = 3
    moveButton.Width = 15
    moveButton.BorderLabelFg = ui.ColorWhite
    moveButton.X = ui.TermWidth() - 22
    moveButton.Y = 13

    workspace := ui.NewPar("")
    workspace.Height = ui.TermHeight()
    workspace.BorderLabel = "Workspace"
    workspace.BorderLabelFg = ui.ColorWhite

    controls := ui.NewPar("")
    controls.Height = ui.TermHeight()
    controls.BorderLabel = "Controls"
    controls.BorderLabelFg = ui.ColorWhite

    blocks := make([]Block, 0, 100)
    colorButtons := make([]Button, 7, 7)

    colorButtons[0].Init(ui.TermWidth()-22, 18, 3, 15, ui.ColorRed, " Red ")
    colorButtons[1].Init(ui.TermWidth()-22, 21, 3, 15, ui.ColorGreen, " Green ")
    colorButtons[2].Init(ui.TermWidth()-22, 24, 3, 15, ui.ColorYellow, " Yellow ")
    colorButtons[3].Init(ui.TermWidth()-22, 27, 3, 15, ui.ColorBlue, " Blue ")
    colorButtons[4].Init(ui.TermWidth()-22, 30, 3, 15, ui.ColorMagenta, " Magenta ")
    colorButtons[5].Init(ui.TermWidth()-22, 33, 3, 15, ui.ColorCyan, " Cyan ")
    colorButtons[6].Init(ui.TermWidth()-22, 36, 3, 15, ui.ColorWhite|ui.AttrBold, " White ")

    colorButtons[0].Active = true
    newMode := true
    moveMode := false

    var selectedBlock *Block
    var selectedIndex int
    selectedColor := ui.ColorRed

    ui.Body.AddRows(
        ui.NewRow(
            ui.NewCol(9, 0, workspace),
            ui.NewCol(3, 0, controls)))

    ui.Body.Align()

    ui.Render(ui.Body, newBlockButton, moveButton)
    renderButtons(colorButtons)

    ui.Handle("<MouseLeft>", func(e ui.Event) {
        mouseX := e.Payload.(ui.Mouse).X
        mouseY := e.Payload.(ui.Mouse).Y

        if newMode {

            newBlock := ui.NewBlock()
            newBlock.Height = 3
            newBlock.Width = 6
            newBlock.X = mouseX
            newBlock.Y = mouseY

            newBBlock := Block{
                newBlock,
                selectedColor,
                1,
                false,
            }

            _, found := FindIntersectingBlock(blocks, newBBlock)

            if !found {
                blocks = append(blocks, newBBlock)
            }

        } else if moveMode {
            foundIndex, found := FindIntersecting(blocks, mouseX, mouseY)

            if found {
                if selectedBlock != nil {
                    selectedBlock.Focused = false
                }

                blocks[foundIndex].Focused = true
                selectedBlock = &blocks[foundIndex]
                selectedIndex = foundIndex

            } else if selectedBlock != nil {
                selectedBlock.SetPosition(mouseX, mouseY)
            }
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("t", func(ui.Event) {
        newMode = !newMode
        moveMode = !moveMode

        if newMode {
            newBlockButton.BorderFg = ui.ColorRed
        } else {
            newBlockButton.BorderFg = ui.ColorWhite
        }

        if moveMode {
            moveButton.BorderFg = ui.ColorRed
        } else {
            moveButton.BorderFg = ui.ColorWhite
        }

        if selectedBlock != nil {
            selectedBlock.Focused = false
            selectedBlock = nil
            selectedIndex = -1
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("<Delete>", func(ui.Event) {
        if selectedIndex >= 0 {
            blocks = append(blocks[:selectedIndex], blocks[selectedIndex+1:]...)
            selectedIndex = -1
            selectedBlock = nil
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("1", func(ui.Event) {
        selectedColor = ui.ColorRed
        activateButton(colorButtons, 0)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[0].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("2", func(ui.Event) {
        selectedColor = ui.ColorGreen
        activateButton(colorButtons, 1)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[1].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("3", func(ui.Event) {
        selectedColor = colorButtons[2].Color
        activateButton(colorButtons, 2)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[2].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("4", func(ui.Event) {
        selectedColor = colorButtons[3].Color
        activateButton(colorButtons, 3)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[3].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("5", func(ui.Event) {
        selectedColor = colorButtons[4].Color
        activateButton(colorButtons, 4)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[4].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("6", func(ui.Event) {
        selectedColor = colorButtons[5].Color
        activateButton(colorButtons, 5)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[5].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("7", func(ui.Event) {
        selectedColor = colorButtons[6].Color
        activateButton(colorButtons, 6)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[6].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("q", func(ui.Event) {
        ui.StopLoop()
    })

    ui.Handle("<Enter>", func(ui.Event) {
        dump(blocks)
    })

    ui.Loop()
}
