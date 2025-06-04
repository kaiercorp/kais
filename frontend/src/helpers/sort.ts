import { convertDurationToSecond } from './times'

type TCustomOrder = {
  [key: string]: number
}

const customOrder: TCustomOrder = {
  train: 0,
  additional_train: 0,
  test: 1,
  finish: 2,
  'finish-fail': 3,
  finish_test: 4,
  cancel: 5,
  test_cancel: 6,
  fail: 7,
  test_fail: 8,
}

export const customTrialSort = (oldTrials: any) => {
  let newTrials = oldTrials.slice()
  newTrials = newTrials.map((trial: any) => {
    return {
      ...trial,
      customOrder: customOrder??[trial.state],
    }
  })

  newTrials = newTrials.sort(function (a: any, b: any) {
    // 상태가 같거나, 두 row가 모두 train/test 가 아니면 최신 순
    if (a.state === b.state || (customOrder[a.state] > 1 && customOrder[b.state] > 1)) {
      return a.created_at > b.created_at ? -1 : a.created_at < b.created_at ? 1 : 0
    }

    // train/test 를 우선순위
    return customOrder[a.state] - customOrder[b.state]
  })

  return newTrials
}

export const customSortFixedStatus = (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) => {
  const direction = order === 'asc' ? 1 : -1

  if (
    (customOrder[rowA.state] > 1 && customOrder[rowB.state] > 1) ||
    (customOrder[rowA.state] < 2 && customOrder[rowB.state] < 2) ||
    rowA.state === rowB.state
  ) {
    if (dataField === 'duration') {
      return convertDurationToSecond(a).toString().localeCompare(convertDurationToSecond(b).toString(), undefined, { numeric: true, sensitivity: 'base' }) * direction     
    } else if (dataField === 'inference_time') {
      return (a * 1000).toFixed(0).toString().localeCompare((b * 1000).toFixed(0).toString(), undefined, { numeric: true, sensitivity: 'base' }) * direction
    } else {
      return a.toString().localeCompare(b.toString(), undefined, { numeric: true, sensitivity: 'base' }) * direction
    }
  }

  return customOrder[rowA.state] - customOrder[rowB.state]
}

export const customModelsort = (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) => {
  let direction = order === 'asc' ? 1 : -1
  return direction * (Number(a) - Number(b))
}

export const customModelsortByPerf = (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) => {
  let direction = order === 'asc' ? -1 : 1
  if (['mse', 'mae', 'rmse'].includes(dataField)) {
    direction = direction * -1
  }
  return direction * (Number(rowA['all_result'][dataField]) - Number(rowB['all_result'][dataField]))
}

export const customModelsortByCorrect = (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) => {
  let direction = order === 'asc' ? 1 : -1
  return direction * ((a > b)?-1:1)
}

export const customSortFixedStatusAcc = (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) => {
  const direction = order === 'asc' ? 1 : -1

  if (
    (customOrder[rowA.state] > 1 && customOrder[rowB.state] > 1) ||
    (customOrder[rowA.state] < 2 && customOrder[rowB.state] < 2) ||
    rowA.state === rowB.state
  ) {
    const accA = rowA.perf.wa || 0
    const accB = rowB.perf.wa || 0
    return (accA - accB) * direction
  }

  return customOrder[rowA.state] - customOrder[rowB.state]
}