package ui2d

import (
	"fmt"

	"github.com/ahmadfarhanstwn/rpg/game-logic"
	"github.com/veandco/go-sdl2/sdl"
)


func (ui *ui) getInventoryRect() *sdl.Rect {
	invWidth := int32(float32(ui.winWidth)*.60)
	invHeight := int32(float32(ui.winHeight)*.75)
	offsetX := (int32(ui.winWidth)-invWidth)/2
	offsetY := (int32(ui.winHeight)-invHeight)/2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}

func (ui *ui) getInventoryItemRect(i int) *sdl.Rect {
	invRect := ui.getInventoryRect()
	itemSize := itemSizeRatio * float32(ui.winWidth)
	return  &sdl.Rect{10+invRect.X+int32(i)*int32(itemSize), invRect.Y+invRect.H-int32(itemSize),int32(itemSize),int32(itemSize)}
}

func(ui *ui) DrawInventory(level *game.Level) {
	invRect := ui.getInventoryRect()
	playerRect := ui.textureIndex[level.Player.Rune][0]
	helperX := (int32(float32(invRect.W)*1.01)-invRect.W)/2
	helperY := (int32(float32(invRect.H)*1.01)-invRect.H)/2
	// helperCharX := (int32((float32(invRect.W)/1.25)*1.01)-invRect.W)/2
	// helperCharY := (int32((float32(invRect.H)/2.25)*1.01)-invRect.H)/2
	ui.renderer.Copy(ui.inventoryBorder, nil, &sdl.Rect{invRect.X-helperX,invRect.Y-helperY,int32(float32(invRect.W)*1.01),int32(float32(invRect.H)*1.01)})
	ui.renderer.Copy(ui.inventoryBackground, nil, &sdl.Rect{invRect.X,invRect.Y,invRect.W,invRect.H})
	// ui.renderer.Copy(ui.characterBorder, nil, &sdl.Rect{int32(float32(invRect.X)*1.25)-helperCharX, int32(float32(invRect.Y)*1.15)-helperCharY, int32((float32(invRect.W)/1.25)*1.01), int32((float32(invRect.H)/2.25)*1.01)})
	// ui.renderer.Copy(ui.characterSlotBackground, nil, &sdl.Rect{int32(float32(invRect.X)*1.25), int32(float32(invRect.Y)*1.15), int32(float32(invRect.W)/1.25), int32(float32(invRect.H)/2.25)})
	ui.renderer.Copy(ui.imageAtlas, &playerRect, &sdl.Rect{int32(float32(invRect.X)*1.65), int32(float32(invRect.Y)*1.25), int32(float32(invRect.W)/1.75), int32(float32(invRect.H)/1.75)})
	ui.renderer.Copy(ui.helmetSlotBackground, nil, ui.getHelmetSlotRect())
	if level.Player.Helmet != nil {
		ui.renderer.Copy(ui.imageAtlas, &ui.textureIndex[level.Player.Helmet.Rune][0], ui.getHelmetSlotRect())
	}
	ui.renderer.Copy(ui.swordSlotBackground, nil, ui.getSwordSlotRect())
	if level.Player.Sword != nil {
		ui.renderer.Copy(ui.imageAtlas, &ui.textureIndex[level.Player.Sword.Rune][0], ui.getSwordSlotRect())
	}
	ui.renderer.Copy(ui.armourSlotBackground, nil, ui.getArmorSlotRect())
	if level.Player.Armour != nil {
		ui.renderer.Copy(ui.imageAtlas, &ui.textureIndex[level.Player.Armour.Rune][0], ui.getArmorSlotRect())
	}

	for i, item := range level.Player.Items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		if item == ui.draggedItem {
			itemSize := itemSizeRatio * float32(ui.winWidth)
			ui.renderer.Copy(ui.imageAtlas, &itemSrcRect,&sdl.Rect{int32(ui.currMouseState.pos.X)-(int32(itemSize/2)), int32(ui.currMouseState.pos.Y)-(int32(itemSize/2)), int32(itemSize), int32(itemSize)})
		} else {
			ui.renderer.Copy(ui.imageAtlas, &itemSrcRect,ui.getInventoryItemRect(i))
		}
	}
}

