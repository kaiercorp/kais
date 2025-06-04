import { PopoverLabel } from "components"
import CustomSortIcon from "./CustomSortIcon"


const InferenceTimeFormatter = (column: any, colIndex: any, sort: any, t: any) => {
    return (
      <span>
        {column.text} {CustomSortIcon(sort.dataField !== 'inference_time' ? '' : sort.order, column)}
        <PopoverLabel name='inference_time' variant='default'>
          {t('ui.formatter.inftime.desc')}
        </PopoverLabel>
      </span>
    )
  }

export default InferenceTimeFormatter