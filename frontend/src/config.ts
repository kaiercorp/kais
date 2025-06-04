const config = {
    ENV: process.env.REACT_APP_APP_ENV,
    baseURL: (process.env.REACT_APP_API_URL || '') + process.env.REACT_APP_API_PORT + process.env.REACT_APP_API_ROOT,
    wsURL: process.env.REACT_APP_WS_ROOT,
    staticURL: (process.env.REACT_APP_API_URL || '') + (process.env.REACT_APP_API_PORT),
    basePort: process.env.REACT_APP_API_PORT
}

export default config
