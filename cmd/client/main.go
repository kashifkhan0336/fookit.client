package main

// // import (
// // 	"runtime"

// // 	"github.com/go-gl/glfw/v3.3/glfw"
// // )

// // func init() {
// // 	// This is needed to arrange that main() runs on main thread.
// // 	// See documentation for functions that are only allowed to be called from the main thread.
// // 	runtime.LockOSThread()
// // }

// // func main() {
// // 	err := glfw.Init()
// // 	if err != nil {
// // 		panic(err)
// // 	}
// // 	defer glfw.Terminate()

// // 	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	window.MakeContextCurrent()
// // 	window.SetOpacity(0)
// // 	window.Maximize()
// // 	for !window.ShouldClose() {
// // 		// Do OpenGL stuff.
// // 		window.SwapBuffers()
// // 		glfw.PollEvents()

// // 	}
// // }

// package main

import (
	"context"
	"fmt"

	"image"
	"image/png"
	"math/rand"
	"os"
	"syscall"
	"unsafe"

	"fyne.io/systray/example/icon"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/daspoet/gowinkey"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kbinani/screenshot"
)

func credentials() (*cloudinary.Cloudinary, context.Context) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, _ := cloudinary.New()
	cld.Config.URL.Secure = true
	cld.Config.Cloud.APIKey = "483785888746772"
	cld.Config.Cloud.APISecret = "dpIhboL06S4usggEA7_lIi5O9oU"
	cld.Config.Cloud.CloudName = "dsryinzqb"
	ctx := context.Background()
	return cld, ctx
}

func getPos() (int32, int32) {
	userDll := syscall.NewLazyDLL("user32.dll")
	getWindowRectProc := userDll.NewProc("GetCursorPos")
	type POINT struct {
		X, Y int32
	}
	var pt POINT
	_, _, eno := syscall.SyscallN(getWindowRectProc.Addr(), uintptr(unsafe.Pointer(&pt)))
	if eno != 0 {
		fmt.Println(eno)
	}
	return pt.X, pt.Y
}

var (
	x1 int32
	x2 int32
	y1 int32
	y2 int32
)

func snap(x1, x2, y1, y2 int) {

	n := screenshot.NumActiveDisplays()
	println(x1, y1, x2, y1)
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(image.Rect(x1, x2, y1, y2))
		if err != nil {
			fmt.Printf("error occured : %s", err)
			return
		}
		fileName := fmt.Sprintf("%d_%dx%d_%d.png", i, bounds.Dx(), bounds.Dy(), rand.Intn(50000-1))
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)
		err = beeep.Notify("System", "Screenshot taken", "assets/information.png")
		if err != nil {
			panic(err)
		}
		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
}

const (
	screenWidth  = 300
	screenHeight = 300
)

var (
	isDisplayed = true
	text        = "Yeeeeeeeeeeeeeeesh!"
)

type Game struct{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func (g *Game) Init() {

}
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		isDisplayed = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		isDisplayed = true
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if isDisplayed {
		ebitenutil.DebugPrint(screen, text)
	}

}

const (
	width  = 200
	height = 200
)

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuitOrig.SetIcon(icon.Data)
}

func onExit() {
	print("onExit called!")
	os.Exit(0)
	// clean up here
}
func changeText(textStr string) {
	isDisplayed = false
	text = textStr
	isDisplayed = true
}
func main() {

	go func() {
		keys := []gowinkey.VirtualKey{
			gowinkey.VK_RBUTTON,
			gowinkey.VK_Q,
		}
		events, _ := gowinkey.Listen(gowinkey.Selective(keys...))

		for e := range events {
			pressedKey := e.VirtualKey
			switch e.State {
			case gowinkey.KeyDown:
				if pressedKey == gowinkey.VK_RBUTTON {
					changeText("Started!")
					// Right button was pressed
					print("Right click")
					x1, y1 = getPos()
				} else if pressedKey == gowinkey.VK_Q {
					// Q key was pressed
					print("Q Presssed")
					systray.Quit()
				}

				print("Started!")
			case gowinkey.KeyUp:
				if pressedKey == gowinkey.VK_RBUTTON {
					changeText("Stopped!")
					// Right button was pressed
					print("Right click")
					x2, y2 = getPos()
				} else if pressedKey == gowinkey.VK_Q {
					// Q key was pressed
					print("Q Presssed")
					systray.Quit()
				}
				print("Ended!")

				go snap(int(x1), int(y1), int(x2), int(y2))
			}
		}
	}()
	go systray.Run(onReady, onExit)
	ebiten.SetWindowTitle("Ebiten Test")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowMousePassthrough(true)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true

	if err := ebiten.RunGameWithOptions(&Game{}, op); err != nil {
		panic(err)
	}
}
