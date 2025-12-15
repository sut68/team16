import { Get, Post, Put, Delete } from './https';
import type { ScholarshipTypeCreate, ScholarshipTypeUpdate } from '@/interfaces';


export const TypescholarshipAPI = {
  create: (data: ScholarshipTypeCreate) => Post("/typescholarship", data),
  getAll: () => Get("/typescholarship"),
  getById: (id: number) => Get(`/typescholarship/${id}`),
  update: (id: number, data: ScholarshipTypeUpdate) => Put(`/typescholarship/${id}`, data),
  delete: (id: number) => Delete(`/typescholarship/${id}`),
};