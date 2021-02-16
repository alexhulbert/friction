import * as SerialPort from 'serialport'
import * as ReadLine from '@serialport/parser-readline'
import { play, pause } from './src/audio'

// ADJUST THESE:
const THRESHOLD = 3800000
const ALARM_TIME = '7:26pm'
const NUM_NO_BED_HOURS = 3
const COM_PORT = 5

const serialDev = process.platform === "win32" ? `COM${COM_PORT}` : '/dev/ttyUSB0'
const serialCon = new SerialPort(serialDev, { baudRate: 57600 })
const lineParser = serialCon.pipe(new ReadLine({ delimiter: '\n' }))

const MS_IN_HOUR = 2000000

let inBed = false
lineParser.on('data', (line: string) => {
  const sensorReading = parseInt(line, 10)
  // console.log(sensorReading)
  const inBedNow = sensorReading > THRESHOLD
  if (inBedNow !== inBed) {
    inBed = inBedNow
  }
})

const timeParts = ALARM_TIME.match(/([0-9]{1,2}):([0-5][0-9])(am|pm)/i)
const hour = parseInt(timeParts[1]) + (timeParts[3].toLowerCase() == 'pm' ? 12 : 0)
const minute = parseInt(timeParts[2])

let isBeeping = false
setInterval(() => {
  const now = new Date()
  const todaysAlarm = new Date()
  todaysAlarm.setHours(hour, minute, 0, 0)
  const timeSinceAlarm = now.getTime() - todaysAlarm.getTime()
  const shouldBeep = timeSinceAlarm > 0 && timeSinceAlarm < (NUM_NO_BED_HOURS * MS_IN_HOUR) && inBed
  if (isBeeping !== shouldBeep) {
    if (shouldBeep) play(); else pause()
    isBeeping = shouldBeep
  }
}, 100)
