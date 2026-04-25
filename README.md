# 🌐 Crompressor WASM

WebAssembly build of the [Crompressor](https://github.com/MrJc01/crompressor) compression engine.

## Overview

This module compiles the full CROM engine to WebAssembly, enabling browser-based compression and decompression without any server-side processing.

## JavaScript API

```javascript
// Load WASM module (requires wasm_exec.js from Go distribution)
const go = new Go();
const result = await WebAssembly.instantiateStreaming(
  fetch("crompressor.wasm"), go.importObject
);
go.run(result.instance);

// Compress data
const packed = cromPack(inputBytes, codebookBytes, "archive");
// packed.data    → Uint8Array (.crom bytes)
// packed.metrics → {originalSize, packedSize, hitRate, entropy, ...}

// Decompress data
const original = cromUnpack(packed.data, codebookBytes);
// original.data → Uint8Array (restored bytes)
// original.size → number

// Analyze entropy
const analysis = cromAnalyze(dataBytes);
// JSON: {entropy, size, classification, isHighEntropy}

// Version info
cromInfo(); // → "crompressor-wasm v0.2.0"
```

## Building

```bash
GOOS=js GOARCH=wasm go build -o crompressor.wasm ./cmd/wasm/
```

## Requirements

- Go 1.25+
- Local checkout of [crompressor](https://github.com/MrJc01/crompressor) (sibling directory)

## License

MIT
