import { Get, Post, Put, Delete } from './https';
import type { FeatureScholarshipCreate, FeatureScholarshipUpdate } from '@/interfaces';


export const FeatureScholarshipAPI = {
  create: (data: FeatureScholarshipCreate) => Post("/featurescholarship", data),
  getAll: () => Get("/featurescholarship"),
  getById: (id: number) => Get(`/featurescholarship/${id}`),
  update: (id: number, data: FeatureScholarshipUpdate) => Put(`/featurescholarship/${id}`, data),
  delete: (id: number) => Delete(`/featurescholarship/${id}`),
};