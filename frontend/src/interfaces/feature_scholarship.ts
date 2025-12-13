import type { TypeFeatureResponse } from './type_feature'

export interface FeatureScholarshipResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt?: string | null

  feature_scholarship_name: string 
  
  operator: string
  value: string

  scholarship_id: number
  typefeature_id: number

  typefeature?: TypeFeatureResponse 
  
}