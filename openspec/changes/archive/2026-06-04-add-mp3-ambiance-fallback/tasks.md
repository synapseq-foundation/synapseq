## 1. Path and Format Resolution

- [x] 1.1 Rename `MaxWavFileSize` to `MaxAmbianceFileSize` and apply it to all ambiance audio reads.
- [x] 1.2 Add an ambiance audio format type for WAV and MP3 without changing public DSL syntax.
- [x] 1.3 Update native local ambiance resolution to try `.wav` first, then `.mp3` only when the WAV file is missing.
- [x] 1.4 Preserve existing local path validation that rejects explicit extensions, absolute paths, drive paths, backslashes, stdin, and parent traversal.
- [x] 1.5 Add errors that report both attempted local paths when neither WAV nor MP3 exists.
- [x] 1.6 Add remote ambiance format detection from URL extension, with MIME type fallback only for extensionless URLs.
- [x] 1.7 Reject remote ambiance resources with unsupported extensions or unsupported MIME types.

## 2. Decoder and Runtime

- [x] 2.1 Vendor `github.com/gopxl/beep/v2/mp3` and its required indirect decoder dependencies.
- [x] 2.2 Generalize ambiance cached source metadata to retain path, format, cached bytes, decoder state, and decoded format properties per source.
- [x] 2.3 Add decoder construction that routes WAV sources to `wav.Decode` and MP3 sources to `mp3.Decode`.
- [x] 2.4 Adapt decoder reader handling so MP3 sources use a closeable, seekable reader over cached bytes.
- [x] 2.5 Preserve independent read positions for each ambiance source index across WAV and MP3 sources.
- [x] 2.6 Preserve current EOF looping behavior by restarting MP3 playback from cached bytes when a source ends.
- [x] 2.7 Generalize sample-rate resampling so MP3 and WAV sources can both be resampled to the sequence sample rate.
- [x] 2.8 Ensure invalid existing WAV files fail with a WAV decode error and do not fall back to MP3.

## 3. Tests

- [x] 3.1 Add sequence loading tests for WAV-present, WAV-missing/MP3-present, both-missing, explicit-extension rejection, and remote URL preservation.
- [x] 3.2 Add resource or resolver tests for remote WAV/MP3 extension detection and extensionless MIME type detection.
- [x] 3.3 Add unsupported remote extension and unsupported MIME type tests.
- [x] 3.4 Add ambiance audio tests that decode and render MP3 sources.
- [x] 3.5 Add MP3 loop tests covering EOF restart and sources shorter than a render buffer.
- [x] 3.6 Add MP3 resampling coverage for mismatched source and sequence sample rates.
- [x] 3.7 Add a regression test proving an invalid existing WAV does not fall back to MP3.
- [x] 3.8 Run focused audio, sequence, resource, and core tests.

## 4. Documentation and Verification

- [x] 4.1 Update syntax documentation to describe WAV-first local fallback and MP3 as the fallback ambiance format.
- [x] 4.2 Update architecture or how-it-works documentation where ambiance format handling is described.
- [x] 4.3 Update template comments or examples if needed to clarify that local paths remain extensionless.
- [x] 4.4 Run the project test suite through the standard project command.
