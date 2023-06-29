package shotcutProjectBuilder

import "time"

var SHOTCUT_VERSION = "23.05.14"
var MLT_VERSION = "7.17.0"
var MLT_LC_NUMERIC = "C"

var INTRO_DURATION = 3500 * time.Millisecond                      // in ms
var OUTRO_DURATION = 7500 * time.Millisecond                      // in ms
var OUTRO_OVERLAY_DURATION = 6000 * time.Millisecond              // in ms
var OUTRO_OVERLAY_DELAY = OUTRO_DURATION - OUTRO_OVERLAY_DURATION // in ms
var BG_MUSIC_TRUNCATE_DURATION = 10000 * time.Millisecond         // truncate last 10 seconds of background music
var BG_MUSIC_FADE_DURATION = 2000 * time.Millisecond              // fade in and fade out duration of background music
