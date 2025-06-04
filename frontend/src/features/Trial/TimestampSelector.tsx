import { LabelInput } from "components"
import { useTranslation } from "react-i18next"


const TimestampSelector = ({ value, onChange }: any) => {
    const [t] = useTranslation('translation')

    return (
        <LabelInput
            title={t('ui.train.dateformat')}
            name='date_format'
            value={value || ''}
            onChange={(e: any) => onChange('date_format', e.target.value)} 
            errors={undefined} 
        />
    )
}

export default TimestampSelector