import { useEffect, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import { logger } from 'helpers'
import { LocationContext } from 'contexts'

const User = () => {
    const [t] = useTranslation('translation')
    const { updateLocationContextValue } = useContext(LocationContext)

    useEffect(() => {
        logger.log(`Change Location to ${t('title.user')}`)
        updateLocationContextValue({ location: 'user' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (<div>Not implemented yet</div>)
}

export default User