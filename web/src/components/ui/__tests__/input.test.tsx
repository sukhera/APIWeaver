import { describe, it, expect, vi } from 'vitest'
import { renderWithProviders, screen, userEvent } from '@/test/utils'
import { Input } from '../input'

describe('Input', () => {
  it('renders with default props', () => {
    renderWithProviders(<Input />)
    
    const input = screen.getByRole('textbox')
    expect(input).toBeInTheDocument()
    expect(input).toHaveClass(
      'flex',
      'h-10',
      'w-full',
      'rounded-md',
      'border',
      'border-input',
      'bg-background'
    )
  })

  it('handles value and onChange', async () => {
    const user = userEvent.setup()
    const handleChange = vi.fn()
    
    renderWithProviders(<Input value="" onChange={handleChange} />)
    
    const input = screen.getByRole('textbox')
    await user.type(input, 'test input')
    
    expect(handleChange).toHaveBeenCalled()
  })

  it('accepts different input types', () => {
    const { rerender } = renderWithProviders(<Input type="email" />)
    
    let input = screen.getByRole('textbox')
    expect(input).toHaveAttribute('type', 'email')

    rerender(<Input type="password" />)
    input = screen.getByDisplayValue('')
    expect(input).toHaveAttribute('type', 'password')

    rerender(<Input type="number" />)
    input = screen.getByRole('spinbutton')
    expect(input).toHaveAttribute('type', 'number')
  })

  it('handles placeholder text', () => {
    renderWithProviders(<Input placeholder="Enter text here" />)
    
    const input = screen.getByPlaceholderText('Enter text here')
    expect(input).toBeInTheDocument()
  })

  it('can be disabled', () => {
    renderWithProviders(<Input disabled />)
    
    const input = screen.getByRole('textbox')
    expect(input).toBeDisabled()
    expect(input).toHaveClass('disabled:cursor-not-allowed', 'disabled:opacity-50')
  })

  it('accepts custom className', () => {
    renderWithProviders(<Input className="custom-input-class" />)
    
    const input = screen.getByRole('textbox')
    expect(input).toHaveClass('custom-input-class')
  })

  it('forwards ref correctly', () => {
    let inputRef: HTMLInputElement | null = null
    
    renderWithProviders(
      <Input ref={(ref) => { inputRef = ref }} />
    )
    
    expect(inputRef).toBeInstanceOf(HTMLInputElement)
  })

  it('handles focus and blur events', async () => {
    const user = userEvent.setup()
    const handleFocus = vi.fn()
    const handleBlur = vi.fn()
    
    renderWithProviders(<Input onFocus={handleFocus} onBlur={handleBlur} />)
    
    const input = screen.getByRole('textbox')
    
    await user.click(input)
    expect(handleFocus).toHaveBeenCalledTimes(1)
    
    await user.tab()
    expect(handleBlur).toHaveBeenCalledTimes(1)
  })

  it('supports controlled input', async () => {
    const user = userEvent.setup()
    let value = ''
    const handleChange = vi.fn((e) => {
      value = e.target.value
    })
    
    const { rerender } = renderWithProviders(
      <Input value={value} onChange={handleChange} />
    )
    
    const input = screen.getByRole('textbox')
    expect(input).toHaveValue('')
    
    await user.type(input, 'a')
    expect(handleChange).toHaveBeenCalled()
    
    // Simulate controlled update
    value = 'ab'
    rerender(<Input value={value} onChange={handleChange} />)
    expect(input).toHaveValue('ab')
  })
})