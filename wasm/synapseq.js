/*/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 *
 *
 * JavaScript wrapper for SynapSeq WASM
 *
 * @class SynapSeq
 * @example
 * const synapseq = new SynapSeq();
 * await synapseq.load(content);
 * await synapseq.play();
 */
class SynapSeq {
  /**
   * Creates a new SynapSeq instance
   * @constructor
   * @param {Object} options - Configuration options
   * @param {string} [options.wasmPath='https://synapseq.org/lib/synapseq.wasm'] - Path or URL to the WASM file
   * @param {string} [options.wasmExecPath='https://synapseq.org/lib/wasm_exec.js'] - Path or URL to the wasm_exec.js file
   * @example
   * // Use the official hosted runtime (default)
   * const synapseq = new SynapSeq();
   *
   * // Use custom paths
   * const synapseq = new SynapSeq({
   *   wasmPath: './dist/synapseq.wasm',
   *   wasmExecPath: './dist/wasm_exec.js'
   * });
   *
   * // Use remote CDN
   * const synapseq = new SynapSeq({
   *   wasmPath: 'https://cdn.example.com/synapseq.wasm',
   *   wasmExecPath: 'https://cdn.example.com/wasm_exec.js'
   * });
   */
  constructor(options = {}) {
    /**
     * @private
     * @type {string}
     */
    this._wasmPath =
      options.wasmPath || "https://synapseq.org/lib/synapseq.wasm";

    /**
     * @private
     * @type {string}
     */
    this._wasmExecPath =
      options.wasmExecPath || "https://synapseq.org/lib/wasm_exec.js";

    /**
     * @private
     * @type {string|null}
     */
    this._sequence = null;

    /**
     * @private
     * @type {string|null}
     */
    this._format = "text";

    /**
     * @private
     * @type {Worker|null}
     */
    this._worker = null;

    /**
     * @private
     * @type {boolean}
     */
    this._workerReady = false;

    /**
     * @private
     * @type {string}
     */
    this._version = "unknown";

    /**
     * @private
     * @type {string}
     */
    this._buildDate = "";

    /**
     * @private
     * @type {string}
     */
    this._hash = "";

    /**
     * @private
     * @type {Promise<void>|null}
     */
    this._initPromise = null;

    /**
     * @private
     * @type {AudioContext|null}
     */
    this._audioContext = null;

    /**
     * @private
     * @type {AudioWorkletNode|null}
     */
    this._audioWorkletNode = null;

    /**
     * @private
     * @type {boolean}
     */
    this._isStreaming = false;

    /**
     * @private
     * @type {number}
     */
    this._sampleRate = 44100;

    /**
     * @private
     * @type {number}
     */
    this._playStartTime = 0;

    this._initializeWorker();
  }

  /**
   * Creates Web Worker from function
   * @private
   * @returns {Worker}
   */
  _createWorker() {
    // Worker code as a standalone function
    function workerFunction() {
      function toWorkerErrorMessage(error) {
        if (error && typeof error.message === "string") {
          return error.message;
        }

        return error || "Unknown error occurred";
      }

      let wasmReady = false;
      let wasmPath = "";
      let wasmExecPath = "";

      self.onmessage = async function (e) {
        if (e.data.type === "init") {
          wasmPath = e.data.wasmPath;
          wasmExecPath = e.data.wasmExecPath;

          try {
            const baseUrl = e.data.baseUrl || "";
            const absoluteExecPath =
              wasmExecPath.startsWith("http") || wasmExecPath.startsWith("/")
                ? wasmExecPath
                : baseUrl + wasmExecPath;
            const absoluteWasmPath =
              wasmPath.startsWith("http") || wasmPath.startsWith("/")
                ? wasmPath
                : baseUrl + wasmPath;

            importScripts(absoluteExecPath);

            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(
              fetch(absoluteWasmPath),
              go.importObject,
            );

            go.run(result.instance);
            wasmReady = true;

            self.postMessage({
              type: "ready",
              version: synapseqVersion,
              buildDate: synapseqBuildDate,
              hash: synapseqHash,
            });
          } catch (error) {
            self.postMessage({
              type: "error",
              error: "Failed to load WASM: " + toWorkerErrorMessage(error),
            });
          }
          return;
        }

        if (e.data.type === "stream") {
          if (!wasmReady) {
            self.postMessage({
              type: "error",
              error: "WASM not initialized yet",
            });
            return;
          }

          try {
            let streamErrorReported = false;
            const contentBytes = e.data.contentBytes;
            const format = e.data.format || "text";

            const onChunk = (chunk) => {
              self.postMessage(
                {
                  type: "chunk",
                  chunk: chunk,
                },
                [chunk.buffer],
              );
            };

            const onDone = () => {
              self.postMessage({
                type: "stream-done",
              });
            };

            const onError = (error) => {
              streamErrorReported = true;
              self.postMessage({
                type: "error",
                error: toWorkerErrorMessage(error),
              });
            };

            await synapseqStreamPcm(
              onChunk,
              onDone,
              onError,
              contentBytes,
              format,
            );
          } catch (error) {
            if (!streamErrorReported) {
              self.postMessage({
                type: "error",
                error: toWorkerErrorMessage(error),
              });
            }
          }
        }
      };
    }

    const workerCode = `(${workerFunction.toString()})();`;
    const blob = new Blob([workerCode], { type: "application/javascript" });
    const url = URL.createObjectURL(blob);
    const worker = new Worker(url);
    URL.revokeObjectURL(url);
    return worker;
  }

