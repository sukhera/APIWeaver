import { describe, it, expect } from 'vitest'
import { cn } from '../utils'

describe('cn utility function', () => {
  it('combines class names correctly', () => {
    const result = cn('btn', 'btn-primary', 'rounded')
    expect(result).toBe('btn btn-primary rounded')
  })

  it('handles conditional classes', () => {
    const isActive = true
    const isDisabled = false
    
    const result = cn(
      'btn',
      isActive && 'btn-active',
      isDisabled && 'btn-disabled'
    )
    
    expect(result).toBe('btn btn-active')
  })

  it('handles undefined and null values', () => {
    const result = cn('btn', undefined, null, 'btn-primary')
    expect(result).toBe('btn btn-primary')
  })

  it('merges conflicting Tailwind classes', () => {
    // tailwind-merge should handle conflicting classes
    const result = cn('bg-red-500', 'bg-blue-500')
    expect(result).toBe('bg-blue-500') // Should keep the last one
  })

  it('handles arrays of classes', () => {
    const result = cn(['btn', 'btn-primary'], 'rounded')
    expect(result).toBe('btn btn-primary rounded')
  })

  it('handles objects with conditional classes', () => {
    const result = cn({
      btn: true,
      'btn-primary': true,
      'btn-disabled': false,
    })
    
    expect(result).toBe('btn btn-primary')
  })

  it('handles empty input', () => {
    const result = cn()
    expect(result).toBe('')
  })

  it('handles complex combinations', () => {
    const variant: string = 'primary'
    const size: string = 'lg'
    const disabled = false
    
    const result = cn(
      'btn',
      {
        'btn-primary': variant === 'primary',
        'btn-secondary': variant === 'secondary',
        'btn-lg': size === 'lg',
        'btn-sm': size === 'sm',
      },
      disabled && 'btn-disabled',
      'focus:outline-none'
    )
    
    expect(result).toBe('btn btn-primary btn-lg focus:outline-none')
  })

  it('handles Tailwind responsive classes', () => {
    const result = cn(
      'text-sm',
      'md:text-base',
      'lg:text-lg',
      'text-primary'
    )
    
    expect(result).toBe('text-sm md:text-base lg:text-lg text-primary')
  })

  it('handles hover and focus states', () => {
    const result = cn(
      'bg-blue-500',
      'hover:bg-blue-600',
      'focus:bg-blue-700',
      'transition-colors'
    )
    
    expect(result).toBe('bg-blue-500 hover:bg-blue-600 focus:bg-blue-700 transition-colors')
  })
})