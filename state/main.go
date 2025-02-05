package main

import (
	"errors"
	"image/color"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil" // required for debug text
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil" // required for isKeyJustPressed
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type gameState int

const (
	titleScreen gameState = iota
	options
	play
	quit
)

// define some kind of palette?
var (
	green1  = &color.NRGBA{0x00, 0x38, 0x40, 0xff}
	green2  = &color.NRGBA{0x00, 0x5a, 0x5b, 0xff}
	green3  = &color.NRGBA{0x00, 0x73, 0x69, 0xff}
	green4  = &color.NRGBA{0x00, 0x8c, 0x72, 0xff}
	green5  = &color.NRGBA{0x02, 0xa6, 0x76, 0xff}
	purple1 = &color.NRGBA{0x30, 0x28, 0x40, 0xff}
	purple2 = &color.NRGBA{0x47, 0x39, 0x5b, 0xff}
	purple3 = &color.NRGBA{0x5f, 0x49, 0x73, 0xff}
	purple4 = &color.NRGBA{0x7b, 0x58, 0x8c, 0xff}
	purple5 = &color.NRGBA{0x99, 0x69, 0xa6, 0xff}
)

// MenuItem represents an item in a menu list
type MenuItem struct {
	Name         string
	Text         string
	TxtX         int           // optional X location to draw text, if not provided x is 0
	TxtY         int           // optional Y location to draw text, it not provided y is the menu list height - 5
	image        *ebiten.Image // used to store the autogenerated image for the menu item
	BgColour     *color.NRGBA  // optional background colour, overrides default colour
	TxtColour    *color.NRGBA  // optional text colour, overrides default text colour
	SelBgColour  *color.NRGBA  // optional selected background colour, overrides default selected colour
	SelTxtColour *color.NRGBA  // optional selected text colour, overrides default selected text colour
}

// TO DO - Text features
// Support default text alignment: left, right, centre

// MenuList is a navigatable, selectable menu
type MenuList struct {
	Tx                  float64      // x translation of the menu
	Ty                  float64      // y translation of the menu
	Width               int          // width of all menu items
	Height              int          // height of all menu items
	Offx                float64      // x offset of subsequent menu items
	Offy                float64      // y offset of subsequent menu items
	DefaultBgColour     *color.NRGBA // default background colour
	DefaultTxtColour    *color.NRGBA // default text colour
	DefaultSelBgColour  *color.NRGBA // default selected background colour
	DefaultSelTxtColour *color.NRGBA // default selected text colour
	SelectedIndex       *int         // index of the item in list which is selected
	MenuItems           []MenuItem   // menu items
}

// MenuListInput is an object used to create a menu list
type MenuListInput struct {
	Tx                  float64      // optional, x translation of the menu, if not provided will be 0
	Ty                  float64      // optional, y translation of the menu, if not provided will be 0
	Width               int          // mandatory, width of all menu items
	Height              int          // mandatory, height of all menu items
	Offx                float64      // optional, offset of subsequent menu items, if not provided will 0
	Offy                float64      // optional, offset of subsequent menu items, if not provided will be menu item height
	DefaultBgColour     *color.NRGBA // optional, default background colour of menu, if not provided will be cyan
	DefaultTxtColour    *color.NRGBA // optional, default text colour, if not provided will be black
	DefaultSelBGColour  *color.NRGBA // optional, default selected background colour of menu, if not provided will be magenta
	DefaultSelTxtColour *color.NRGBA //optional, default selected text colour of menu, if not provided it will be white
	MenuItems           []MenuItem   // mandtory, list of menu items
}

//NewMenu constructs a new menu from a MenuListInput
func NewMenu(input MenuListInput) (MenuList, error) {

	if input.Width == 0 {
		return MenuList{}, errors.New("Mandatory input field width is missing")
	}
	if input.Height == 0 {
		return MenuList{}, errors.New("Mandatory input field height is missing")
	}
	if len(input.MenuItems) < 1 {
		return MenuList{}, errors.New("Mandatory input field MenuItems is missing")
	}

	if input.Offy == 0 {
		input.Offy = float64(input.Height)
	}

	if input.DefaultBgColour == nil {
		input.DefaultBgColour = &color.NRGBA{0x00, 0xff, 0xff, 0xff}
	}

	if input.DefaultTxtColour == nil {
		input.DefaultTxtColour = &color.NRGBA{0x00, 0x00, 0x00, 0xff}
	}

	if input.DefaultSelBGColour == nil {
		input.DefaultSelBGColour = &color.NRGBA{0xff, 0x00, 0xff, 0xff}
	}

	if input.DefaultSelTxtColour == nil {
		input.DefaultSelTxtColour = &color.NRGBA{0xff, 0xff, 0xff, 0xff}
	}

	defaultSelectedIndex := 0

	ml := MenuList{
		Tx:                  input.Tx,
		Ty:                  input.Ty,
		Width:               input.Width,
		Height:              input.Height,
		Offx:                input.Offx,
		Offy:                input.Offy,
		DefaultBgColour:     input.DefaultBgColour,
		DefaultTxtColour:    input.DefaultTxtColour,
		DefaultSelBgColour:  input.DefaultSelBGColour,
		DefaultSelTxtColour: input.DefaultSelTxtColour,
		SelectedIndex:       &defaultSelectedIndex,
		MenuItems:           input.MenuItems,
	}

	// set override colours if needed otherwise use default colours
	for i, item := range input.MenuItems {
		if item.BgColour != nil {
			ml.MenuItems[i].BgColour = item.BgColour
		} else {
			ml.MenuItems[i].BgColour = ml.DefaultBgColour
		}

		if item.TxtColour != nil {
			ml.MenuItems[i].TxtColour = item.TxtColour
		} else {
			ml.MenuItems[i].TxtColour = ml.DefaultTxtColour
		}

		if item.SelBgColour != nil {
			ml.MenuItems[i].SelBgColour = item.SelBgColour
		} else {
			ml.MenuItems[i].SelBgColour = ml.DefaultSelBgColour
		}

		if item.SelTxtColour != nil {
			ml.MenuItems[i].SelTxtColour = item.SelTxtColour
		} else {
			ml.MenuItems[i].SelTxtColour = ml.DefaultSelTxtColour
		}

		if item.TxtY == 0 {
			ml.MenuItems[i].TxtY = ml.Height - 5  // default value for text y height
		}

	}

	// initialise images for each menu item
	for i := range ml.MenuItems {
		newImage, _ := ebiten.NewImage(ml.Width, ml.Height, ebiten.FilterNearest)
		ml.MenuItems[i].image = newImage
	}
	return ml, nil
}

//GetSelectedItem returns then name of the selected item
func (m *MenuList) GetSelectedItem() string {
	return m.MenuItems[*m.SelectedIndex].Name
}

//IncrementSelected increments the selected index provided it is not already at maximum
func (m *MenuList) IncrementSelected() {
	maxIndex := len(m.MenuItems) - 1
	if *m.SelectedIndex < maxIndex {
		*m.SelectedIndex++
	}
}

//DecrementSelected decrements the selected index provided it is not already at minimum
func (m *MenuList) DecrementSelected() {
	minIndex := 0
	if *m.SelectedIndex > minIndex {
		*m.SelectedIndex--
	}
}

//Draw draws the menu to the screen
func (m *MenuList) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(m.Tx, m.Ty)

	for index, item := range m.MenuItems {

		if index == *m.SelectedIndex {
			item.image.Fill(item.SelBgColour)
		} else {
			item.image.Fill(item.BgColour)
		}

		if index == *m.SelectedIndex {
			text.Draw(item.image, item.Text, mplusNormalFont, item.TxtX, item.TxtY, item.SelTxtColour)
		} else {
			text.Draw(item.image, item.Text, mplusNormalFont, item.TxtX, item.TxtY, item.TxtColour)
		}

		screen.DrawImage(item.image, opts)
		opts.GeoM.Translate(m.Offx, m.Offy)
	}
}

