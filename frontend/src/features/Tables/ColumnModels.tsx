import { useTranslation } from "react-i18next"
import styled from "styled-components"

import CustomSortIcon from "./ColumnFormatters/CustomSortIcon"
import { customModelsortByCorrect, customModelsortByPerf } from "helpers"
import { IconButtonWithPopover } from "components"
import { ResultOrder } from "common"

const ButtonRow = styled.div`
display: flex;

& button {
  margin-right: 3px;
  margin-left: 3px;
}
`

export const ColumnModels = ({
  sample,
  openMultiTest,
  downloadModel
}: any) => {
  const [t] = useTranslation('translation')
  let columns = new Array<any>()

  columns.push({
    dataField: 'model_id',
    text: '',
    events: {
      onClick: (e: any, column: any, columnIndex: any, row: any, rowIndex: any) => {
        e.stopPropagation()
      },
    },
    formatter: (cell: any, row: any, rowIndex: any) => {
      if (!row.is_model_saved) return
      return (
        <ButtonRow>
          <IconButtonWithPopover
            name={`btn-stop-train-${rowIndex}`}
            variant='success'
            onClick={() => downloadModel(row)}
            popTitle={t('button.save.model')}
            icon='mdi-download'
          />
          <IconButtonWithPopover
            name={`btn-opentest-${rowIndex}`}
            variant='info'
            onClick={() => openMultiTest(row)}
            popTitle={t('button.test.model')}
            icon='mdi mdi-folder'
          />
        </ButtonRow>
      )
    },
  })

  columns.push({
    dataField: 'train_local_id',
    text: 'Trial',
    headerStyle: { textAlign: 'center' },
    align: () => {
      return 'center'
    },
    sort: true,
    sortCaret: CustomSortIcon,
    formatter: (cell: any, row: any) => {
      return `Trial${cell}`
    },
  })

  columns.push({
    dataField: 'name',
    text: 'Model',
    headerStyle: { textAlign: 'center' },
    events: {
      onClick: (e: any, column: any, columnIndex: any, row: any, rowIndex: any) => {
        e.stopPropagation()
      },
    },
    formatter: (cell: any, row: any, rowIndex: any) => {
      return cell
    },
    align: () => {
      return 'center'
    },
    sort: true,
    sortCaret: CustomSortIcon,
  })
  
  columns.push({
    dataField: 'inference_time',
    text: t('ui.formatter.inftime'),
    headerStyle: { textAlign: 'center' },
    formatter: (cell: any, row: any, rowIndex: any) => {
      return (cell * 1000).toFixed(2) + 'ms'
    },
    align: () => {
      return 'center'
    },
    sort: true,
    sortCaret: CustomSortIcon,
  })

  if (sample === undefined || !sample) return columns

  if (sample.hasOwnProperty('all_result') && !sample['all_result'].hasOwnProperty('String')) {
    const result_keys = Object.keys(sample.all_result).sort(function (x, y) {
      const xidx = ResultOrder[x] || 999
      const yidx = ResultOrder[y] || 999
      return xidx - yidx
    })

    result_keys.forEach((key: string) => {
      if (key.toLowerCase() === 'model' || key.toLowerCase() === 'epoch') return
      if (key.startsWith('c_') || key.startsWith('b_')) return

      const keyTitle = t(`metric.${key}`)
      if (keyTitle === undefined) return
      columns.push({
        dataField: key,
        text: ['mse', 'mae', 'rmse'].includes(key) ? keyTitle : keyTitle + '(%)',
        headerStyle: { textAlign: 'center' },
        align: () => {
          return 'center'
        },
        sort: true,
        sortCaret: CustomSortIcon,
        sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
          customModelsortByPerf(a, b, order, dataField, rowA, rowB),
        formatter: (cell: any, row: any, rowIndex: any) => {
          if (['mse', 'mae', 'rmse'].includes(key)) {
            return (Number(row['all_result'][key])).toFixed(3)
          }
          return (Number(row['all_result'][key]) * 100).toFixed(2)
        }
      })
    })
  }

  if (sample.hasOwnProperty('class_result') && !sample['class_result'].hasOwnProperty('String')) {
    const class_keys = Object.keys(sample.class_result)
    class_keys.forEach((key: string) => {
      const keyTitleAcc = key + ' ' + t(`ui.result.car.acc`)
      columns.push({
        dataField: `${key}_acc`,
        text: keyTitleAcc,
        headerStyle: { textAlign: 'center' },
        align: () => {
          return 'center'
        },
        sort: true,
        sortCaret: CustomSortIcon,
        sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
          customModelsortByCorrect(a, b, order, dataField, rowA, rowB),
      })

      const keyTitleCorrect = key + ' ' + t(`ui.result.car.n_correct_over_total`)
      columns.push({
        dataField: `${key}_correct`,
        text: keyTitleCorrect,
        headerStyle: { textAlign: 'center' },
        align: () => {
          return 'center'
        },
        sort: true,
        sortCaret: CustomSortIcon,
        sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
          customModelsortByCorrect(a, b, order, dataField, rowA, rowB),
      })
    })
  }
  
  if (sample.hasOwnProperty('accuracy_info') && !sample.accuracy_info.hasOwnProperty('String')) {
    const accuracyInfoKeys = Object.keys(sample.accuracy_info)
    accuracyInfoKeys.forEach((key: string) => {
      const keyTitleCorrect = key + ' ' + t(`ui.result.car.correct_${key}_over_total`)
      columns.push({
        dataField: `correct_${key}_over_total`,
        text: keyTitleCorrect,
        headerStyle: { textAlign: 'center' },
        align: () => {
          return 'center'
        },
        sort: true,
        sortCaret: CustomSortIcon,
        sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
          customModelsortByCorrect(a, b, order, dataField, rowA, rowB),
      })
    })
  }

  return columns
}
