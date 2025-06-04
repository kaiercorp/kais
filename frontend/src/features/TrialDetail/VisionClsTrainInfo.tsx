import { useTranslation } from 'react-i18next'

const VisionClsTrainInfo = ({ config }: any) => {
    const [t] = useTranslation('translation')
    const width = config.train_config.width === -1 ? t('ui.train.image.resolution.origin') : config.train_config.width
    const height = config.train_config.height === -1 ? t('ui.train.image.resolution.origin') : config.train_config.height
    return (
        <table>
            <tbody>
                <tr>
                    <th colSpan={2}>{t('ui.train.config')}</th>
                </tr>
                <tr>
                    <th>{t('ui.train.image.title')}</th>
                    <td>
                        {width} x {height}
                    </td>
                </tr>
                <tr>
                    <th>{t('ui.train.base_lr')}</th>
                    <td>{(config.train_config.base_lr || 0).toFixed(6)}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.epochs')}</th>
                    <td>
                        {config.train_config.epochs < 1 ? (
                            <>{t('ui.train.max')}</>
                        ) : (
                            config.train_config.epochs
                        )}
                    </td>
                </tr>
                <tr>
                    <th>{t('ui.train.batch_size')}</th>
                    <td>
                        {config.train_config.train_batch_size < 1 ? (
                            <>{t('ui.train.max')}</>
                        ) : (
                            config.train_config.train_batch_size
                        )}
                    </td>
                </tr>
            </tbody>
        </table>
    )
}

export default VisionClsTrainInfo