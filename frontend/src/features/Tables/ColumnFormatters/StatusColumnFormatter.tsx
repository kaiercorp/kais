import { PopoverLabel, StatusLabel, StatusRow } from "components"
import CustomSortIcon from './CustomSortIcon'


const StatusColumnFormatter = (column: any, colIndex: any, sort: any, t: any) => {
    return (
      <span>
        {column.text} {CustomSortIcon(sort.dataField !== 'state' ? '' : sort.order, column)}
        <PopoverLabel name='state'>
          <StatusRow>
            <StatusLabel state={'train'} />
            <span>{t('state.train.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'additional_train'} />
            <span>{t('state.additional_train.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'finish'} />
            <span>{t('state.finish.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'finish-fail'} />
            <span>{t('state.finish_fail.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'cancel'} />
            <span>{t('state.cancel.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'fail'} />
            <span>{t('state.fail.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'test'} />
            <span>{t('state.test.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'finish_test'} />
            <span>{t('state.finish_test.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'test_cancel'} />
            <span>{t('state.test_cancel.desc')}</span>
          </StatusRow>
          <StatusRow>
            <StatusLabel state={'test_fail'} />
            <span>{t('state.test_fail.desc')}</span>
          </StatusRow>
        </PopoverLabel>
      </span>
    )
  }

export default StatusColumnFormatter