package menulist

import(
	"log"
	"errors"
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
)

var mplusNormalFont font.Face

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
