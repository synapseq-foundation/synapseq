# SynapSeq JavaScript API Reference

Browser-facing JavaScript wrapper for the SynapSeq WebAssembly runtime.

`SynapSeq` loads the Go WASM runtime in a Web Worker, streams generated PCM audio to the main thread, and plays it through the Web Audio API.

## Class: `SynapSeq`

### Constructor

```js
const synapseq = new SynapSeq(options);
```

Creates a new SynapSeq instance and starts worker initialization immediately.

Parameters:

- `options?: object`
- `options.wasmPath?: string` - Path or URL to `synapseq.wasm`. Default: `https://synapseq.org/lib/synapseq.wasm`
- `options.wasmExecPath?: string` - Path or URL to `wasm_exec.js`. Default: `https://synapseq.org/lib/wasm_exec.js`

Default behavior:

- If no options are provided, the instance fetches the latest official runtime assets from `synapseq.org`.
- If custom paths are provided, absolute URLs are used as-is.
- Relative custom paths are resolved against `window.location.href`.

Example:

```js
const synapseq = new SynapSeq();
```

Custom assets:

```js
const synapseq = new SynapSeq({
  wasmPath: "./dist/synapseq.wasm",
  wasmExecPath: "./dist/wasm_exec.js",
});
```

## Methods

### `load(input, format = "text")`

```js
await synapseq.load(input, format);
```

Loads a sequence into memory.

Parameters:

- `input: string | File` - Sequence content as text or a browser `File` object.
- `format?: string` - Accepted only as `"text"`. Kept for compatibility.

Returns:

- `Promise<void>`

Behavior:

- Waits for worker initialization if needed.
- Stores the sequence for future playback.
- Dispatches `onloaded` after a successful load.

Throws or rejects:

- `Error("Input is required")` when `input` is missing.
- `Error("Unsupported format: ... Only text is supported.")` when `format !== "text"`.
- `Error("Failed to read file")` when file reading fails.
- `Error("Input must be a string or File object")` when `input` has an unsupported type.

### `play()`

```js
await synapseq.play();
```

Starts playback of the loaded sequence.

Returns:

- `Promise<void>`

Behavior:

- Requires the worker to already be ready.
- Requires a previously loaded sequence.
- Stops current playback before starting a new one.
- Extracts sample rate from the text sequence using the `@samplerate` directive.
- Defaults to `44100` when no sample rate is declared.
- Creates or recreates the `AudioContext` when required.
- Creates an `AudioWorkletNode` for streaming playback.
- Dispatches `ongenerating` before streaming starts.
- Dispatches `onplaying` after playback startup is triggered.
- Dispatches `onended` when playback finishes naturally.

Throws:

- `Error("Worker is not ready. Please wait for initialization.")`
- `Error("No sequence loaded. Call load() first.")`
- `Error("Failed to start streaming: ...")`

### `stop()`

```js
synapseq.stop();
```

Stops the current playback session.

Returns:

- `void`

Behavior:

- Stops the active stream when playback is running.
- Disconnects the current `AudioWorkletNode`.
- Resets playback timing state.
- Dispatches `onstopped`.

Notes:

- Does not terminate the worker.
- Does not unload the sequence.

### `destroy()`

```js
synapseq.destroy();
```

Releases runtime resources owned by the instance.

Returns:

- `void`

Behavior:

- Calls `stop()`.
- Closes the current `AudioContext`.
- Terminates the worker.
- Clears internal state such as readiness and loaded content.

Use this when the instance will no longer be reused.

### `getCurrentTime()`

```js
const seconds = synapseq.getCurrentTime();
```

Returns the elapsed playback time in seconds.

Returns:

- `number`

Notes:

- Returns `0` when not streaming.

### `getState()`

```js
const state = synapseq.getState();
```

Returns the current playback state.

Returns:

- `"playing" | "idle"`

### `isLoaded()`

```js
const loaded = synapseq.isLoaded();
```

Returns whether a sequence is currently loaded.

Returns:

- `boolean`

### `getSampleRate()`

```js
const sampleRate = synapseq.getSampleRate();
```

Returns the sample rate associated with the currently loaded or active sequence.

Returns:

- `number`

Notes:

- Defaults to `44100` before playback.

### `isReady()`

```js
const ready = synapseq.isReady();
```

Returns whether the worker has finished WASM initialization.

Returns:

- `boolean`

### `getVersion()`

