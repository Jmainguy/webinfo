package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

func writeText(id, text string) {
	doc := js.Global().Get("document")
	elem := doc.Call("getElementById", id)
	if elem.Truthy() {
		elem.Set("innerText", text)
	} else {
		js.Global().Get("console").Call("error", "Element not found:", id)
	}
}

func getUserAgent() string {
	return js.Global().Get("navigator").Get("userAgent").String()
}

func getLanguage() string {
	return js.Global().Get("navigator").Get("language").String()
}

func getPlatform() string {
	return js.Global().Get("navigator").Get("platform").String()
}

func getTimezone() string {
	tz := js.Global().Get("Intl").Get("DateTimeFormat").New().Call("resolvedOptions").Get("timeZone")
	js.Global().Get("console").Call("log", "getTimezone", tz)
	return tz.String()
}

func getScreenSize() string {
	width := js.Global().Get("screen").Get("width")
	height := js.Global().Get("screen").Get("height")
	js.Global().Get("console").Call("log", "getScreenSize", width, height)
	return fmt.Sprintf("%dx%d", width.Int(), height.Int())
}

func getGPU() string {
	canvas := js.Global().Get("document").Call("createElement", "canvas")
	gl := canvas.Call("getContext", "webgl")
	if !gl.Truthy() {
		return "No WebGL"
	}
	debugInfo := gl.Call("getExtension", "WEBGL_debug_renderer_info")
	if debugInfo.Truthy() {
		renderer := gl.Call("getParameter", debugInfo.Get("UNMASKED_RENDERER_WEBGL")).String()
		return renderer
	}
	return "Unknown GPU"
}

func getConnectionType() string {
	connection := js.Global().Get("navigator").Get("connection")
	if connection.Truthy() {
		effectiveType := connection.Get("effectiveType")
		if effectiveType.Truthy() {
			return effectiveType.String()
		}
	}
	return "Not Supported"
}

func getDeviceMemory() string {
	mem := js.Global().Get("navigator").Get("deviceMemory")
	if mem.Truthy() {
		return fmt.Sprintf("%.0f GB", mem.Float())
	}
	return "Not Supported"
}

func getHardwareConcurrency() string {
	cores := js.Global().Get("navigator").Get("hardwareConcurrency")
	if cores.Truthy() {
		return fmt.Sprintf("%d", cores.Int())
	}
	return "Unknown"
}

func getCookiesEnabled() string {
	enabled := js.Global().Get("navigator").Get("cookieEnabled")
	if enabled.Truthy() {
		return fmt.Sprintf("%v", enabled.Bool())
	}
	return "Unknown"
}

func getOnlineStatus() string {
	online := js.Global().Get("navigator").Get("onLine")
	if online.Truthy() {
		return fmt.Sprintf("%v", online.Bool())
	}
	return "Unknown"
}

// Parse browser name/version from user agent (very basic)
func getBrowser() string {
	ua := getUserAgent()
	switch {
	case js.Global().Get("navigator").Get("brave").Truthy():
		return "Brave (Chromium-based)"
	case js.Global().Get("navigator").Get("userAgentData").Truthy():
		// Chromium-based browsers may expose userAgentData
		brands := js.Global().Get("navigator").Get("userAgentData").Get("brands")
		if brands.Length() > 0 {
			brand := brands.Index(0).Get("brand").String()
			version := brands.Index(0).Get("version").String()
			return brand + " " + version
		}
	}
	// Fallback: simple substring matching
	switch {
	case strings.Contains(ua, "Edg/"):
		return "Microsoft Edge"
	case strings.Contains(ua, "OPR/"):
		return "Opera"
	case strings.Contains(ua, "Chrome/"):
		return "Chrome"
	case strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/"):
		return "Safari"
	case strings.Contains(ua, "Firefox/"):
		return "Firefox"
	default:
		return "Unknown"
	}
}

func getBrowserVersion() string {
	ua := getUserAgent()
	// Try to extract version for common browsers
	switch {
	case strings.Contains(ua, "Edg/"):
		return extractVersionGo(ua, "Edg/")
	case strings.Contains(ua, "OPR/"):
		return extractVersionGo(ua, "OPR/")
	case strings.Contains(ua, "Chrome/"):
		return extractVersionGo(ua, "Chrome/")
	case strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/"):
		return extractVersionGo(ua, "Version/")
	case strings.Contains(ua, "Firefox/"):
		return extractVersionGo(ua, "Firefox/")
	default:
		return ""
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (indexOf(s, substr) != -1)))
}

func indexOf(s, substr string) int {
	return js.Global().Get("String").New(s).Call("indexOf", substr).Int()
}

func extractVersion(ua, marker string) string {
	idx := indexOf(ua, marker)
	if idx == -1 {
		return ""
	}
	start := idx + len(marker)
	end := start
	for end < len(ua) && ((ua[end] >= '0' && ua[end] <= '9') || ua[end] == '.') {
		end++
	}
	return ua[start:end]
}

