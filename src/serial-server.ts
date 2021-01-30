import * as SerialPort from 'serialport'
import * as ReadLine from '@serialport/parser-readline'

import { state, refreshAlarm } from './shared'

const serialDev = process.platform === "win32" ? 'COM5' : '/dev/ttyUSB0'
const serialCon = new SerialPort(serialDev, { baudRate: 57600 })
const lineParser = serialCon.pipe(new ReadLine({ delimiter: '\n' }))

lineParser.on('data', (line: String) => {
  if (line.includes('PRESSED')) state.inBed = true
  if (line.includes('RELEASED')) state.inBed = false
  refreshAlarm()
})