  /**
   * Normalizes worker/runtime errors into Error instances.
   * @private
   * @param {unknown} error - Error value from worker or browser APIs
   * @returns {Error}
   */
  _toError(error) {
    if (error instanceof Error) {
      return error;
    }

    if (error && typeof error.message === "string") {
      return new Error(error.message);
    }

    return new Error(String(error || "Unknown error occurred"));
  }

  /**
   * Initializes the Web Worker for WASM processing
   * @private
   * @returns {Promise<void>}
   */
  _initializeWorker() {
    this._initPromise = new Promise((resolve, reject) => {
      try {
        this._worker = this._createWorker();

        this._worker.onmessage = (e) => {
          const data = e.data;

          if (data.type === "ready") {
            this._workerReady = true;

            this._version = data.version || "unknown";
            this._buildDate = data.buildDate || "";
            this._hash = data.hash || "";

            resolve();
          } else if (data.type === "chunk") {
            this._handleAudioChunk(data.chunk);
          } else if (data.type === "stream-done") {
            this._handleStreamDone();
          } else if (data.type === "error") {
            this._handleError(this._toError(data.error));
          }
        };

        this._worker.onerror = (error) => {
          reject(new Error("Worker initialization failed: " + error.message));
        };

        // Send initialization message with paths
        const baseUrl =
          typeof window !== "undefined"
            ? window.location.href.substring(
                0,
                window.location.href.lastIndexOf("/") + 1,
              )
            : "";

        this._worker.postMessage({
          type: "init",
          wasmPath: this._wasmPath,
          wasmExecPath: this._wasmExecPath,
          baseUrl: baseUrl,
        });
      } catch (error) {
        reject(new Error("Failed to create worker: " + error.message));
      }
    });

    return this._initPromise;
  }

  /**
   * Handles audio chunk from streaming
   * @private
   * @param {Uint8Array} chunk - PCM audio chunk as Uint8Array (stereo interleaved)
   */
  _handleAudioChunk(chunk) {
    if (!this._audioWorkletNode) {
      return;
    }

    // Convert Uint8Array (Int16 bytes, stereo interleaved) to separate Float32Arrays
    const numSamples = chunk.length / 4; // 4 bytes per stereo sample (2 bytes per channel)
    const leftChannel = new Float32Array(numSamples);
    const rightChannel = new Float32Array(numSamples);

    for (let i = 0; i < numSamples; i++) {
      // Read left channel Int16 (little-endian)
      const leftByte1 = chunk[i * 4];
      const leftByte2 = chunk[i * 4 + 1];
      const leftInt16 = leftByte1 | (leftByte2 << 8);
      const signedLeft = leftInt16 > 32767 ? leftInt16 - 65536 : leftInt16;
      leftChannel[i] = signedLeft / 32768.0;

      // Read right channel Int16 (little-endian)
      const rightByte1 = chunk[i * 4 + 2];
      const rightByte2 = chunk[i * 4 + 3];
      const rightInt16 = rightByte1 | (rightByte2 << 8);
      const signedRight = rightInt16 > 32767 ? rightInt16 - 65536 : rightInt16;
      rightChannel[i] = signedRight / 32768.0;
    }

    // Send both channels to AudioWorklet
    this._audioWorkletNode.port.postMessage({
      type: "chunk",
      left: leftChannel,
      right: rightChannel,
    });
  }

