/**
 * Transactions & Credits API
 */

import { get, post } from '@/utils/request'
import type { 
  GetUserCreditsResponse,
  PurchaseCreditsRequest,
  PurchaseCreditsResponse,
  GetTransactionsRequest,
  GetTransactionsResponse
} from '@/types/api'

// Get user credits
export const getUserCredits = () => {
  return get<GetUserCreditsResponse>('/credits/balance')
}

// Purchase credits
export const purchaseCredits = (data: PurchaseCreditsRequest) => {
  return post<PurchaseCreditsResponse>('/credits/purchase', data)
}

// Get credit packages
export const getCreditPackages = () => {
  return get('/credits/packages')
}

// Get transaction history
export const getTransactions = (params?: GetTransactionsRequest) => {
  return get<GetTransactionsResponse>('/transactions', params)
}

// Get transaction details
export const getTransactionById = (id: string) => {
  return get(`/transactions/${id}`)
} 