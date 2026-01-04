import { Get, Post, Put, Delete } from './https';
import type { AssistanceCreate, AssistanceUpdate } from '@/interfaces';


export const AssistanceAPI = {
  create: (data: AssistanceCreate) => Post("/assistance", data),
  getAll: () => Get("/assistance"),
  getById: (id: number) => Get(`/assistance/${id}`),
  update: (id: number, data: AssistanceUpdate) => Put(`/assistance/${id}`, data),
  delete: (id: number) => Delete(`/assistance/${id}`),
};

// import { Get, Post, Put, Delete } from './https';
// import type { AssistanceCreate, AssistanceUpdate, AssistanceRespone } from '@/interfaces';

// export const AssistanceAPI = {
//   create: (data: AssistanceCreate) => Post("/assistance", data),
//   getAll: (): Promise<{ data: AssistanceRespone[] }> => Get("/assistance"),
//   getById: (id: number) => Get(`/assistance/${id}`),
//   update: (id: number, data: AssistanceUpdate) => Put(`/assistance/${id}`, data),
//   delete: (id: number) => Delete(`/assistance/${id}`),
// };
