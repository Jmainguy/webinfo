package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

// Write text to an element by id, with error logging if not found.
func writeText(id string, text string) {
	document := js.Global().Get("document")
	element := document.Call("getElementById", id)
	if element.Truthy() {
		element.Set("innerText", text)
	} else {
		console := js.Global().Get("console")
		console.Call("error", "Element not found:", id)
	}
}

// Get the user agent string from the browser.
func getUserAgent() string {
	navigator := js.Global().Get("navigator")
	userAgent := navigator.Get("userAgent")
	return userAgent.String()
}

// Get the browser language.
func getLanguage() string {
	navigator := js.Global().Get("navigator")
	language := navigator.Get("language")
	return language.String()
}

// Get the platform string.
func getPlatform() string {
	navigator := js.Global().Get("navigator")
	platform := navigator.Get("platform")
	return platform.String()
}

// Get the timezone string.
func getTimezone() string {
	intl := js.Global().Get("Intl")
	dateTimeFormat := intl.Get("DateTimeFormat").New()
	resolvedOptions := dateTimeFormat.Call("resolvedOptions")
	timeZone := resolvedOptions.Get("timeZone")
	js.Global().Get("console").Call("log", "getTimezone", timeZone)
	return timeZone.String()
}

// Get the screen size as a string.
func getScreenSize() string {
	screen := js.Global().Get("screen")
	width := screen.Get("width")
	height := screen.Get("height")
	js.Global().Get("console").Call("log", "getScreenSize", width, height)
	return fmt.Sprintf("%dx%d", width.Int(), height.Int())
}

// Get the GPU renderer string using WebGL.
func getGPU() string {
	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	var gl js.Value
	contexts := []string{"webgl2", "webgl", "experimental-webgl"}
	for i := 0; i < len(contexts); i++ {
		gl = canvas.Call("getContext", contexts[i])
		if gl.Truthy() {
			break
		}
	}
	if !gl.Truthy() {
		return "No WebGL"
	}
	debugInfo := gl.Call("getExtension", "WEBGL_debug_renderer_info")
	if debugInfo.Truthy() {
		renderer := gl.Call("getParameter", debugInfo.Get("UNMASKED_RENDERER_WEBGL"))
		if renderer.Type() == js.TypeString && renderer.String() != "" {
			return renderer.String()
		}
		vendor := gl.Call("getParameter", debugInfo.Get("UNMASKED_VENDOR_WEBGL"))
		if vendor.Type() == js.TypeString && vendor.String() != "" {
			return vendor.String()
		}
	}
	return "Unknown GPU"
}

// Get the connection type string.
func getConnectionType() string {
	navigator := js.Global().Get("navigator")
	connection := navigator.Get("connection")
	if connection.Truthy() {
		effectiveType := connection.Get("effectiveType")
		if effectiveType.Truthy() && effectiveType.Type() == js.TypeString {
			return effectiveType.String()
		}
		typ := connection.Get("type")
		if typ.Truthy() && typ.Type() == js.TypeString {
			return typ.String()
		}
		return "Unknown"
	}
	return "Not Supported"
}

// Get the device memory as a string.
func getDeviceMemory() string {
	navigator := js.Global().Get("navigator")
	mem := navigator.Get("deviceMemory")
	if mem.Truthy() && mem.Type() == js.TypeNumber {
		val := mem.Float()
		if val > 0 {
			return fmt.Sprintf("%.0f GB", val)
		}
	}
	return "Not Supported"
}

// Get the number of hardware concurrency (CPU cores) as a string.
func getHardwareConcurrency() string {
	navigator := js.Global().Get("navigator")
	cores := navigator.Get("hardwareConcurrency")
	if cores.Truthy() {
		return fmt.Sprintf("%d", cores.Int())
	}
	return "Unknown"
}

// Get whether cookies are enabled as a string.
func getCookiesEnabled() string {
	navigator := js.Global().Get("navigator")
	enabled := navigator.Get("cookieEnabled")
	if enabled.Truthy() {
		return fmt.Sprintf("%v", enabled.Bool())
	}
	return "Unknown"
}

