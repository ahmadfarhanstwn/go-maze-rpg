package ui2d

import (
	"math/rand"
	"strconv"

	"github.com/ahmadfarhanstwn/rpg/game-logic"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type fontSize int

type UiState int

const (
	MainUI UiState = iota
	InventoryUI
)

const (
	smallSize fontSize = iota
	mediumSize
	largeSize
)

const itemSizeRatio = .066

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
}


type ui struct {
	state UiState
	sounds
	winWidth int
	winHeight int
	renderer *sdl.Renderer
	imageAtlas *sdl.Texture
	textureIndex map[rune][]sdl.Rect
	keyboardState []uint8
	prevKeyboardState []uint8
	centerX int
	centerY int
	cameraLimit int
	window *sdl.Window
	levelChannel chan *game.Level
	inputChannel chan *game.Input
	r *rand.Rand
	smallFont *ttf.Font
	mediumFont *ttf.Font
	largeFont *ttf.Font
	texSmall map[string]*sdl.Texture
	texMedium map[string]*sdl.Texture
	texLarge map[string]*sdl.Texture
	eventBackground *sdl.Texture
	inventoryBackground *sdl.Texture
	helmetSlotBackground *sdl.Texture
	swordSlotBackground *sdl.Texture
	armourSlotBackground *sdl.Texture
	draggedItem *game.Items
	currMouseState *mouseState
	prevMouseState *mouseState
}

func NewUi(levelChannel chan *game.Level, inputChannel chan *game.Input) *ui {
	ui := &ui{}
	ui.state = MainUI
	ui.r = rand.New(rand.NewSource(1))
	ui.levelChannel = levelChannel
	ui.inputChannel = inputChannel
	ui.winWidth = 800
	ui.winHeight = 600
	ui.cameraLimit = 5
	window, err := sdl.CreateWindow("RPG BABY!!!", 100, 100,
		int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window
	ui.texSmall = make(map[string]*sdl.Texture)
	ui.texMedium = make(map[string]*sdl.Texture)
	ui.texLarge = make(map[string]*sdl.Texture)

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.imageAtlas = ui.imgFileToTexture("game-ui-2d/assets/tiles.png")

	ui.loadTextureIdx()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

	ui.centerX = -1
	ui.centerY = -1

	ui.smallFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 10)
	if err != nil {
		panic(err)
	}

	ui.mediumFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 18)
	if err != nil {
		panic(err)
	}

	ui.largeFont, err = ttf.OpenFont("game-ui-2d/assets/Kingthings.ttf", 26)
	if err != nil {
		panic(err)
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{0,0,0,128})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.inventoryBackground = ui.GetSinglePixelTex(sdl.Color{85,52,165,128})
	ui.inventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.helmetSlotBackground = ui.GetSinglePixelTex(sdl.Color{0,0,0,128})
	ui.helmetSlotBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.swordSlotBackground = ui.GetSinglePixelTex(sdl.Color{0,0,0,128})
	ui.swordSlotBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.armourSlotBackground = ui.GetSinglePixelTex(sdl.Color{0,0,0,128})
	ui.armourSlotBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.initSound()

	return ui
}

