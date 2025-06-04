import '@testing-library/jest-dom'
import { render } from '@testing-library/react'
import App, { sample } from '../src/App'

describe('<App /> sample', () => {
  // 렌더링 테스트
  it('<App />', () => {
    render(
        <App />
    )
  })

  // 로직 테스트
  it('sample function', () => {
    expect(sample(1, 2)).toBe(3)
  })
})
