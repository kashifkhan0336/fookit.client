// package main

// import (
// 	"runtime"

// 	"github.com/go-gl/glfw/v3.3/glfw"
// )

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

// func main() {
// 	err := glfw.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer glfw.Terminate()

// 	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	window.MakeContextCurrent()
// 	window.SetOpacity(0)
// 	window.Maximize()
// 	for !window.ShouldClose() {
// 		// Do OpenGL stuff.
// 		window.SwapBuffers()
// 		glfw.PollEvents()

// 	}
// }

package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/daspoet/gowinkey"
	"github.com/gen2brain/beeep"
	"github.com/kbinani/screenshot"
)

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
func main() {
	keys := []gowinkey.VirtualKey{
		gowinkey.VK_RBUTTON,
		gowinkey.VK_Q,
	}
	err := beeep.Notify("System", "Started", "assets/information.png")
	if err != nil {
		panic(err)
	}
	events, stopFn := gowinkey.Listen(gowinkey.Selective(keys...))

	time.AfterFunc(time.Minute, func() {
		stopFn()
	})

	for e := range events {
		switch e.State {
		case gowinkey.KeyDown:

			x1, y1 = getPos()
		case gowinkey.KeyUp:
			x2, y2 = getPos()

			go snap(int(x1), int(y1), int(x2), int(y2))
			fmt.Println("released", e)
		}
	}
}

// package main

// import (
// 	"time"

// 	"github.com/daspoet/gowinkey"
// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// )

// const (
// 	screenWidth  = 300
// 	screenHeight = 300
// )

// var (
// 	isDisplayed = false
// )

// type Game struct{}

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }
// func (g *Game) Init() {

// }
// func (g *Game) Update() error {
// 	return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	if isDisplayed {
// 		ebitenutil.DebugPrint(screen, "This is a test.")
// 	}

// }

// const (
// 	width  = 200
// 	height = 200
// )

// func main() {
// 	go func() {
// 		keys := []gowinkey.VirtualKey{
// 			gowinkey.VK_Q,
// 		}
// 		events, stopFn := gowinkey.Listen(gowinkey.Selective(keys...))

// 		time.AfterFunc(time.Minute, func() {
// 			stopFn()
// 		})

// 		for e := range events {
// 			switch e.State {
// 			case gowinkey.KeyDown:

// 				// x1, x2 = getPos()qqqqq

// 			case gowinkey.KeyUp:
// 				// y1, y2 = getPos()
// 				println("Pressed!")
// 				// width := math.Abs(float64(x2) - float64(x1))
// 				// height := math.Abs(float64(y2) - float64(y1))
// 				// fmt.Println(width)
// 				// fmt.Println(height)
// 				// go snap(width, height)qqq
// 				// fmt.Println("released", e)
// 			}
// 		}
// 	}()
// 	ebiten.SetWindowTitle("Ebiten Test")
// 	ebiten.SetWindowDecorated(false)
// 	ebiten.SetWindowFloating(true)
// 	ebiten.SetWindowSize(width, height)
// 	ebiten.SetWindowMousePassthrough(true)

//		op := &ebiten.RunGameOptions{}
//		op.ScreenTransparent = true
//		op.SkipTaskbar = false
//		if err := ebiten.RunGameWithOptions(&Game{}, op); err != nil {
//			panic(err)
//		}
//	}