```js
const version = await synapseq.getVersion();
```

Returns the SynapSeq version string exported by the WASM runtime.

Returns:

- `Promise<string>`

Notes:

- Waits for initialization if needed.

### `getBuildDate()`

```js
const buildDate = await synapseq.getBuildDate();
```

Returns the build date string exported by the WASM runtime.

Returns:

- `Promise<string>`

Notes:

- Waits for initialization if needed.

### `getHash()`

```js
const hash = await synapseq.getHash();
```

Returns the Git commit hash exported by the WASM runtime.

Returns:

- `Promise<string>`

Notes:

- Waits for initialization if needed.

## Event Callback Properties

The wrapper exposes callback properties instead of DOM events.

### `onloaded`

```js
synapseq.onloaded = () => {};
```

Called after a sequence is successfully loaded.

### `ongenerating`

```js
synapseq.ongenerating = () => {};
```

Called when audio generation is about to start.

### `onplaying`

```js
synapseq.onplaying = () => {};
```

Called after playback startup is triggered.

### `onstopped`

```js
synapseq.onstopped = () => {};
```

Called when playback is stopped explicitly.

### `onended`

```js
synapseq.onended = () => {};
```

Called when playback ends naturally.

### `onerror`

```js
synapseq.onerror = ({ error }) => {};
```

Called when the wrapper forwards a runtime or playback error.

Callback payload:

- `{ error: Error }`

## Input Support

Supported input:

- SynapSeq text sequences as `string`
- Browser `File` objects read as text

Unsupported input:

- JSON sequence payloads
- Binary sequence formats

## Runtime Notes

### Initialization model

- Worker initialization starts in the constructor.
- `load()` waits for initialization automatically.
- `play()` requires the worker to already be ready.

### Playback model

- Audio generation runs in a `Worker`.
- PCM chunks are streamed back to the main thread.
- Playback uses `AudioContext` and `AudioWorklet`.

### Sample rate handling

- `@samplerate <number>` in the text sequence defines the playback sample rate.
- If omitted, `44100` Hz is used.
- If the sample rate changes, the current `AudioContext` is closed and recreated.

### Browser requirements

The wrapper depends on:

- `WebAssembly`
- `Worker`
- `Blob`
- `URL.createObjectURL`
- `AudioContext` or `webkitAudioContext`
- `AudioWorklet`
- `FileReader`
- `TextEncoder`

### Deployment notes

- Use HTTP or HTTPS, not a raw `file://` page.
- If you self-host `synapseq.wasm`, serve it as `application/wasm`.
- If you use custom cross-origin asset URLs, they must allow CORS.
- Start playback from a user gesture when required by the browser.

### Limitations

- Browser-only playback model
- Text sequences only
- Callback properties instead of `EventTarget`
- No seek support
- No pause/resume API
- No multiple simultaneous playback nodes within one instance

## Examples

### Basic playback

```js
const synapseq = new SynapSeq();

await synapseq.load(`@samplerate 44100
@volume 75

alpha
  noise pink amplitude 20
  tone 250 isochronic 8 amplitude 12

00:00:00 silence
00:00:10 alpha
00:10:00 alpha ease-out
00:10:20 silence`);

await synapseq.play();
```

### Using a local runtime build

```js
const synapseq = new SynapSeq({
  wasmPath: "./dist/synapseq.wasm",
  wasmExecPath: "./dist/wasm_exec.js",
});

await synapseq.load(`alpha
  noise pink amplitude 18
  tone 220 binaural 10 amplitude 10

00:00:00 silence
00:00:15 alpha
00:08:00 alpha ease-out
00:08:20 silence`);

await synapseq.play();
```

### Loading from a file input

```js
const synapseq = new SynapSeq();
const file = document.getElementById("fileInput").files[0];

await synapseq.load(file, "text");
await synapseq.play();
```

### Registering callbacks

```js
const synapseq = new SynapSeq();

synapseq.onloaded = () => {
  console.log("Sequence loaded");
};

synapseq.onplaying = () => {
  console.log("Playback started");
};

synapseq.onended = () => {
  console.log("Playback finished");
};

synapseq.onerror = ({ error }) => {
  console.error(error);
};
```

Interactive example: see `example.html` in this directory.
It uses the local files `./synapseq.wasm` and `./wasm_exec.js`.
If those files are missing, run `make wasm` from the repository root before opening the example.
