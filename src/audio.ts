import { AudioIO, IoStreamWrite, SampleFormat16Bit } from 'naudiodon'
import { createReadStream, ReadStream } from 'fs'

let isPlaying = false
let alarmStream: ReadStream
let audioDevice: IoStreamWrite

function reset() {
  if (alarmStream) alarmStream.close()
  if (audioDevice) audioDevice.quit()
}

function loopOnce() {
  reset()
  audioDevice = AudioIO({
    outOptions: {
      channelCount: 2,
      sampleFormat: SampleFormat16Bit,
      sampleRate: 48000,
      closeOnError: false
    }
  })
  alarmStream = createReadStream('alarm.raw')
  alarmStream.on('end', loopOnce)
  alarmStream.pipe(audioDevice)
  audioDevice.start()
}

export function play() {
  if (isPlaying) return
  isPlaying = true
  loopOnce()
}

export function pause() {
  if (!isPlaying) return
  isPlaying = false
  reset()
}
