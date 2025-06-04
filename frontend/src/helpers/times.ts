export const convertUtcTime = function (utc: Date | undefined | string) {
  if (utc === undefined)
    return ''

  const date = new Date(utc)
  // const offset = new Date().getTimezoneOffset()
  // date.setTime(date.getTime() + (offset * 60 * 1000))
  date.setTime(date.getTime())
  var hour = ('00' + date.getHours()).slice(-2)
  var min = ('00' + date.getMinutes()).slice(-2)
  var sec = ('00' + date.getSeconds()).slice(-2)

  var month = ('00' + (date.getMonth() + 1)).slice(-2)
  var day = ('00' + date.getDate()).slice(-2)

  return date.getFullYear() + '/' + month + '/' + day + ' ' + hour + ':' + min + ':' + sec
}

export const convertUtcTimeToLocalTimeMonthly = function (utc: number) {
  const date = new Date(utc)
  const offset = new Date().getTimezoneOffset()
  date.setTime(date.getTime() - (offset * 60))
  var hour = ('00' + date.getHours()).slice(-2)
  var min = ('00' + date.getMinutes()).slice(-2)
  var sec = ('00' + date.getSeconds()).slice(-2)

  return (date.getMonth() + 1) + '/' + date.getDate() + ' ' + hour + ':' + min + ':' + sec
}

export const convertNowTimeForFilename = function () {
  const date = new Date()
  const offset = new Date().getTimezoneOffset()
  date.setTime(date.getTime() - (offset * 60))

  var month = ('00' + (date.getMonth() + 1)).slice(-2)
  var day = ('00' + date.getDate()).slice(-2)
  var hour = ('00' + date.getHours()).slice(-2)
  var min = ('00' + date.getMinutes()).slice(-2)
  var sec = ('00' + date.getSeconds()).slice(-2)

  return date.getFullYear() + '.' + month + day + '.' + hour + min + sec
}

export const getDuration = function (start: string, end: string) {
  const minute = 60
  const hour = minute * 60
  const day = hour * 24

  var time_sec = 0
  if (end === 'now') {
    time_sec = Math.floor((Date.now() - Date.parse(start)) / 1000)
  } else {
    time_sec = Math.floor((Date.parse(end) - Date.parse(start)) / 1000)
  }

  if (time_sec >= day) {
    const dd = Math.floor(time_sec / day)
    const hh = (Math.floor(time_sec / day)) % hour
    return `${dd}d ${hh}h`
  } else if (time_sec >= hour) {
    const hh = Math.floor(time_sec / hour)
    const mm = (Math.floor(time_sec / hour)) % minute
    return `${hh}h ${mm}m`
  } else if (time_sec >= minute) {
    const mm = Math.floor(time_sec / minute)
    const ss = Math.floor(time_sec % minute)
    return `${mm}m ${ss}s`
  } 
  
  return `${time_sec}s`
}

export const getDurationRealtime = function (start: string, end: string, state: string) {
  if (['train', 'additional_train', 'test'].includes(state)) {
    return getDuration(start, "now")
  }

  return getDuration(start, end)
}

export const convertDurationToSecond = (duration: string) => {
  const minute = 60
  const hour = minute * 60
  const day = hour * 24

  return duration.split('').reduce((acc, c) => {
    if (['d', 'h', 'm', 's'].includes(c)) {
      switch (c) {
        case 'd':
          return { time: '', second: acc.second + Number.parseInt(acc.time) * day }
        case 'h':
          return { time: '', second: acc.second + Number.parseInt(acc.time) * hour }
        case 'm':
          return { time: '', second: acc.second + Number.parseInt(acc.time) * minute }
        case 's':
          return { time: '', second: acc.second + Number.parseInt(acc.time) }
        default:
          return acc
      }
    } else {
      return { time: acc.time + c, second: acc.second }
    }
  }, { time: '', second: 0 }).second
}