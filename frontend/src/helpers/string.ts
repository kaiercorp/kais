export const trimAndRemoveDoubleSpace = (source: string) => {
  let result = source
  result = result.replace(/\s\s/gi, ' ')
  result = result.trim()
  return result
}

export const checkOnlyNumber = (v: string | number) => {
  let number = v.toString()

  number = number.replace(/[^0-9]/g,'')

  return Number(number)
}
