import { describe, it, expect } from 'vitest'
import { renderWithProviders, screen } from '@/test/utils'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../card'

describe('Card Components', () => {
  describe('Card', () => {
    it('renders with default props', () => {
      renderWithProviders(<Card>Card content</Card>)
      
      const card = screen.getByText('Card content')
      expect(card).toBeInTheDocument()
      expect(card).toHaveClass(
        'rounded-lg',
        'border',
        'bg-card',
        'text-card-foreground',
        'shadow-sm'
      )
    })

    it('accepts custom className', () => {
      renderWithProviders(<Card className="custom-card-class">Content</Card>)
      
      const card = screen.getByText('Content')
      expect(card).toHaveClass('custom-card-class')
    })

    it('forwards ref correctly', () => {
      let cardRef: HTMLDivElement | null = null
      
      renderWithProviders(
        <Card ref={(ref) => { cardRef = ref }}>
          Content
        </Card>
      )
      
      expect(cardRef).toBeInstanceOf(HTMLDivElement)
    })
  })

  describe('CardHeader', () => {
    it('renders with correct styling', () => {
      renderWithProviders(<CardHeader>Header content</CardHeader>)
      
      const header = screen.getByText('Header content')
      expect(header).toBeInTheDocument()
      expect(header).toHaveClass('flex', 'flex-col', 'space-y-1.5', 'p-6')
    })
  })

  describe('CardTitle', () => {
    it('renders as h3 with correct styling', () => {
      renderWithProviders(<CardTitle>Card Title</CardTitle>)
      
      const title = screen.getByRole('heading', { level: 3 })
      expect(title).toBeInTheDocument()
      expect(title).toHaveTextContent('Card Title')
      expect(title).toHaveClass(
        'text-2xl',
        'font-semibold',
        'leading-none',
        'tracking-tight'
      )
    })
  })

  describe('CardDescription', () => {
    it('renders with correct styling', () => {
      renderWithProviders(<CardDescription>Card description</CardDescription>)
      
      const description = screen.getByText('Card description')
      expect(description).toBeInTheDocument()
      expect(description).toHaveClass('text-sm', 'text-muted-foreground')
    })
  })

  describe('CardContent', () => {
    it('renders with correct styling', () => {
      renderWithProviders(<CardContent>Card content</CardContent>)
      
      const content = screen.getByText('Card content')
      expect(content).toBeInTheDocument()
      expect(content).toHaveClass('p-6', 'pt-0')
    })
  })

  describe('CardFooter', () => {
    it('renders with correct styling', () => {
      renderWithProviders(<CardFooter>Footer content</CardFooter>)
      
      const footer = screen.getByText('Footer content')
      expect(footer).toBeInTheDocument()
      expect(footer).toHaveClass('flex', 'items-center', 'p-6', 'pt-0')
    })
  })

  describe('Complete Card Structure', () => {
    it('renders all components together correctly', () => {
      renderWithProviders(
        <Card>
          <CardHeader>
            <CardTitle>Test Card</CardTitle>
            <CardDescription>This is a test card description</CardDescription>
          </CardHeader>
          <CardContent>
            <p>Main card content goes here</p>
          </CardContent>
          <CardFooter>
            <button>Action Button</button>
          </CardFooter>
        </Card>
      )
      
      // Check all parts are rendered
      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Test Card')
      expect(screen.getByText('This is a test card description')).toBeInTheDocument()
      expect(screen.getByText('Main card content goes here')).toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'Action Button' })).toBeInTheDocument()
    })
  })
})