// Get the online status as a string.
func getOnlineStatus() string {
	navigator := js.Global().Get("navigator")
	online := navigator.Get("onLine")
	if online.Truthy() {
		return fmt.Sprintf("%v", online.Bool())
	}
	return "Unknown"
}

// Get the browser name using user agent and userAgentData.
func getBrowser() string {
	navigator := js.Global().Get("navigator")
	ua := getUserAgent()
	if navigator.Get("brave").Truthy() {
		return "Brave (Chromium-based)"
	}
	if navigator.Get("userAgentData").Truthy() {
		brands := navigator.Get("userAgentData").Get("brands")
		if brands.Length() > 0 {
			brand := brands.Index(0).Get("brand").String()
			version := brands.Index(0).Get("version").String()
			return brand + " " + version
		}
	}
	if strings.Contains(ua, "Edg/") {
		return "Microsoft Edge"
	}
	if strings.Contains(ua, "OPR/") {
		return "Opera"
	}
	if strings.Contains(ua, "Chrome/") {
		return "Chrome"
	}
	if strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/") {
		return "Safari"
	}
	if strings.Contains(ua, "Firefox/") {
		return "Firefox"
	}
	return "Unknown"
}

// Get the browser version string.
func getBrowserVersion() string {
	ua := getUserAgent()
	if strings.Contains(ua, "Edg/") {
		return extractVersionGo(ua, "Edg/")
	}
	if strings.Contains(ua, "OPR/") {
		return extractVersionGo(ua, "OPR/")
	}
	if strings.Contains(ua, "Chrome/") {
		return extractVersionGo(ua, "Chrome/")
	}
	if strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/") {
		return extractVersionGo(ua, "Version/")
	}
	if strings.Contains(ua, "Firefox/") {
		return extractVersionGo(ua, "Firefox/")
	}
	return ""
}

// Helper: check if s contains substr.
func contains(s string, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (indexOf(s, substr) != -1)))
}

// Helper: index of substr in s using JS String.indexOf.
func indexOf(s string, substr string) int {
	return js.Global().Get("String").New(s).Call("indexOf", substr).Int()
}

// Extract version string from user agent.
func extractVersion(ua string, marker string) string {
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

// Extract version string from user agent using Go's strings.Index.
func extractVersionGo(ua string, marker string) string {
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

// Get the browser history length as a string.
func getHistoryLength() string {
	history := js.Global().Get("history")
	length := history.Get("length")
	return fmt.Sprintf("%d", length.Int())
}

// Fetch the user's location and update the UI (no map).
func fetchLocation() {
	successCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pos := args[0]
		coords := pos.Get("coords")
		lat := coords.Get("latitude").Float()
		lon := coords.Get("longitude").Float()
		writeText("location", fmt.Sprintf("Latitude: %.5f, Longitude: %.5f", lat, lon))
		return nil
	})
	errorCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		writeText("location", "Location unavailable")
		return nil
	})
	navigator := js.Global().Get("navigator")
	geolocation := navigator.Get("geolocation")
	geolocation.Call("getCurrentPosition", successCallback, errorCallback)
}

// Fetch battery info and update the UI.
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

// Start a digital clock that updates every second.
func startClock() {
	go func() {
		for {
			now := time.Now()
			writeText("clock", now.Format("15:04:05"))
			time.Sleep(time.Second)
		}
	}()
}

// Fetch contacts info and update the UI.
func fetchContacts() {
	navigator := js.Global().Get("navigator")
	contacts := navigator.Get("contacts")
	if contacts.Truthy() && contacts.Get("select").Type() == js.TypeFunction {
		writeText("contacts", "Contacts API available (requires user gesture)")
	} else {
		writeText("contacts", "Contacts API not available")
	}
}

// Get the preferred color scheme.
func getPreferredColorScheme() string {
	window := js.Global().Get("window")
	mql := window.Call("matchMedia", "(prefers-color-scheme: dark)")
	if mql.Truthy() && mql.Get("matches").Bool() {
		return "Dark"
	}
	return "Light"
}