func (ui *ui) CheckEquippedItem() *game.Items {
	if ui.draggedItem.Type == game.Sword {
		r := ui.getSwordSlotRect()
		if r.HasIntersection(&sdl.Rect{int32(ui.currMouseState.pos.X), int32(ui.currMouseState.pos.Y), 1, 1}) {
			return ui.draggedItem
		}
	} else if ui.draggedItem.Type == game.Armour {
		r := ui.getArmorSlotRect()
		if r.HasIntersection(&sdl.Rect{int32(ui.currMouseState.pos.X), int32(ui.currMouseState.pos.Y), 1, 1}) {
			return ui.draggedItem
		}
	} else if ui.draggedItem.Type == game.Helmet {
		r := ui.getHelmetSlotRect()
		if r.HasIntersection(&sdl.Rect{int32(ui.currMouseState.pos.X), int32(ui.currMouseState.pos.Y), 1, 1}) {
			return ui.draggedItem
		}
	}
	return nil
}

func (ui *ui) CheckDroppedItem() *game.Items{
	invRect := ui.getInventoryRect()
	mousePos := ui.currMouseState.pos
	if invRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1,1}) {
		return nil
	} else {
		return ui.draggedItem
	}
}

func (ui *ui) getHelmetSlotRect() *sdl.Rect {
	itemSize := itemSizeRatio * float32(ui.winWidth) * 1.05
	return &sdl.Rect{448,164, int32(itemSize), int32(itemSize)}
}

func (ui *ui) getSwordSlotRect() *sdl.Rect {
	itemSize := itemSizeRatio * float32(ui.winWidth) * 1.05
	return &sdl.Rect{356,248, int32(itemSize), int32(itemSize)}
}

func (ui *ui) getArmorSlotRect() *sdl.Rect {
	itemSize := itemSizeRatio * float32(ui.winWidth) * 1.05
	return &sdl.Rect{332,150, int32(itemSize), int32(itemSize)}
}

func (ui *ui) getBackgroundRect(index int) *sdl.Rect {
	itemSize := itemSizeRatio * float32(ui.winWidth)
	return &sdl.Rect{int32(ui.winWidth-int(itemSize)-(int(itemSize)*index)),int32(ui.winHeight-int(itemSize)),int32(itemSize),int32(itemSize)}
} 

func (ui *ui) CheckInventoryItems(level *game.Level) *game.Items {
	if ui.currMouseState.leftButton {
		mousePos := ui.currMouseState.pos
		for i, item := range level.Player.Items {
			itemRect := ui.getInventoryItemRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				fmt.Println(item.Name + "in inventory has been clicked!")
				return item
			}
		}
	}
	return nil
}

func (ui *ui) DrinkPotion(level *game.Level) {
	if ui.currMouseState.rightButton {
		mousePos := ui.currMouseState.pos
		for i, item := range level.Player.Items {
			if item.Type == game.Potion {
				itemRect := ui.getInventoryItemRect(i)
				if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y),1,1}) {
					level.Player.Items = append(level.Player.Items[:i], level.Player.Items[i+1:]...)
					level.Player.Character.Hp += int(item.Power)
					playRandomSounds(ui.burpSound, 100)
				}
			}
		}
	}
}

func (ui *ui) CheckBackgroundItems(level *game.Level) *game.Items {
	if !ui.currMouseState.leftButton && ui.prevMouseState.leftButton {
		fmt.Println("Clicked, x, y : ", ui.currMouseState.pos)
		items := level.Items[level.Player.Pos]
		mousePos := ui.currMouseState.pos
		for i, item := range items {
			itemRect := ui.getBackgroundRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				return item
			}
		}
	}
	return nil
}