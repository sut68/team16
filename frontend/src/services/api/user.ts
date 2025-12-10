import type { UserResponse, CreateUserRequest } from '@/interfaces/user'
import { https } from '@/services/api'

export const listUsers = async (): Promise<UserResponse[]> => {
  const r = await https.get<UserResponse[]>('/users')
  return r.data
}

export const createUser = async (payload: CreateUserRequest & { role?: string }) => {
  const r = await https.post('/users', payload)
  return r.data
}

export const deleteUser = async (id: number) => {
  const r = await https.delete(`/users/${id}`)
  return r.data
}
