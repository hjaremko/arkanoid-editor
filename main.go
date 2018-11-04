package main

import (
    // "fmt"
    ui "github.com/gizak/termui"
    "io/ioutil"
    "strconv"
)

type Block struct {
    TermBlock *ui.Block
    Color     ui.Attribute
    Breakable int
    Focused   bool
}

func (block *Block) SetPosition(x, y int) {
    block.TermBlock.X = x
    block.TermBlock.Y = y
}

func (block *Block) Render() {
    if block.Focused {
        block.TermBlock.BorderFg = ui.AttrReverse | block.Color
    } else {
        block.TermBlock.BorderFg = block.Color
    }

    ui.Render(block.TermBlock)
}

func (block *Block) Dump() string {
    output := strconv.Itoa(block.TermBlock.Y) + " "
    output += strconv.Itoa(block.TermBlock.X) + " "
    output += strconv.Itoa(block.TermBlock.Width) + " "
    output += strconv.Itoa(block.TermBlock.Height) + " "
    output += strconv.Itoa(int(block.Color) - 1)
    output += " 1\n"

    return output
}

func renderBlocks(blocks []Block) {
    for _, block := range blocks {
        block.Render()
    }
}

func (block *Block) Intersects(x int, y int) bool {
    return (y >= block.TermBlock.Y && y < (block.TermBlock.Y+block.TermBlock.Height) && x >= block.TermBlock.X && x < (block.TermBlock.X+block.TermBlock.Width))
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

type ColorBlock struct {
    Block  *ui.Par
    Color  ui.Attribute
    Active bool
}

func (block *ColorBlock) Render() {
    if block.Active {
        block.Block.BorderFg = block.Color
    } else {
        block.Block.BorderFg = ui.ColorWhite
    }

    ui.Render(block.Block)
}

func (block *ColorBlock) Toggle() {
    block.Active = !block.Active
}

func (block *ColorBlock) Init(x, y, height, width int, color ui.Attribute, colorName string) {
    if block.Block == nil {
        block.Block = ui.NewPar(colorName)
        block.Block.X = x
        block.Block.Y = y
        block.Block.Height = height
        block.Block.Width = width
        block.Block.BorderFg = color
        block.Color = color
    }
}

func renderColorButtons(blocks []ColorBlock) {
    for _, block := range blocks {
        block.Render()
    }
}

func activateColorButton(buttons []ColorBlock, index int) {
    for i, _ := range buttons {
        if i == index {
            buttons[i].Active = true
        } else {
            buttons[i].Active = false
        }
    }
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
    colorButtons := make([]ColorBlock, 7, 7)

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
    renderColorButtons(colorButtons)

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

            blocks = append(blocks, newBBlock)

        } else if moveMode {
            found := false

            for i, block := range blocks {
                if block.Intersects(mouseX, mouseY) {

                    if selectedBlock != nil {
                        selectedBlock.Focused = false
                    }

                    blocks[i].Focused = true
                    selectedBlock = &blocks[i]
                    selectedIndex = i

                    found = true
                    break
                }
            }

            if selectedBlock != nil && !found {
                selectedBlock.SetPosition(mouseX, mouseY)
            }
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
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
        renderColorButtons(colorButtons)
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
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("1", func(ui.Event) {
        selectedColor = ui.ColorRed
        activateColorButton(colorButtons, 0)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[0].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("2", func(ui.Event) {
        selectedColor = ui.ColorGreen
        activateColorButton(colorButtons, 1)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[1].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("3", func(ui.Event) {
        selectedColor = colorButtons[2].Color
        activateColorButton(colorButtons, 2)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[2].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("4", func(ui.Event) {
        selectedColor = colorButtons[3].Color
        activateColorButton(colorButtons, 3)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[3].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("5", func(ui.Event) {
        selectedColor = colorButtons[4].Color
        activateColorButton(colorButtons, 4)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[4].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("6", func(ui.Event) {
        selectedColor = colorButtons[5].Color
        activateColorButton(colorButtons, 5)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[5].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
        renderBlocks(blocks)

    })

    ui.Handle("7", func(ui.Event) {
        selectedColor = colorButtons[6].Color
        activateColorButton(colorButtons, 6)

        if selectedBlock != nil {
            selectedBlock.Color = colorButtons[6].Color
        }

        ui.Clear()
        ui.Render(ui.Body, newBlockButton, moveButton)
        renderColorButtons(colorButtons)
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
