# go-whisper-cpp-server-example

Example using Silero VAD and whisper.cpp for speech recognition using go. For the article [Local, all-in-one Go speech-to-text solution with Silero VAD and whisper.cpp server](https://medium.com/@etolkachev93/local-all-in-one-go-speech-to-text-solution-with-silero-vad-and-whisper-cpp-server-94a69fa51b04).

## Dependencies (for mac)

* Install whisper.cpp
* Download the whisper model converted to ggml format: [ggerganov/whisper.cpp](https://huggingface.co/ggerganov/whisper.cpp)
* Install onnxruntime: `brew install onnxruntime`

## How to run

```bash
export LIBRARY_PATH=/opt/homebrew/Cellar/onnxruntime/1.17.1/lib
C_INCLUDE_PATH=/opt/homebrew/Cellar/onnxruntime/1.17.1/include/onnxruntime go run main.go
```