import { ref } from 'vue';
import type { AudioChunkMessage } from '@/messages/AudioChunkMessage';

import radioChimeStart from '@assets/audio/radio_chime_start.mp3'
import radioChimeEnd from '@assets/audio/radio_chime_end.mp3'

class AudioPlayer {
    private static instance: AudioPlayer;
    private audioContext: AudioContext | null = null;
    private queue: AudioBuffer[] = [];
    private isPlaying: boolean = false;
    private radioChimeStartBuffer: AudioBuffer | null = null;
    private radioChimeEndBuffer: AudioBuffer | null = null;
    private currentSource: AudioBufferSourceNode | null = null;
    private chimeBuffersLoadedPromise: Promise<void> | null = null;

    public isAudioReady = ref(false);

    private constructor() { }

    public static getInstance(): AudioPlayer {
        if (!AudioPlayer.instance) {
            AudioPlayer.instance = new AudioPlayer();
        }
        return AudioPlayer.instance;
    }

    public async initialize(): Promise<void> {
        /**
         * Initializes the AudioPlayer by creating and resuming the AudioContext and loading chime buffers.
         */
        try {
            if (!this.audioContext) {
                this.audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
                this.isAudioReady.value = true;
            }

            if (this.audioContext.state === 'suspended') {
                await this.audioContext.resume();
            }
        } catch (error) {
            console.error('Error initializing AudioContext:', error);
            this.isAudioReady.value = false;
            return;
        }

        if (!this.chimeBuffersLoadedPromise) {
            this.chimeBuffersLoadedPromise = (async () => {
                const startChimePromise = this.loadAudioBuffer(this.audioContext!, radioChimeStart);
                const endChimePromise = this.loadAudioBuffer(this.audioContext!, radioChimeEnd);
                [this.radioChimeStartBuffer, this.radioChimeEndBuffer] = await Promise.all([startChimePromise, endChimePromise]);
            })();
    }

    public async playAudio(audioData: AudioChunkMessage, playChimes: boolean = true): Promise<void> {
        /**
         * Plays an audio message composed of multiple audio chunks.
         * 
         * @param audioData The AudioChunkMessage containing audio data.
         * @param playChimes Whether to play chimes before and after the message.
         */
        const context = this.audioContext;

        if (!context) {
            console.warn('AudioContext is not initialized. Call initialize() first.');
            return;
        }

        let parsedAudioData: AudioBuffer[] = [];
        audioData.audio_data.forEach((b64String, idx) => {
            const sampleRate = audioData.sample_rate[idx]
            const channels = audioData.channels[idx]

            // Decode base64 string to ArrayBuffer
            const arrayBuffer = this.base64ToArrayBuffer(b64String);

            parsedAudioData.push(this.decodePCMChunk(context, arrayBuffer, sampleRate, channels));
        });

        if (playChimes) {
            parsedAudioData = [
            this.radioChimeStartBuffer,
            ...parsedAudioData,
            this.radioChimeEndBuffer
            ].filter((b): b is AudioBuffer => !!b);
        }

        this.queue.push(...parsedAudioData);

        if (!this.isPlaying) {
            this.processAudioQueue();
        }
    }

    public async stopAudio(): Promise<void> {
        /**
         * Stops the currently playing audio and clears the queue.
         */
        if (this.currentSource) {
            this.currentSource.stop(0);
            this.currentSource = null;
        }
        this.queue = [];
        this.isPlaying = false;
    }

    public async playRadioStartChime(): Promise<void> {
        /**
         * Plays the radio start chime.
         */
        if (!this.audioContext) {
            console.warn('AudioContext is not initialized. Call initialize() first.');
            return;
        }
        if (!this.radioChimeStartBuffer) {
            console.warn('Radio chime start buffer is not loaded.');
            return;
        }

        this.queue.push(this.radioChimeStartBuffer);
        
        if (!this.isPlaying) {
            this.processAudioQueue();
        }
    }

    private async processAudioQueue(): Promise<void> {
        /**
         * Processes the audio queue and plays AudioBuffers sequentially.
         */
        if (this.queue.length === 0 || !this.audioContext) {
            this.isPlaying = false;
            return;
        }

        this.isPlaying = true;

        let source = this.audioContext.createBufferSource();
        source.buffer = this.queue.shift()!;
        source.connect(this.audioContext.destination);
        this.currentSource = source;
        source.start(0);

        source.onended = () => {
            this.currentSource = null;
            if (this.audioContext) this.processAudioQueue();
        };
    }

    private async loadAudioBuffer(audioContext: AudioContext, url: string): Promise<AudioBuffer> {
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

    private base64ToArrayBuffer(b64string: string): ArrayBuffer {
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

    private decodePCMChunk(audioContext: AudioContext,
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
}

export const audioPlayer = AudioPlayer.getInstance();