func (ui *ui) Draw(level *game.Level) {

	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	if level.Player.X > ui.centerX+ui.cameraLimit {
		diff := level.Player.X - (ui.centerX+ui.cameraLimit)
		ui.centerX += diff
	} else if level.Player.X < ui.centerX-ui.cameraLimit {
		diff := (ui.centerX-ui.cameraLimit) - level.Player.X
		ui.centerX -= diff
	} else if level.Player.Y > ui.centerY+ui.cameraLimit {
		diff := level.Player.Y - (ui.centerY+ui.cameraLimit)
		ui.centerY += diff
	} else if level.Player.Y < ui.centerY-ui.cameraLimit {
		diff := (ui.centerY-ui.cameraLimit) - level.Player.Y
		ui.centerY -= diff
	}

	offsetX := (ui.winWidth/2) - ui.centerX*32
	offsetY := (ui.winHeight/2) - ui.centerY*32

	ui.renderer.Clear()
	ui.r.Seed(1)
	for y, rows := range level.Map {
		for x, cols := range rows {
			if cols.Rune != game.Blank {
				r := ui.r.Intn(len(ui.textureIndex[cols.Rune]))
				srcRect := ui.textureIndex[cols.Rune][r]
				if level.Map[y][x].Visible || level.Map[y][x].Seen {
					destRect := sdl.Rect{int32(x*32)+int32(offsetX),int32(y*32)+int32(offsetY),32,32}

					pos := game.Pos{x,y}
					if level.Debug[pos] {
						ui.imageAtlas.SetColorMod(128,0,0)
					} else if level.Map[y][x].Seen && !level.Map[y][x].Visible {
						ui.imageAtlas.SetColorMod(128,128,128)
					} else {
						ui.imageAtlas.SetColorMod(255,255,255)
					}
					ui.renderer.Copy(ui.imageAtlas, &srcRect, &destRect)

					if level.Map[y][x].OverlayRune != game.Blank {
						srcRect = ui.textureIndex[cols.OverlayRune][0]
						ui.renderer.Copy(ui.imageAtlas, &srcRect, &destRect)
					}
				}
			}
		}
	}

	ui.imageAtlas.SetColorMod(255,255,255)

	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {
			monsterRect := ui.textureIndex[monster.Rune][0]
			ui.renderer.Copy(ui.imageAtlas, &monsterRect, &sdl.Rect{int32(pos.X*32)+int32(offsetX),int32(pos.Y*32)+int32(offsetY),32,32})
		}
	}

	ui.r.Seed(1)
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				r := ui.r.Intn(len(ui.textureIndex[item.Rune]))
				itemRect := ui.textureIndex[item.Rune][r]
				ui.renderer.Copy(ui.imageAtlas, &itemRect, &sdl.Rect{int32(pos.X*32)+int32(offsetX),int32(pos.Y*32)+int32(offsetY),32,32})
			}
		}
	}

	playerSrc := ui.textureIndex['@'][0]
	ui.renderer.Copy(ui.imageAtlas, &playerSrc, &sdl.Rect{int32(level.Player.X*32)+int32(offsetX),int32(level.Player.Y*32)+int32(offsetY),32,32})
	
	startEventH := int32(float64(ui.winHeight)*0.75)
	backgroundWidth := int32(float64(ui.winWidth)*0.30)
	
	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{0, startEventH, backgroundWidth, int32(ui.winHeight)-startEventH})

	for i, event := range level.Events {
		if event != "" {
			tex := ui.stringToFont(event, mediumSize, sdl.Color{255,0,0,0})
			_,_,w,h,err := tex.Query()
			if err != nil {
				panic(err)
			}
			ui.renderer.Copy(tex, nil, &sdl.Rect{0,int32(i*12)+startEventH,w,h})
		}
	}

	hp := ui.stringToFont("Player HP : "+strconv.FormatInt(int64(level.Player.Hp),10),mediumSize, sdl.Color{255,255,255,0})
	_,_,w,h,err := hp.Query()
	if err != nil {
		panic(err)
	}
	ui.renderer.Copy(hp, nil, &sdl.Rect{0,0,w,h})

	// inventoryStart := int32(float64(ui.winWidth)*.4)
	// inventoryWidth := int32(ui.winWidth)-inventoryStart
	// // itemSize := int32(itemSizeRatio * float32(ui.winWidth))
	// ui.renderer.Copy(ui.inventoryBackground, nil, &sdl.Rect{inventoryStart, int32(ui.winHeight), inventoryWidth, int32(32)})

	items := level.Player.Character.Items
	for i, item := range items {
		itemRect := ui.textureIndex[item.Rune][0]
		ui.renderer.Copy(ui.imageAtlas, &itemRect, ui.getBackgroundRect(i))
	}
	// sdl.Delay(16)
}

func (ui *ui) Run() {
	ui.prevMouseState = getMouseState()
	var newLevel *game.Level
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChannel <- &game.Input{Input: game.Quit}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChannel <- &game.Input{Input: game.CloseWindow,LevelChannel: ui.levelChannel}
				}
			}
		}

		ui.currMouseState = getMouseState()
		var input game.Input
		var ok bool

		select {
		case newLevel, ok = <-ui.levelChannel:
			if ok {
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSounds(ui.sounds.footStep,20)
				case game.OpenDoor:
					playRandomSounds(ui.openDoor, 75)
				case game.Attacking:
					playRandomSounds(ui.attackingSound, 75)
				case game.PickUpItems, game.DropItems:
					playRandomSounds(ui.pickUpItems, 75)
				case game.Portal:
					playRandomSounds(ui.enteringPortals, 100)
				}
			}
		default:
		}

		ui.Draw(newLevel)
		if ui.state == InventoryUI {
			if ui.draggedItem != nil && !ui.currMouseState.leftButton && ui.prevMouseState.leftButton {
				item := ui.CheckDroppedItem()
				if item != nil {
					input.Input = game.DropItem
					input.Item = item
					ui.draggedItem = nil 
				}

				item = ui.CheckEquippedItem()
				if item != nil {
					input.Input = game.EquipItem
					input.Item = item
					ui.draggedItem = nil
				}
			}
			if !ui.currMouseState.leftButton || ui.draggedItem == nil {
				ui.draggedItem = ui.CheckInventoryItems(newLevel)	
			}
			ui.DrawInventory(newLevel)
		}
		ui.renderer.Present()

		item := ui.CheckBackgroundItems(newLevel)
		if item != nil {
			input.Input = game.TakeItem
			input.Item = item
		}

		ui.checkInput(input)
		ui.prevMouseState = ui.currMouseState
		sdl.Delay(10)
	}
}