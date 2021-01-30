import { play, pause } from './audio'

const HOURS = 60 * 60 * 1000

export const state = {
  inBed: false,
  curAlarm: 0
}

export function refreshAlarm() {
  const now = new Date().getTime()
  console.log('Alarm Refresh:', { ...state, now })
  if (now - state.curAlarm < 3 * HOURS && state.inBed) {
    play()
  } else {
    pause()
  }
}
