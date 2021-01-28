import Fastify from 'fastify'
import * as SerialPort from 'serialport'
import * as ReadLine from '@serialport/parser-readline'
import { play, pause } from './audio'

const HOURS = 60 * 60 * 1000
const SERIAL_LOC = '/dev/ttyUSB0' // 'COM5'
const HTTP_PORT = 3000

const http = Fastify()
const port = new SerialPort(SERIAL_LOC, { baudRate: 57600 })
const lineParser = port.pipe(new ReadLine({ delimiter: '\n' }))

let inBed = false
let lastAlarm = 0

function refreshAlarm() {
  const now = new Date().getTime()
  console.log({ now, lastAlarm, inBed })
  if (now - lastAlarm < 3 * HOURS && inBed) {
    play()
  } else {
    pause()
  }
}

lineParser.on('data', (line: String) => {
  if (line.includes('PRESSED')) inBed = true
  if (line.includes('RELEASED')) inBed = false
  refreshAlarm()
})

http.post('/', async () => {
  lastAlarm = new Date().getTime()
  refreshAlarm()
  return 'OK'
})

http.listen(HTTP_PORT, '0.0.0.0')
