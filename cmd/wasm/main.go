//go:build js && wasm

// Package main implements the Crompressor WASM module.
//
// It exposes three JavaScript functions:
//   - cromPack(inputBytes, codebookBytes, mode) → packed bytes
//   - cromUnpack(packedBytes, codebookBytes) → original bytes
//   - cromAnalyze(inputBytes) → JSON metrics string
//   - cromInfo() → version string
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/MrJc01/crompressor/pkg/entropy"
	"github.com/MrJc01/crompressor/pkg/cromlib"
)

func main() {
	fmt.Println("Crompressor WASM module loaded — engine wired")

	js.Global().Set("cromPack", js.FuncOf(pack))
	js.Global().Set("cromUnpack", js.FuncOf(unpack))
	js.Global().Set("cromAnalyze", js.FuncOf(analyze))
	js.Global().Set("cromInfo", js.FuncOf(info))

	// Block forever
	select {}
}

// pack receives (Uint8Array input, Uint8Array codebook, string mode) from JS
// and returns the packed .crom bytes as a Uint8Array.
func pack(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return jsError("cromPack: requires at least 2 arguments (input, codebook)")
	}

	input := jsToBytes(args[0])
	codebookData := jsToBytes(args[1])

	mode := "archive"
	if len(args) > 2 && args[2].Type() == js.TypeString {
		mode = args[2].String()
	}

	opts := cromlib.DefaultPackOptions()
	opts.Mode = mode

	packed, metrics, err := cromlib.PackBytes(input, codebookData, opts)
	if err != nil {
		return jsError(fmt.Sprintf("cromPack: %v", err))
	}

	// Return an object with { data: Uint8Array, metrics: {...} }
	result := js.Global().Get("Object").New()

	dataArray := js.Global().Get("Uint8Array").New(len(packed))
	js.CopyBytesToJS(dataArray, packed)
	result.Set("data", dataArray)

	metricsObj := js.Global().Get("Object").New()
	metricsObj.Set("originalSize", metrics.OriginalSize)
	metricsObj.Set("packedSize", metrics.PackedSize)
	metricsObj.Set("durationMs", metrics.Duration.Milliseconds())
	metricsObj.Set("hitRate", metrics.HitRate)
	metricsObj.Set("literalChunks", metrics.LiteralChunks)
	metricsObj.Set("totalChunks", metrics.TotalChunks)
	metricsObj.Set("avgSimilarity", metrics.AvgSimilarity)
	metricsObj.Set("entropy", metrics.Entropy)
	result.Set("metrics", metricsObj)

	return result
}

// unpack receives (Uint8Array packed, Uint8Array codebook) from JS
// and returns the reconstructed original bytes as a Uint8Array.
func unpack(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return jsError("cromUnpack: requires 2 arguments (cromData, codebook)")
	}

	cromData := jsToBytes(args[0])
	codebookData := jsToBytes(args[1])

	opts := cromlib.DefaultUnpackOptions()

	original, err := cromlib.UnpackBytes(cromData, codebookData, opts)
	if err != nil {
		return jsError(fmt.Sprintf("cromUnpack: %v", err))
	}

	result := js.Global().Get("Object").New()
	dataArray := js.Global().Get("Uint8Array").New(len(original))
	js.CopyBytesToJS(dataArray, original)
	result.Set("data", dataArray)
	result.Set("size", len(original))

	return result
}

// analyze receives a Uint8Array and returns entropy analysis as a JSON string.
func analyze(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("cromAnalyze: requires 1 argument (data)")
	}

	data := jsToBytes(args[0])
	if len(data) == 0 {
		return jsError("cromAnalyze: empty input")
	}

	ent := entropy.Shannon(data)

	result := map[string]interface{}{
		"entropy":     ent,
		"size":        len(data),
		"isHighEntropy": ent > 7.0,
		"isLowEntropy":  ent < 3.0,
		"classification": classifyEntropy(ent),
	}

	jsonBytes, _ := json.Marshal(result)
	return js.ValueOf(string(jsonBytes))
}

func info(this js.Value, args []js.Value) interface{} {
	return js.ValueOf("crompressor-wasm v0.2.0 — engine wired to pkg/cromlib.PackBytes")
}

// --- Helpers ---

func jsToBytes(val js.Value) []byte {
	length := val.Get("length").Int()
	buf := make([]byte, length)
	js.CopyBytesToGo(buf, val)
	return buf
}

func jsError(msg string) interface{} {
	result := js.Global().Get("Object").New()
	result.Set("error", msg)
	return result
}

func classifyEntropy(e float64) string {
	switch {
	case e < 1.0:
		return "trivial"
	case e < 3.0:
		return "low"
	case e < 5.0:
		return "medium"
	case e < 7.0:
		return "high"
	default:
		return "random"
	}
}