var (
	state           gameState
	playImage       *ebiten.Image
	optionsImage    *ebiten.Image
	quitImage       *ebiten.Image
	square          *ebiten.Image
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	mainMenu        MenuList
)

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func update(screen *ebiten.Image) error {

	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})

	if state == titleScreen {

		ebitenutil.DebugPrint(screen, "Title screen")
		mainMenu.Draw(screen)

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			mainMenu.DecrementSelected()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			mainMenu.IncrementSelected()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch mainMenu.GetSelectedItem() {
			case "playButton":
				state = play
			case "optionButton":
				state = options
			case "quitButton":
				os.Exit(0)
			}
			return nil
		}

	}

	if state == play {
		ebitenutil.DebugPrint(screen, "Play screen")

		if square == nil {
			square, _ = ebiten.NewImage(32, 32, ebiten.FilterNearest)
		}
		someColor := &color.NRGBA{0x7f, 0xff, 0x00, 0xff}
		square.Fill(someColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(64.0, 64.0)
		screen.DrawImage(square, opts)

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			state = titleScreen
			return nil
		}

	}

	if state == options {
		ebitenutil.DebugPrint(screen, "Options screen")

		if square == nil {
			square, _ = ebiten.NewImage(32, 32, ebiten.FilterNearest)
		}
		someColor := &color.NRGBA{0x8a, 0x2b, 0xe2, 0xff}
		square.Fill(someColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(64.0, 64.0)
		screen.DrawImage(square, opts)

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			state = titleScreen
			return nil
		}
	}

	return nil
}

func main() {

	newMenuItems := []MenuItem{
		{Name: "playButton",
			Text:     "PLAY",
			TxtX: 36,
			TxtY: 25,
			BgColour: green1},
		{Name: "optionButton",
			Text:     "OPTIONS",
			TxtX: 12,
			TxtY: 25,
			BgColour: green2},
		{Name: "quitButton",
			Text:     "QUIT",
			TxtX: 36,
			TxtY: 25,
			BgColour: green3},
	}

	newMenuInput := MenuListInput{
		Width:              128,
		Height:             36,
		Tx:                 128,
		Ty:                 128,
		DefaultSelBGColour: purple3,
		MenuItems:          newMenuItems,
	}

	newMenu, err := NewMenu(newMenuInput)

	if err != nil {
		log.Printf("unable to create menu: %+v\n", err)
	}

	mainMenu = newMenu

	state = titleScreen

	if err := ebiten.Run(update, 400, 300, 2, "State!"); err != nil {
		panic(err)
	}
}
