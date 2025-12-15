import { createApp } from 'vue'
import App from '@/App.vue'
import { store } from '@/store'
import type { AudioChunkMessage } from '@/messages/AudioChunkMessage';
import radioChimeStart from '@assets/audio/radio_chime_start.mp3'
import radioChimeEnd from '@assets/audio/radio_chime_end.mp3'

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string,
      RE_BACKEND_PORT: string,
    }
  }
}

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? '8282'
const reWsPort = window.APP_CONFIG?.RE_BACKEND_PORT ?? '8283'
const reconnectInterval = 1000

// TODO: Initialize audio context on user gesture to comply with browser autoplay policies
const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();

// Lazy audio initialization: created and resumed only after an explicit user gesture
let radioChimeStartBuffer: AudioBuffer;
let radioChimeEndBuffer: AudioBuffer;

const startChimePromise = loadAudioBuffer(audioContext, radioChimeStart);
const endChimePromise = loadAudioBuffer(audioContext, radioChimeEnd);
[radioChimeStartBuffer, radioChimeEndBuffer] = await Promise.all([startChimePromise, endChimePromise]);


// Create WebSocket connection
connectWebSocket(wsPort, handleWSMessage, (v) => store.wsConnected = v)
connectWebSocket(reWsPort, handleREWSMessage, (v) => store.reWsConnected = v)

function connectWebSocket(port: string, messageHandler: (event: MessageEvent) => void, setConnected?: (connected: boolean) => void) {
  /**
   * Connect to WebSocket server and set up handlers for messages, close, and errors.
   *
   * @param port The port to connect to.
   * @param messageHandler The function to handle incoming messages.
   * @param setConnected Optional function to update connection status.
   */
  console.log(`Trying to connect to WebSocket on port ${port}...`)
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    console.log('WebSocket connected')
    setConnected?.(true)
  }

  ws.onmessage = (event) => {
    messageHandler(event)
  };

  ws.onclose = function (event) {
    console.log('Socket was closed.', event.reason);
    setConnected?.(false)
    setTimeout(function () {
      connectWebSocket(port, messageHandler, setConnected);
    }, reconnectInterval);
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    ws.close()
    setConnected?.(false)
  };
}

function handleWSMessage(event: MessageEvent) {
  /**
   * Handle incoming WebSocket telemetry messages.
   *
   * @param event The WebSocket message event.
   */
  let data
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  switch (data?.topic) {
    case 'telemetry.car_telemetry': {
      // Handle car telemetry data
      const newInnerTemps = data['data']['M_carTelemetry'][0]['M_tyresInnerTemperature']
      const newOuterTemps = data['data']['M_carTelemetry'][0]['M_tyresSurfaceTemperature']

      // Mutate array in place to keep reactivity
      if (Array.isArray(newInnerTemps) && newInnerTemps.length === 4) {
        store.innerTyreTemps.splice(0, 4, ...newInnerTemps)
        console.log('Updated inner tyre temperatures:', store.innerTyreTemps)
      } else {
        console.warn('Invalid inner tyre temperature data:', newInnerTemps)
      }

      if (Array.isArray(newOuterTemps) && newOuterTemps.length === 4) {
        store.outerTyreTemps.splice(0, 4, ...newOuterTemps)
        console.log('Updated outer tyre temperatures:', store.outerTyreTemps)
      } else {
        console.warn('Invalid outer tyre temperature data:', newOuterTemps)
      }
      break
    }
    case 'telemetry.car_status': {
      // Handle car status data
      const compoundId = data['data']['M_carStatusData'][0]['M_actualTyreCompound']
      if (typeof compoundId === 'number') {
        store.actualTyreCompoundId = compoundId
      } else {
        console.warn('Invalid tyre compound ID data:', compoundId)
      }
      break
    }
    default:
      // ignore other topics
      break
  }
}

