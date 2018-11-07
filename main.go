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

func (b *Block) Get() string {
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

func RenderBlocks(blocks []Block) {
    for _, block := range blocks {
        block.Render()
    }
}

func DumpToFile(blocks []Block) {
    var blockAttrs string

    for _, block := range blocks {

        blockAttrs += block.Get()
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

func RenderButtons(blocks []Button) {
    for _, block := range blocks {
        block.Render()
    }
}

func ActivateButton(buttons []Button, index int) {
    for i := 0; i < 7; i++ {
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

    workspace := ui.NewPar("")
    workspace.Height = ui.TermHeight()
    workspace.BorderLabel = "Workspace"
    workspace.BorderLabelFg = ui.ColorWhite

    controls := ui.NewPar("")
    controls.Height = ui.TermHeight()
    controls.BorderLabel = "Controls"
    controls.BorderLabelFg = ui.ColorWhite

    blocks := make([]Block, 0, 100)
    buttons := make([]Button, 9, 9)

    buttons[0].Init(ui.TermWidth()-22, 18, 3, 15, ui.ColorRed, " Red ")
    buttons[1].Init(ui.TermWidth()-22, 21, 3, 15, ui.ColorGreen, " Green ")
    buttons[2].Init(ui.TermWidth()-22, 24, 3, 15, ui.ColorYellow, " Yellow ")
    buttons[3].Init(ui.TermWidth()-22, 27, 3, 15, ui.ColorBlue, " Blue ")
    buttons[4].Init(ui.TermWidth()-22, 30, 3, 15, ui.ColorMagenta, " Magenta ")
    buttons[5].Init(ui.TermWidth()-22, 33, 3, 15, ui.ColorCyan, " Cyan ")
    buttons[6].Init(ui.TermWidth()-22, 36, 3, 15, ui.ColorWhite|ui.AttrBold, " White ")

    buttons[7].Init(ui.TermWidth()-22, 10, 3, 15, ui.ColorRed, " New Block ")
    buttons[8].Init(ui.TermWidth()-22, 13, 3, 15, ui.ColorRed, " Move or Del ")

    buttons[0].Active = true
    buttons[7].Active = true
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

    ui.Render(ui.Body)
    RenderButtons(buttons)

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
        ui.Render(ui.Body)
        RenderButtons(buttons)
        RenderBlocks(blocks)

    })

    ui.Handle("t", func(ui.Event) {
        newMode = !newMode
        moveMode = !moveMode
        buttons[7].Toggle()
        buttons[8].Toggle()

        if selectedBlock != nil {
            selectedBlock.Focused = false
            selectedBlock = nil
            selectedIndex = -1
        }

    })

    ui.Handle("<Delete>", func(ui.Event) {
        if selectedIndex >= 0 {
            blocks = append(blocks[:selectedIndex], blocks[selectedIndex+1:]...)
            selectedIndex = -1
            selectedBlock = nil
        }
    })

    ui.Handle("<Keyboard>", func(e ui.Event) {
        i, err := strconv.Atoi(e.ID)

        if err == nil && i >= 1 && i <= 7 {

            selectedColor = buttons[i-1].Color
            ActivateButton(buttons, i-1)

            if selectedBlock != nil {
                selectedBlock.Color = buttons[i-1].Color
            }

        }

        ui.Clear()
        ui.Render(ui.Body)
        RenderButtons(buttons)
        RenderBlocks(blocks)
    })

    ui.Handle("q", func(ui.Event) {
        ui.StopLoop()
    })

    ui.Handle("<Enter>", func(ui.Event) {
        DumpToFile(blocks)
    })

    ui.Loop()
}
