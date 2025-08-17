import { UseMutationResult, UseQueryResult } from '@tanstack/react-query'
import type { ValidateResponse, ValidateRequest, AmendResponse, AmendRequest, DiffResponse, DiffRequest } from '@/types/api'

// Type assertion helper for test mocks
// This avoids using 'any' while still allowing incomplete mocks for testing
export const asMockMutation = <TData, TError = Error, TVariables = unknown>(
  mock: Record<string, unknown>
): UseMutationResult<TData, TError, TVariables, unknown> => 
  mock as unknown as UseMutationResult<TData, TError, TVariables, unknown>

export const asMockQuery = <TData, TError = Error>(
  mock: Record<string, unknown>
): UseQueryResult<TData, TError> => 
  mock as unknown as UseQueryResult<TData, TError>

// Specific type assertion helpers for our API
export const asMockValidateMutation = (mock: Record<string, unknown>) =>
  asMockMutation<ValidateResponse, Error, ValidateRequest>(mock)

export const asMockAmendMutation = (mock: Record<string, unknown>) =>
  asMockMutation<AmendResponse, Error, AmendRequest>(mock)

export const asMockDiffMutation = (mock: Record<string, unknown>) =>
  asMockMutation<DiffResponse, Error, DiffRequest>(mock)

export const asMockRealtimeValidation = (mock: Record<string, unknown>) =>
  asMockQuery<ValidateResponse, Error>(mock)
