import { createContext } from 'react'

type Theme = 'light' | 'dark' | 'system'

interface ThemeProviderState {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const initialState: ThemeProviderState = {
  theme: 'system',
  setTheme: () => null,
}

export const ThemeProviderContext = createContext<ThemeProviderState>(initialState)
export type { Theme, ThemeProviderState }
