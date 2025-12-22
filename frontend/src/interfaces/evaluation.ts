// Response from backend
// ----------------------
export interface EvaluationCriterionResponse {
  ID: number
  name: string
  description: string
  score_type: 'numeric' | 'grade' | 'pass_fail'
  max_score: number
  weight: number
  is_active: boolean
  CreatedAt?: string
  UpdatedAt?: string
}

export interface InterviewRoundCriterionResponse {
  ID: number
  interview_round_id: number
  evaluation_criterion_id: number
  weight: number
  is_enabled: boolean
  evaluation_criterion?: EvaluationCriterionResponse
}

export interface EvaluationScoreResponse {
  ID: number
  evaluation_id: number
  evaluation_criterion_id: number
  score_value: number
  comment: string
  scoring_admin_id: number
  evaluation_criterion?: EvaluationCriterionResponse
  scoring_admin?: any
}

export interface EvaluationResponse {
  ID: number
  interview_round_id: number
  application_scholarship_id: number
  admin_id: number
  total_score: number
  status: 'pending' | 'in_progress' | 'completed' | 'approved' | 'rejected'
  remark: string
  interview_round?: any
  application_scholarship?: any
  admin_profile?: any
  evaluation_scores?: EvaluationScoreResponse[]
  CreatedAt?: string
  UpdatedAt?: string
}
// ----------------------


// Payload send to Backend
// ----------------------
export interface EvaluationCriterionPayload {
  name: string
  description?: string
  score_type?: 'numeric' | 'grade' | 'pass_fail'
  max_score?: number
  weight?: number
  is_active?: boolean
}

export interface InterviewRoundCriterionPayload {
  evaluation_criterion_id: number
  weight?: number
  is_enabled?: boolean
}

export interface EvaluationPayload {
  interview_round_id: number
  application_scholarship_id: number
  admin_id: number
  remark?: string
}

export interface EvaluationScorePayload {
  evaluation_criterion_id: number
  score_value: number
  comment?: string
  scoring_admin_id: number
}

export interface CompleteEvaluationPayload {
  final_decision: 'approved' | 'rejected' | 'waitlist'
}
// ----------------------


// For UI
// ----------------------
export interface EvaluationCriterionView {
  ID: number
  name: string
  description: string
  score_type: string
  max_score: number
  weight: number
  is_active: boolean
}

export interface EvaluationView {
  ID: number
  student_name: string
  scholarship_name: string
  interview_round_name: string
  total_score: number
  status: string
  admin_name: string
  created_at: string
}
// ----------------------
