package modelhooks

import (
	"github.com/pocketbase/pocketbase"
)

func videoQuoteAudio_RegisterHooks(app *pocketbase.PocketBase) {
	videoQuoteAudio_UpdateAudioFileMd5(app)
}

func videoQuoteAudio_UpdateAudioFileMd5(app *pocketbase.PocketBase) {
	// app.OnModelBeforeCreate().Add(func(e *core.ModelEvent) error {
	// 	rec := e.Model.(*models.Record)
	// 	if e.Model.TableName() == customModels.VideoQuoteAudiosTableName && rec.Get("title") == "" {
	// 		rec.Set("title", "Draft")
	// 	}
	// 	return nil
	// })
}