// Get whether touch is supported.
func getTouchSupport() string {
	window := js.Global().Get("window")
	if window.Get("ontouchstart").Type() != js.TypeUndefined {
		return "Yes"
	}
	navigator := js.Global().Get("navigator")
	maxTouch := navigator.Get("maxTouchPoints")
	if maxTouch.Truthy() && maxTouch.Int() > 0 {
		return "Yes"
	}
	return "No"
}

// Poll the window size and update the UI if it changes.
var lastWindowSize string

func pollWindowSize() {
	go func() {
		for {
			window := js.Global().Get("window")
			width := window.Get("innerWidth")
			height := window.Get("innerHeight")
			size := "Unknown"
			if width.Truthy() && height.Truthy() {
				size = fmt.Sprintf("%dx%d", width.Int(), height.Int())
			}
			if size != lastWindowSize {
				lastWindowSize = size
				writeText("windowSize", size)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// Get the display mode (standalone/PWA or browser tab).
func getDisplayMode() string {
	window := js.Global().Get("window")
	mql := window.Call("matchMedia", "(display-mode: standalone)")
	if mql.Truthy() && mql.Get("matches").Bool() {
		return "Standalone/PWA"
	}
	return "Browser Tab"
}

// Get the cookies string.
func getCookies() string {
	document := js.Global().Get("document")
	cookies := document.Get("cookie").String()
	if cookies == "" {
		return "None"
	}
	return cookies
}

// Keep references to prevent GC
var jsCallbacks []js.Func

// Setup camera and microphone, and update the UI.
func setupCameraAndMic() {
	document := js.Global().Get("document")
	video := document.Call("getElementById", "cameraVideo")
	cameraStatus := document.Call("getElementById", "cameraStatus")
	micVolume := document.Call("getElementById", "micVolume")
	micStatus := document.Call("getElementById", "micStatus")
	cameraBtn := document.Call("getElementById", "cameraStartBtn")

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
			dataArray := js.Global().Get("Uint8Array").New(analyser.Get("frequencyBinCount"))
			source.Call("connect", analyser)

			var updateVolume js.Func
			updateVolume = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				if js.Global().Get("wasmGoExited").Truthy() && js.Global().Get("wasmGoExited").Bool() {
					updateVolume.Release()
					return nil
				}
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
				if js.Global().Get("wasmGoExited").Truthy() && js.Global().Get("wasmGoExited").Bool() {
					updateVolume.Release()
					return nil
				}
				js.Global().Call("requestAnimationFrame", updateVolume)
				return nil
			})
			jsCallbacks = append(jsCallbacks, updateVolume)
			js.Global().Call("requestAnimationFrame", updateVolume)

			playPromise := video.Call("play")
			if playPromise.Type() == js.TypeObject {
				playPromise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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
	startStream.Invoke(js.Null(), nil)
	cameraBtn.Call("addEventListener", "click", startStream)
}

// Clipboard API (polling, async)
var lastClipboard string
var clipboardErrorCount int

// Poll the clipboard for changes and update the UI.
func pollClipboard() {
	navigator := js.Global().Get("navigator")
	clipboard := navigator.Get("clipboard")
	document := js.Global().Get("document")
	if clipboard.Truthy() && clipboard.Get("readText").Type() == js.TypeFunction {
		go func() {
			for {
				// Stop polling if Go program has exited
				if js.Global().Get("wasmGoExited").Truthy() && js.Global().Get("wasmGoExited").Bool() {
					return
				}
				if document.Truthy() && document.Get("hasFocus").Type() == js.TypeFunction && document.Call("hasFocus").Bool() {
					var then, catch js.Func
					then = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						// Stop if Go program has exited
						if js.Global().Get("wasmGoExited").Truthy() && js.Global().Get("wasmGoExited").Bool() {
							then.Release()
							catch.Release()
							return nil
						}
						text := args[0].String()
						clipboardErrorCount = 0
						if text != lastClipboard {
							lastClipboard = text
							writeText("clipboard", text)
						}
						then.Release()
						catch.Release()
						return nil
					})
					catch = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						// Stop if Go program has exited
						if js.Global().Get("wasmGoExited").Truthy() && js.Global().Get("wasmGoExited").Bool() {
							then.Release()
							catch.Release()
							return nil
						}
						clipboardErrorCount++
						if clipboardErrorCount > 3 {
							writeText("clipboard", "Clipboard access denied or unavailable")
						}
						then.Release()
						catch.Release()
						return nil
					})
					promise := clipboard.Call("readText")
					promise.Call("then", then).Call("catch", catch)
					time.Sleep(500 * time.Millisecond)
				} else {
					writeText("clipboard", "")
					time.Sleep(500 * time.Millisecond)
				}
			}
		}()
	} else {
		writeText("clipboard", "Clipboard API not available")
	}
}