  /**
   * Handles stream completion
   * @private
   */
  _handleStreamDone() {
    if (this._audioWorkletNode) {
      this._audioWorkletNode.port.postMessage({
        type: "done",
      });
    }
  }

  /**
   * Initializes AudioContext and AudioWorklet for streaming
   * @private
   * @param {number} sampleRate - Sample rate for the AudioContext
   * @returns {Promise<void>}
   */
  async _initializeAudioContext(sampleRate) {
    // Close existing context if sample rate changed
    if (this._audioContext && this._audioContext.sampleRate !== sampleRate) {
      await this._audioContext.close();
      this._audioContext = null;
    }

    if (this._audioContext) {
      return;
    }

    this._audioContext = new (window.AudioContext || window.webkitAudioContext)(
      { sampleRate: sampleRate },
    );

    // Create AudioWorklet processor inline
    const processorCode = `
      class StreamProcessor extends AudioWorkletProcessor {
        constructor() {
          super();
          this.leftChunks = [];
          this.rightChunks = [];
          this.currentLeftChunk = null;
          this.currentRightChunk = null;
          this.currentIndex = 0;
          this.done = false;
          
          this.port.onmessage = (e) => {
            if (e.data.type === 'chunk') {
              this.leftChunks.push(e.data.left);
              this.rightChunks.push(e.data.right);
            } else if (e.data.type === 'done') {
              this.done = true;
            } else if (e.data.type === 'stop') {
              this.leftChunks = [];
              this.rightChunks = [];
              this.currentLeftChunk = null;
              this.currentRightChunk = null;
              this.currentIndex = 0;
              this.done = false;
            }
          };
        }
        
        process(inputs, outputs, parameters) {
          const output = outputs[0];
          
          if (!output || output.length === 0) {
            return true;
          }
          
          const leftChannel = output[0];
          const rightChannel = output.length > 1 ? output[1] : output[0];
          
          for (let i = 0; i < leftChannel.length; i++) {
            if (!this.currentLeftChunk || this.currentIndex >= this.currentLeftChunk.length) {
              if (this.leftChunks.length > 0) {
                this.currentLeftChunk = this.leftChunks.shift();
                this.currentRightChunk = this.rightChunks.shift();
                this.currentIndex = 0;
              } else if (this.done) {
                this.port.postMessage({ type: 'ended' });
                return false;
              } else {
                leftChannel[i] = 0;
                rightChannel[i] = 0;
                continue;
              }
            }
            
            if (this.currentLeftChunk && this.currentIndex < this.currentLeftChunk.length) {
              leftChannel[i] = this.currentLeftChunk[this.currentIndex];
              rightChannel[i] = this.currentRightChunk[this.currentIndex];
              this.currentIndex++;
            } else {
              leftChannel[i] = 0;
              rightChannel[i] = 0;
            }
          }
          
          return true;
        }
      }
      
      registerProcessor('stream-processor', StreamProcessor);
    `;

    const blob = new Blob([processorCode], { type: "application/javascript" });
    const url = URL.createObjectURL(blob);

    await this._audioContext.audioWorklet.addModule(url);
    URL.revokeObjectURL(url);
  }

  /**
   * Handles errors during processing or playback
   * @private
   * @param {Error} error - The error that occurred
   */
  _handleError(error) {
    this._dispatchEvent("error", { error });
  }

  /**
   * Dispatches custom events
   * @private
   * @param {string} eventName - Name of the event
   * @param {Object} detail - Event detail data
   */
  _dispatchEvent(eventName, detail = {}) {
    if (typeof this[`on${eventName}`] === "function") {
      this[`on${eventName}`](detail);
    }
  }

