import { Button, Alert } from 'react-bootstrap'
import { Location, useLocation } from 'react-router-dom'
import { Navigate } from 'react-router-dom'
import { FormInput, VerticalForm } from 'components'
import AccountLayout from './AccountLayout'
import useLogin from 'hooks/useLogin'
import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { ApiLogout } from 'helpers/api'
import { APICore } from 'helpers/api/apiCore'

export type UserData = {
    username: string
    password: string
};

const api = new APICore()

const Login = () => {
    const { isPending, isLoggedIn, error, schemaResolver, onSubmit, redirectUrl } = useLogin()

    const logout = ApiLogout()
    const user = api.getLoggedInUser()

    useEffect(() => {
        if (user) {
            logout.mutate()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])


    const location: Location = useLocation()

    return (
        <>
            {(location.pathname !== "/auth/logout" && (isLoggedIn || user)) && <Navigate to={redirectUrl} replace />}

            <AccountLayout>

                {error && (
                    <Alert variant="danger" className="my-2">
                        {error.toString()}
                    </Alert>
                )}

                <VerticalForm<UserData>
                    onSubmit={onSubmit}
                    resolver={schemaResolver}
                >
                    <FormInput
                        label={'Username'}
                        type="text"
                        name="username"
                        placeholder={'Enter your Username'}
                        containerClass={'mb-3'}
                    />
                    <FormInput
                        label={'Password'}
                        type="password"
                        name="password"
                        placeholder={'Enter your password'}
                        containerClass={'mb-3'}
                    />

                    <div className="mb-1 mb-0 text-center">
                        <Link to='/' style={{ marginRight: '20px' }}>Home</Link>
                        <Button variant="primary" type="submit" disabled={isPending}>
                            {'Log In'}
                        </Button>
                    </div>
                </VerticalForm>
            </AccountLayout>
        </>
    )
}

export default Login
