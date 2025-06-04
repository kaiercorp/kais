import { Form } from "react-bootstrap"
import { useTranslation } from 'react-i18next'

const TableClsTrainInfo = ({ config }: any) => {
    const [t] = useTranslation('translation')
    return (
        <table>
            <tbody>
                <tr>
                    <th>{t('ui.train.indexCol')}</th>
                    <td>{config.train_config.index_column}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.labelCol')}</th>
                    <td>{config.train_config.label_column}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.includeCol')}</th>
                    <td>
                        <Form.Select
                            name='input-columns'
                            id='input-columns'
                            multiple
                            style={{ 'height': '106px' }}
                        >
                            {config.train_config.input_column && config.train_config.input_column.map((input: any) => {
                                return <option key={`input-column-${input}`}>{input}</option>
                            })}
                        </Form.Select>
                    </td>
                </tr>
            </tbody>
        </table>
    )
}

export default TableClsTrainInfo