  /**
   * Loads a sequence from string or File object
   * @param {string|File} input - Sequence content or File object
   * @returns {Promise<void>}
   * @throws {Error} If input is invalid or worker is not ready
   * @example
   * // Load from string
   * await synapseq.load('# Presets\nalpha\n  tone 250 isochronic 8');
   *
   * // Load from File object
   * const file = document.getElementById('fileInput').files[0];
   * await synapseq.load(file, "text");
   */
  async load(input, format = "text") {
    // Wait for worker to be ready
    if (!this._workerReady) {
      await this._initPromise;
    }

    if (!input) {
      throw new Error("Input is required");
    }

    if (format !== "text") {
      throw new Error(
        "Unsupported format: " + format + ". Only text is supported.",
      );
    }

    this._format = format;

    // Handle File object
    if (input instanceof File) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = (e) => {
          this._sequence = e.target.result;
          this._dispatchEvent("loaded");
          resolve();
        };

        reader.onerror = () => {
          reject(new Error("Failed to read file"));
        };

        reader.readAsText(input);
      });
    }

    // Handle string
    if (typeof input === "string") {
      this._sequence = input;
      this._dispatchEvent("loaded");
      return Promise.resolve();
    }

    throw new Error("Input must be a string or File object");
  }

  /**
   * Extracts sample rate from sequence content
   * @private
   * @param {string} spsq - SPSQ sequence text
   * @returns {number} Sample rate in Hz
   */
  _extractSampleRateFromText(spsq) {
    const lines = spsq.split("\n");
    for (const line of lines) {
      const trimmed = line.trim();
      // Look for @samplerate directive
      if (trimmed.startsWith("@samplerate")) {
        const match = trimmed.match(/@samplerate\s+(\d+)/);
        if (match) {
          return parseInt(match[1], 10);
        }
      }
    }
    return 44100; // Default sample rate
  }

  /**
   * Plays the loaded sequence
   * @returns {Promise<void>}
   * @throws {Error} If no sequence is loaded or worker is not ready
   * @example
   * await synapseq.play();
   */
  async play() {
    if (!this._workerReady) {
      throw new Error("Worker is not ready. Please wait for initialization.");
    }

    if (!this._sequence) {
      throw new Error("No sequence loaded. Call load() first.");
    }

    // Stop any current playback
    this.stop();

    try {
      // Get sample rate directly from SPSQ text (no need to call WASM)
      const sampleRate = this._extractSampleRateFromText(this._sequence);
      this._sampleRate = sampleRate;

      // Encode sequence to bytes for streaming
      const encoder = new TextEncoder();
      const contentBytes = encoder.encode(this._sequence);

      await this._initializeAudioContext(sampleRate);

      // Create AudioWorklet node with stereo output
      this._audioWorkletNode = new AudioWorkletNode(
        this._audioContext,
        "stream-processor",
        {
          outputChannelCount: [2],
        },
      );

      // Listen for end event
      this._audioWorkletNode.port.onmessage = (e) => {
        if (e.data.type === "ended") {
          this._dispatchEvent("ended");
          this.stop();
        }
      };

      // Connect to destination
      this._audioWorkletNode.connect(this._audioContext.destination);

      // Resume AudioContext if suspended
      if (this._audioContext.state === "suspended") {
        await this._audioContext.resume();
      }

      this._isStreaming = true;
      this._playStartTime = this._audioContext.currentTime;
      this._dispatchEvent("generating");

      // Start streaming (reuse contentBytes from above)
      this._worker.postMessage({
        type: "stream",
        contentBytes: contentBytes,
        format: this._format,
      });

      this._dispatchEvent("playing");
    } catch (error) {
      throw new Error(
        "Failed to start streaming: " + this._toError(error).message,
      );
    }
  }

  /**
   * Stops the currently playing sequence
   * @returns {void}
   * @example
   * synapseq.stop();
   */
  stop() {
    if (this._isStreaming && this._audioWorkletNode) {
      this._audioWorkletNode.port.postMessage({ type: "stop" });
      this._audioWorkletNode.disconnect();
      this._audioWorkletNode = null;
      this._isStreaming = false;
      this._playStartTime = 0;
      this._dispatchEvent("stopped");
    }
  }

  /**
   * Gets the current playback position in seconds
   * @returns {number} Current time in seconds since playback started
   * @example
   * const currentTime = synapseq.getCurrentTime();
   */
  getCurrentTime() {
    if (!this._isStreaming || !this._audioContext) {
      return 0;
    }
    return this._audioContext.currentTime - this._playStartTime;
  }

  /**
   * Gets the current playback state
   * @returns {string} One of: 'idle', 'playing', 'stopped'
   * @example
   * const state = synapseq.getState();
   */
  getState() {
    if (this._isStreaming) {
      return "playing";
    }
    return "idle";
  }

  /**
   * Checks if a sequence is currently loaded
   * @returns {boolean} True if a sequence is loaded
   * @example
   * if (synapseq.isLoaded()) {
   *   await synapseq.play();
   * }
   */
  isLoaded() {
    return this._sequence !== null;
  }

  /**
   * Gets the sample rate of the loaded sequence
   * @returns {number} Sample rate in Hz
   * @example
   * const sampleRate = synapseq.getSampleRate();
   * console.log('Sample Rate:', sampleRate, 'Hz');
   */
  getSampleRate() {
    return this._sampleRate;
  }

  /**
   * Checks if the worker is ready
   * @returns {boolean} True if worker is initialized and ready
   * @example
   * if (synapseq.isReady()) {
   *   await synapseq.load(sequence);
   * }
   */
  isReady() {
    return this._workerReady;
  }

  /**
   * Gets the SynapSeq version
   * @returns {string} The version string
   * @example
   * const version = synapseq.getVersion();
   * console.log('SynapSeq Version:', version);
   */
  async getVersion() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._version;
  }

  /**
   * Gets the build date of the SynapSeq WASM
   * @returns {string} The build date string
   * @example
   * const buildDate = synapseq.getBuildDate();
   * console.log('SynapSeq Build Date:', buildDate);
   */
  async getBuildDate() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._buildDate;
  }

  /**
   * Gets the hash of the SynapSeq WASM build
   * @returns {string} The hash string
   * @example
   * const hash = synapseq.getHash();
   * console.log('SynapSeq Hash:', hash);
   */
  async getHash() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._hash;
  }

  /**
   * Cleans up resources and terminates the worker
   * @example
   * synapseq.destroy();
   */
  destroy() {
    this.stop();

    if (this._audioContext) {
      this._audioContext.close();
      this._audioContext = null;
    }

    if (this._worker) {
      this._worker.terminate();
      this._worker = null;
    }

    this._workerReady = false;
    this._sequence = null;
    this._initPromise = null;
  }

  /**
   * Event handler called when sequence is loaded
   * @type {Function|null}
   * @example
   * synapseq.onloaded = () => console.log('Sequence loaded');
   */
  onloaded = null;

  /**
   * Event handler called when audio generation starts
   * @type {Function|null}
   * @example
   * synapseq.ongenerating = () => console.log('Generating audio...');
   */
  ongenerating = null;

  /**
   * Event handler called when playback starts
   * @type {Function|null}
   * @example
   * synapseq.onplaying = () => console.log('Now playing');
   */
  onplaying = null;

  /**
   * Event handler called when playback is stopped
   * @type {Function|null}
   * @example
   * synapseq.onstopped = () => console.log('Stopped');
   */
  onstopped = null;

  /**
   * Event handler called when playback ends naturally
   * @type {Function|null}
   * @example
   * synapseq.onended = () => console.log('Playback finished');
   */
  onended = null;

  /**
   * Event handler called when an error occurs
   * @type {Function|null}
   * @example
   * synapseq.onerror = (detail) => console.error('Error:', detail.error);
   */
  onerror = null;
}

// Export for use in modules or global scope
if (typeof module !== "undefined" && module.exports) {
  module.exports = SynapSeq;
}

// Also expose to window for browser usage
if (typeof window !== "undefined") {
  window.SynapSeq = SynapSeq;
}