func extractVersionGo(ua, marker string) string {
	idx := strings.Index(ua, marker)
	if idx == -1 {
		return ""
	}
	start := idx + len(marker)
	end := start
	for end < len(ua) && ((ua[end] >= '0' && ua[end] <= '9') || ua[end] == '.') {
		end++
	}
	return ua[start:end]
}

func getHistoryLength() string {
	return fmt.Sprintf("%d", js.Global().Get("history").Get("length").Int())
}

func showMap(lat, lon float64) {
	js.Global().Call("showMap", lat, lon)
}

func fetchLocation() {
	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pos := args[0]
		coords := pos.Get("coords")
		lat := coords.Get("latitude").Float()
		lon := coords.Get("longitude").Float()
		showMap(lat, lon)
		return nil
	})
	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("showMap", js.Null(), js.Null())
		return nil
	})

	js.Global().Get("navigator").Get("geolocation").Call("getCurrentPosition", successCallback, errorCallback)
}

func fetchBattery() {
	navigator := js.Global().Get("navigator")
	if !navigator.Truthy() || navigator.Get("getBattery").Type() != js.TypeFunction {
		writeText("battery", "Battery API not supported")
		return
	}

	batteryPromise := navigator.Call("getBattery")
	thenCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		battery := args[0]
		level := battery.Get("level").Float() * 100
		isCharging := battery.Get("charging").Bool()
		status := "Charging"
		if !isCharging {
			status = "Not Charging"
		}
		writeText("battery", fmt.Sprintf("%.0f%% (%s)", level, status))
		return nil
	})
	batteryPromise.Call("then", thenCallback)
}

func startClock() {
	go func() {
		for range time.Tick(time.Second) {
			now := time.Now()
			writeText("clock", now.Format("15:04:05"))
		}
	}()
}

// Clipboard API (polling, async)
var lastClipboard string

func pollClipboard() {
	navigator := js.Global().Get("navigator")
	clipboard := navigator.Get("clipboard")
	// Only poll if Clipboard API is available
	if clipboard.Truthy() && clipboard.Get("readText").Type() == js.TypeFunction {
		go func() {
			for {
				document := js.Global().Get("document")
				if document.Truthy() && document.Get("hasFocus").Type() == js.TypeFunction && document.Call("hasFocus").Bool() {
					promise := clipboard.Call("readText")
					then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						text := args[0].String()
						if text != lastClipboard {
							lastClipboard = text
							writeText("clipboard", text)
						}
						return nil
					})
					catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						writeText("clipboard", "Clipboard access denied or unavailable")
						return nil
					})
					promise.Call("then", then).Call("catch", catch)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()
	} else {
		writeText("clipboard", "Clipboard API not available")
	}
}

// Contacts API (very limited support, requires user gesture)
func fetchContacts() {
	navigator := js.Global().Get("navigator")
	contacts := navigator.Get("contacts")
	if contacts.Truthy() && contacts.Get("select").Type() == js.TypeFunction {
		// Only works with user gesture, so just show API is available
		writeText("contacts", "Contacts API available (requires user gesture)")
	} else {
		writeText("contacts", "Contacts API not available")
	}
}

func getPreferredColorScheme() string {
	mql := js.Global().Get("window").Call("matchMedia", "(prefers-color-scheme: dark)")
	if mql.Truthy() && mql.Get("matches").Bool() {
		return "Dark"
	}
	return "Light"
}

func getTouchSupport() string {
	if js.Global().Get("window").Get("ontouchstart").Type() != js.TypeUndefined {
		return "Yes"
	}
	maxTouch := js.Global().Get("navigator").Get("maxTouchPoints")
	if maxTouch.Truthy() && maxTouch.Int() > 0 {
		return "Yes"
	}
	return "No"
}

// Window Size polling
var lastWindowSize string

