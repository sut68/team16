import { Get, Post, Put, Delete } from './https';
import type { ScholarshipStatusCreate, ScholarshipStatusUpdate } from '@/interfaces';


export const StatusscholarshipAPI = {
  create: (data: ScholarshipStatusCreate) => Post("/statusscholarship", data),
  getAll: () => Get("/statusscholarship"),
  getById: (id: number) => Get(`/statusscholarship/${id}`),
  update: (id: number, data: ScholarshipStatusUpdate) => Put(`/statusscholarship/${id}`, data),
  delete: (id: number) => Delete(`/statusscholarship/${id}`),

  
};