/**
 * Defines the audio chunk message structure for the audio message from the backend.
 */

export interface AudioChunkMessage {
    sample_rate: number[];
    channels: number[];
    audio_data: string[];   // base64 encoded audio data
}