// Register the clipboard poller as a JS function.
func registerClipboardPoller() {
	js.Global().Set("startClipboardPolling", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go pollClipboard()
		return nil
	}))
}

// Register the clipboard permission button handler.
func registerClipboardButton() {
	document := js.Global().Get("document")
	btn := document.Call("getElementById", "clipboardPermBtn")
	if !btn.Truthy() {
		return
	}
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		navigator := js.Global().Get("navigator")
		clipboard := navigator.Get("clipboard")
		clipboardDiv := document.Call("getElementById", "clipboard")
		if clipboard.Truthy() && clipboard.Get("readText").Type() == js.TypeFunction {
			var then, catch js.Func
			then = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				text := args[0].String()
				clipboardDiv.Set("innerText", text)
				if js.Global().Get("startClipboardPolling").Type() == js.TypeFunction {
					js.Global().Call("startClipboardPolling")
				}
				then.Release()
				catch.Release()
				return nil
			})
			catch = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				clipboardDiv.Set("innerText", "Clipboard access denied or unavailable")
				then.Release()
				catch.Release()
				return nil
			})
			promise := clipboard.Call("readText")
			promise.Call("then", then).Call("catch", catch)
		} else {
			clipboardDiv.Set("innerText", "Clipboard API not available")
		}
		return nil
	})
	btn.Call("addEventListener", "click", handler)
	jsCallbacks = append(jsCallbacks, handler)
}

// Add: Register location permission button
func registerLocationButton() {
	document := js.Global().Get("document")
	btn := document.Call("getElementById", "locationPermBtn")
	if !btn.Truthy() {
		return
	}
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fetchLocation()
		return nil
	})
	btn.Call("addEventListener", "click", handler)
	jsCallbacks = append(jsCallbacks, handler)
}

// Main entry point.
func main() {
	writeText("userAgent", getUserAgent())
	writeText("language", getLanguage())
	writeText("platform", getPlatform())
	writeText("timezone", getTimezone())
	writeText("screenSize", getScreenSize())
	writeText("gpu", getGPU())
	writeText("connection", getConnectionType())
	writeText("deviceMemory", getDeviceMemory())
	writeText("cpuCores", getHardwareConcurrency())
	writeText("cookiesEnabled", getCookiesEnabled())
	writeText("onlineStatus", getOnlineStatus())
	writeText("browser", fmt.Sprintf("%s %s", getBrowser(), getBrowserVersion()))
	writeText("battery", "")
	writeText("preferredColorScheme", getPreferredColorScheme())
	writeText("touchSupport", getTouchSupport())
	writeText("windowSize", "")
	writeText("displayMode", getDisplayMode())
	writeText("cookies", getCookies())
	writeText("clipboard", "Click the button to request access")

	js.Global().Set("wasmGoExited", false)
	defer js.Global().Set("wasmGoExited", true)

	fetchLocation()
	fetchBattery()
	startClock()
	pollWindowSize()
	fetchContacts()
	setupCameraAndMic()
	registerClipboardPoller()
	registerClipboardButton()
	registerLocationButton()

	select {}
}
