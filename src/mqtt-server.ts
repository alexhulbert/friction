
import * as Aedes from 'aedes'
import { createServer } from 'aedes-server-factory'

import { state, refreshAlarm } from './shared'

type MQTTSubscriptionFn = (value1?: any, value2?: any) => Promise<void>

const mqttHandlers: Record<string, MQTTSubscriptionFn> = {
  alarm_alert_start: async () => {
    state.curAlarm = new Date().getTime()
  },
  alarm_snooze_clicked: async () => {
    state.curAlarm = 0
  }
}

const mqttServer = Aedes()
mqttServer.subscribe('SleepAsAndroid', ({ payload }, cb) => {
  const message = JSON.parse(payload.toString())
  if (mqttHandlers.hasOwnProperty(message.event)) {
    mqttHandlers[message.event]()
    refreshAlarm()
  } else {
    console.log('MQTT Message:', message)
  }
  cb()
}, () => {})

const tcpServer = createServer(mqttServer)
tcpServer.listen(1883, '0.0.0.0', () => console.log('DONE'))