func pollWindowSize() {
	go func() {
		for {
			w := js.Global().Get("window").Get("innerWidth")
			h := js.Global().Get("window").Get("innerHeight")
			size := "Unknown"
			if w.Truthy() && h.Truthy() {
				size = fmt.Sprintf("%dx%d", w.Int(), h.Int())
			}
			if size != lastWindowSize {
				lastWindowSize = size
				writeText("windowSize", size)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func getDisplayMode() string {
	mql := js.Global().Get("window").Call("matchMedia", "(display-mode: standalone)")
	if mql.Truthy() && mql.Get("matches").Bool() {
		return "Standalone/PWA"
	}
	return "Browser Tab"
}

func getCookies() string {
	cookies := js.Global().Get("document").Get("cookie").String()
	if cookies == "" {
		return "None"
	}
	return cookies
}

// Keep references to prevent GC
var jsCallbacks []js.Func

func setupCameraAndMic() {
	doc := js.Global().Get("document")
	video := doc.Call("getElementById", "cameraVideo")
	cameraStatus := doc.Call("getElementById", "cameraStatus")
	micVolume := doc.Call("getElementById", "micVolume")
	micStatus := doc.Call("getElementById", "micStatus")
	cameraBtn := doc.Call("getElementById", "cameraStartBtn")

	// Ensure video element has correct attributes for autoplay
	video.Set("autoplay", true)
	video.Set("playsInline", true)
	video.Set("muted", true)

	navigator := js.Global().Get("navigator")
	mediaDevices := navigator.Get("mediaDevices")
	if !mediaDevices.Truthy() || !mediaDevices.Get("getUserMedia").Truthy() {
		cameraStatus.Set("textContent", "Camera/microphone not supported")
		micStatus.Set("textContent", "Camera/microphone not supported")
		return
	}

	startStream := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cameraBtn.Get("classList").Call("add", "hidden")
		mediaDevices.Call("getUserMedia", map[string]interface{}{
			"video": true,
			"audio": true,
		}).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			stream := args[0]
			video.Set("srcObject", stream)
			cameraStatus.Set("textContent", "Camera active")
			micStatus.Set("textContent", "Listening...")
			micVolume.Set("disabled", false)

			audioCtx := js.Global().Get("AudioContext")
			if !audioCtx.Truthy() {
				audioCtx = js.Global().Get("webkitAudioContext")
			}
			if !audioCtx.Truthy() {
				micStatus.Set("textContent", "No AudioContext support")
				return nil
			}
			ctx := audioCtx.New()
			source := ctx.Call("createMediaStreamSource", stream)
			analyser := ctx.Call("createAnalyser")
			analyser.Set("fftSize", 256)
			source.Call("connect", analyser)
			dataArray := js.Global().Get("Uint8Array").New(analyser.Get("frequencyBinCount"))

			var updateVolume js.Func
			updateVolume = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				analyser.Call("getByteTimeDomainData", dataArray)
				sum := 0.0
				length := dataArray.Get("length").Int()
				for i := 0; i < length; i++ {
					v := (dataArray.Index(i).Float() - 128.0) / 128.0
					sum += v * v
				}
				rms := 0.0
				if length > 0 {
					rms = js.Global().Get("Math").Call("sqrt", sum/float64(length)).Float()
				}
				percent := int(rms * 100 * 2)
				if percent > 100 {
					percent = 100
				}
				micVolume.Set("value", percent)
				js.Global().Call("requestAnimationFrame", updateVolume)
				return nil
			})
			jsCallbacks = append(jsCallbacks, updateVolume)
			js.Global().Call("requestAnimationFrame", updateVolume)

			// Try to play the video (required in some browsers)
			playPromise := video.Call("play")
			if playPromise.Type() == js.TypeObject {
				playPromise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					// If play() fails (likely due to autoplay policy), show the button
					cameraBtn.Get("classList").Call("remove", "hidden")
					cameraStatus.Set("textContent", "Click 'Start Camera' to begin")
					return nil
				}))
			}
			return nil
		})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			cameraStatus.Set("textContent", "Camera access denied")
			micStatus.Set("textContent", "Microphone access denied")
			cameraBtn.Get("classList").Call("remove", "hidden")
			return nil
		}))
		return nil
	})
	jsCallbacks = append(jsCallbacks, startStream)

	// Try to start stream immediately
	startStream.Invoke(js.Null(), nil)

	// Fallback: let user click the button to start camera if autoplay fails
	cameraBtn.Call("addEventListener", "click", startStream)
}

func main() {
	writeText("userAgent", js.Global().Get("navigator").Get("userAgent").String())
	writeText("language", js.Global().Get("navigator").Get("language").String())
	writeText("platform", js.Global().Get("navigator").Get("platform").String())
	writeText("timezone", js.Global().Get("Intl").Get("DateTimeFormat").New().Call("resolvedOptions").Get("timeZone").String())
	writeText("screenSize", fmt.Sprintf("%dx%d", js.Global().Get("screen").Get("width").Int(), js.Global().Get("screen").Get("height").Int()))
	writeText("gpu", getGPU())
	writeText("connection", getConnectionType())
	writeText("deviceMemory", getDeviceMemory())
	writeText("cpuCores", fmt.Sprintf("%d", js.Global().Get("navigator").Get("hardwareConcurrency").Int()))
	writeText("cookiesEnabled", fmt.Sprintf("%v", js.Global().Get("navigator").Get("cookieEnabled").Bool()))
	writeText("onlineStatus", fmt.Sprintf("%v", js.Global().Get("navigator").Get("onLine").Bool()))
	writeText("browser", fmt.Sprintf("%s %s", getBrowser(), getBrowserVersion()))
	writeText("battery", "") // will be set by fetchBattery
	writeText("preferredColorScheme", getPreferredColorScheme())
	writeText("touchSupport", getTouchSupport())
	writeText("windowSize", "") // will be set by pollWindowSize
	writeText("displayMode", getDisplayMode())
	writeText("cookies", getCookies())

	fetchLocation()
	fetchBattery()
	startClock()
	pollClipboard()
	pollWindowSize()
	fetchContacts()
	setupCameraAndMic()

	select {} // <--- This keeps the Go runtime alive!
}
