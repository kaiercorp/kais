import jwtDecode from 'jwt-decode'
import axios from 'axios'
import config from '../../config'
import { logger } from 'helpers/logger/logger'

// content type
axios.defaults.headers.post['Content-Type'] = 'application/json'
axios.defaults.timeout = 1000 * 60 * 3
axios.defaults.baseURL = config.baseURL

// intercepting to capture errors
axios.interceptors.response.use(
    (response) => {
        return response.data
    },
    (error) => {
        if (error && error.response && error.response.data) {
            return Promise.reject(error.response.data.message)
        } else {
            return Promise.reject(`Cant' Connect to Server`)
        }
    }
)

const AUTH_SESSION_KEY = 'kaier_user'

/**
 * Sets the default authorization
 * @param {*} token
 */
const setAuthorization = (token: string | null) => {
    if (token) axios.defaults.headers.common['Authorization'] = 'JWT ' + token
    else delete axios.defaults.headers.common['Authorization']
}

const getUserFromSession = () => {
    const user = sessionStorage.getItem(AUTH_SESSION_KEY)
    return user ? (typeof user == 'object' ? user : JSON.parse(user)) : null
}
class APICore {
    /**
     * Fetches data from given url
     */
    get = (url: string, params: any) => {
        let response
        if (params) {
            var queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : ''
            response = axios.get(`${url}?${queryString}`, {params: params, headers: {Authorization: this.getToken()}})
        } else {
            response = axios.get(`${url}`, {params: params, headers: {Authorization: this.getToken()}})
        }
        return response
    };

    getFile = (url: string, params: any) => {
        let response
        if (params) {
            var queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : ''
            response = axios.get(`${url}?${queryString}`, { responseType: 'blob', headers: {Authorization: this.getToken()}})
        } else {
            response = axios.get(`${url}`, { responseType: 'blob', headers: {Authorization: this.getToken()} })
        }
        return response
    };

    getMultiple = (urls: string, params: any) => {
        const reqs = []
        let queryString = ''
        if (params) {
            queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : ''
        }

        for (const url of urls) {
            reqs.push(axios.get(`${url}?${queryString}`, {headers: {Authorization: this.getToken()}}))
        }
        return axios.all(reqs)
    };

    /**
     * post given data to url
     */
    post = (url: string, data: any) => {
        return axios.post(url, data, {headers: {Authorization: this.getToken()}})
    };

    /**
     * Updates patch data
     */
    updatePatch = (url: string, data: any) => {
        return axios.patch(url, data, {headers: {Authorization: this.getToken()}})
    };

    /**
     * Updates data
     */
    update = (url: string, data: any) => {
        return axios.put(url, data, {headers: {Authorization: this.getToken()}})
    };

    /**
     * Deletes data
     */
    delete = (url: string) => {
        return axios.delete(url, {headers: {Authorization: this.getToken()}})
    };

    /**
     * Deletes data
     */
    deleteWithData = (url: string, config: any) => {
        return axios.delete(url, { data: config, headers: {Authorization: this.getToken()}})
    };

    /**
     * post given data to url with file
     */
    createWithFile = (url: string, data: any) => {
        const formData = new FormData()
        for (const k in data) {
            formData.append(k, data[k])
        }

        const config: any = {
            headers: {
                ...axios.defaults.headers,
                'Authorization': this.getToken(),
                'content-type': 'multipart/form-data',
            },
        }
        return axios.post(url, formData, config)
    };

    /**
     * post given data to url with file
     */
    updateWithFile = (url: string, data: any) => {
        const formData = new FormData()
        for (const k in data) {
            formData.append(k, data[k])
        }

        const config: any = {
            headers: {
                ...axios.defaults.headers,
                'Authorization': this.getToken(),
                'content-type': 'multipart/form-data',
            },
        }
        return axios.patch(url, formData, config)
    };

    isUserAuthenticated = () => {
        const user = this.getLoggedInUser()

        try {
            if (!user || (user && !user.token)) {
                return false
            }
            
            const decoded: any = jwtDecode(user.token)
            const currentTime = Date.now() / 1000
            
            if (decoded.exp < currentTime) {
                logger.warn('access token expired')
                return false
            } else {
                return true
            } 
        } catch (error) {
            logger.error(error)
        }
    };

    setLoggedInUser = (session: any) => {
        if (session) {
            sessionStorage.setItem(AUTH_SESSION_KEY, JSON.stringify(session))
        }
        else {
            sessionStorage.removeItem(AUTH_SESSION_KEY)
        }
    };
    /**
     * Returns the logged in user
     */
    getLoggedInUser = () => {
        const user = getUserFromSession()
        
        try {
            if (!user || (user && !user.token)) {
                this.setLoggedInUser(null)
                return null
            }
            
            const decoded: any = jwtDecode(user.token)
            const currentTime = Date.now() / 1000
            
            if (decoded.exp < currentTime) {
                logger.warn('access token expired')
                this.setLoggedInUser(null)
                return null
            } 
        } catch (error) {
            logger.error(error)
            this.setLoggedInUser(null)
            return null
        }
        
        return user
    };

    setUserInSession = (modifiedUser: any) => {
        let userInfo = sessionStorage.getItem(AUTH_SESSION_KEY)
        if (userInfo) {
            const { token, user } = JSON.parse(userInfo)
            this.setLoggedInUser({ token, ...user, ...modifiedUser })
        }
    };

    getToken = () => {
        const user = this.getLoggedInUser()
        if (user && user.token) {
            return user.token
        }
        return ''
    }
}

export { APICore, setAuthorization }