function handleREWSMessage(event: MessageEvent) {
  /**
   * Handle incoming WebSocket audio chunk messages.
   *
   * @param event The WebSocket message event.
   */
  let data: AudioChunkMessage
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  if (!data.audio_data || data.audio_data.length < 1) return;

  // Handle messages specific to the Race Engineer WebSocket
  console.log('Race Engineer WS message:', data);

  let parsedAudioData: AudioBuffer[] = [];
  data.audio_data.forEach((b64String: string, idx: number) => {
    const sampleRate = data.sample_rate[idx]
    const channels = data.channels[idx]

    // Decode base64 string to ArrayBuffer
    const arrayBuffer = base64ToArrayBuffer(b64String);

    if (audioContext) {
      parsedAudioData.push(decodePCMChunk(audioContext, arrayBuffer, sampleRate, channels));
    } else {
      console.warn('Audio context not initialized — skipping PCM decode');
    }
  });

  let fullAudioBuffer = [
    radioChimeStartBuffer,
    ...parsedAudioData,
    radioChimeEndBuffer
  ].filter((b): b is AudioBuffer => !!b);

  if (audioContext) {
    playAudioBuffers(audioContext, fullAudioBuffer);
  } else {
    console.warn('Audio context not initialized — skipping playback');
  }
}

async function loadAudioBuffer(audioContext: AudioContext, url: string): Promise<AudioBuffer> {
  /**
   * Loads an audio file from a URL and decodes it into an AudioBuffer.
   *
   * @param audioContext The AudioContext to use for decoding.
   * @param url The URL of the audio file to load.
   * @returns A Promise that resolves to the decoded AudioBuffer.
   */
  const response = await fetch(url);

  if (!response.ok) {
    return audioContext.createBuffer(1, audioContext.sampleRate, audioContext.sampleRate);
  }

  const arrayBuffer = await response.arrayBuffer();

  const audioBuffer = await audioContext.decodeAudioData(arrayBuffer);

  return audioBuffer;
}

function playAudioBuffers(audioContext: AudioContext, audioBuffers: AudioBuffer[]) {
  /**
   * Plays a sequence of AudioBuffers one after another.
   *
   * @param audioContext The AudioContext to use for playback.
   * @param audioBuffers An array of AudioBuffers to play in sequence.
   */
  let source = audioContext.createBufferSource();
  source.buffer = audioBuffers.shift()!;
  source.connect(audioContext.destination);
  source.start(0);

  source.onended = () => { if (audioContext) playAudioBuffers(audioContext, audioBuffers); };
}

function base64ToArrayBuffer(b64string: string): ArrayBuffer {
  /**
   * Decodes a base64 string to an ArrayBuffer.
   *
   * @param b64string The base64 encoded string.
   * @returns The decoded ArrayBuffer.
   */
  const binString = atob(b64string);
  const byteArray = new Uint8Array(binString.length);
  for (let i = 0; i < binString.length; i++) {
    byteArray[i] = binString.charCodeAt(i);
  }

  return byteArray.buffer;
}

function decodePCMChunk(audioContext: AudioContext,
  arrayBuffer: ArrayBuffer,
  sampleRate: number,
  channels: number): AudioBuffer {
  /**
   * Decodes a PCM audio chunk from an ArrayBuffer to an AudioBuffer.
   *
   * @param audioContext The AudioContext to use for decoding.
   * @param arrayBuffer The ArrayBuffer containing PCM audio data.
   * @param sampleRate The sample rate of the audio data.
   * @param channels The number of audio channels.
   * @returns The decoded AudioBuffer.
   */
  // Backend uses Int16 PCM
  const int16Array = new Int16Array(arrayBuffer);

  // Normalize Int16 to Float32 range [-1.0, 1.0] for the AudioBuffer
  const maxInt16 = Math.pow(2, 15);
  const float32Array = new Float32Array(int16Array.length);
  for (let i = 0; i < int16Array.length; i++) {
    float32Array[i] = int16Array[i] / maxInt16;
  }

  const audioBuffer = audioContext.createBuffer(
    channels,
    float32Array.length / channels,
    sampleRate
  );

  for (let c = 0; c < channels; c++) {
    const channelData = audioBuffer.getChannelData(c);
    for (let i = 0; i < channelData.length; i++) {
      channelData[i] = float32Array[i * channels + c];
    }
  }

  return audioBuffer;
}


createApp(App).mount('#app')
