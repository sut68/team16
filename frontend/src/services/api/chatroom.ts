import { Get, Post, Put, Delete } from './https';
import type { ChatroomUpdate, ChatroomRespone } from '@/interfaces';

export const ChatroomAPI = {
  create: () => Post("/chatroom", {}),
  getMyOpen: () => Get("/chatroom/my-open"),
  getAllOpen: (): Promise<{ data: ChatroomRespone[] }> => Get("/chatroom/open"),
  getAll: (): Promise<{ data: ChatroomRespone[] }> => Get("/chatroom"),
  getById: (id: number) => Get(`/chatroom/${id}`),
  update: (id: number, data: ChatroomUpdate) => Put(`/chatroom/${id}`, data),
  delete: (id: number) => Delete(`/chatroom/${id}`),
};
