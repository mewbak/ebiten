package main

import (
	"github.com/hajimehoshi/go-ebiten/example/blocks"
	"github.com/hajimehoshi/go-ebiten/graphics"
	"github.com/hajimehoshi/go-ebiten/ui"
	"github.com/hajimehoshi/go-ebiten/ui/cocoa"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Game interface {
	Update(state ui.CanvasState)
	Draw(c graphics.Context)
}

func init() {
	runtime.LockOSThread()
}

func main() {
	const screenWidth = blocks.ScreenWidth
	const screenHeight = blocks.ScreenHeight
	const screenScale = 2
	const fps = 60
	const frameTime = time.Duration(int64(time.Second) / int64(fps))
	const title = "Ebiten Demo"

	u := cocoa.UI()
	canvas := u.CreateCanvas(screenWidth, screenHeight, screenScale, title)

	textureFactory := cocoa.TextureFactory()
	var game Game = blocks.NewGame(NewTextures(textureFactory))
	tick := time.Tick(frameTime)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	u.Start()
	defer u.Terminate()
	for {
		u.DoEvents()
		select {
		default:
			canvas.Draw(game.Draw)
		case <-tick:
			state := canvas.State()
			game.Update(state)
			if state.IsClosed {
				return
			}
		/*case e := <-windowEvents:
			game.HandleEvent(e)
			if _, ok := e.(ui.WindowClosedEvent); ok {
				return
			}*/
		case <-sigterm:
			return
		}
	}
}
