import { Location, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import { UserData } from 'pages/Auth/Login'
import { ApiLogin, QUERY_KEY } from 'helpers'
import { useContext } from 'react'
import { LoginContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'
import { UserType } from 'common'

type LocationState = {
    from?: Location
}

export default function useLogin() {
    const { t } = useTranslation()
    const { loginContextValue } = useContext(LoginContext)
    const login = ApiLogin()

    const location: Location = useLocation()
    let redirectUrl: string = '/'

    if (location.state) {
        const { from } = location.state as LocationState
        redirectUrl = from ? from.pathname : '/'
    }
    
    const { isPending, error } = login
    const queryClient = useQueryClient()
    const user = queryClient.getQueryData<UserType>([QUERY_KEY.login])

    /*
    form validation schema
    */
    const schemaResolver = yupResolver(
        yup.object().shape({
            username: yup.string().required(t('Please enter Username')),
            password: yup.string().required(t('Please enter Password')),
        })
    )

    /*
    handle form submission
    */
    const onSubmit = (formData: UserData) => {
        login.mutate({username: formData['username'], password: formData['password']})
    }

    return {
        isPending,
        isLoggedIn: loginContextValue.isLoggedIn,
        user,
        error,
        schemaResolver,
        onSubmit,
        redirectUrl,
    }
}
