export interface TypeFeatureResponse {
  ID: number
  type_feature_name: string
}

export interface TypeFeatureCreate {
  type_feature_name: string
}

export interface TypeFeatureUpdate {
  ID: number
  type_feature_name?: string
}

export interface FeatureScholarshipResponse {
  ID: number
  operator: string
  value: string

  scholarship_id: number

  typefeature_id: number
  typefeature: TypeFeatureResponse
}

export interface FeatureScholarshipCreate {
  operator: string;
  value: string;

  scholarship_id: number;
  typefeature_id: number;
}

export interface FeatureScholarshipUpdate {
  ID: number;
  operator?: string;
  value?: string;

  scholarship_id?: number;
  typefeature_id?: number;
}