import { Get, Post, Put, Delete } from './https';
import type { TypeFeatureCreate, TypeFeatureUpdate } from '@/interfaces';


export const TypefeatureAPI = {
  create: (data: TypeFeatureCreate) => Post("/typefeature", data),
  getAll: () => Get("/typefeature"),
  getById: (id: number) => Get(`/typefeature/${id}`),
  update: (id: number, data: TypeFeatureUpdate) => Put(`/typefeature/${id}`, data),
  delete: (id: number) => Delete(`/typefeature/${id}`),
};