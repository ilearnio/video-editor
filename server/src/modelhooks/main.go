package modelhooks

import "github.com/pocketbase/pocketbase"

func Register(app *pocketbase.PocketBase) {
	videoQuoteAudio_RegisterHooks(app)
}
