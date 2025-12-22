import type {
  EvaluationCriterionResponse,
  EvaluationCriterionPayload,
  InterviewRoundCriterionResponse,
  InterviewRoundCriterionPayload,
  EvaluationResponse,
  EvaluationPayload,
  EvaluationScorePayload,
  CompleteEvaluationPayload,
  EvaluationScoreResponse,
} from '../../interfaces/evaluation'
import { Get, Post, Patch, Delete } from '../api'

// ========== Evaluation Criteria ==========
export const EvaluationCriteriaService = {
  async getAll(): Promise<EvaluationCriterionResponse[]> {
    return await Get('/evaluation-criteria')
  },

  async getById(id: number): Promise<EvaluationCriterionResponse> {
    return await Get(`/evaluation-criteria/${id}`)
  },

  async create(payload: EvaluationCriterionPayload): Promise<EvaluationCriterionResponse> {
    return await Post('/evaluation-criteria', payload)
  },

  async update(id: number, payload: Partial<EvaluationCriterionPayload>): Promise<EvaluationCriterionResponse> {
    return await Patch(`/evaluation-criteria/${id}`, payload)
  },

  async delete(id: number): Promise<void> {
    await Delete(`/evaluation-criteria/${id}`)
  },
}

// ========== Interview Round Criteria ==========
export const InterviewRoundCriteriaService = {
  async getByRoundId(roundId: number): Promise<InterviewRoundCriterionResponse[]> {
    return await Get(`/interview-rounds/${roundId}/criteria`)
  },

  async addToRound(roundId: number, payload: InterviewRoundCriterionPayload): Promise<InterviewRoundCriterionResponse> {
    return await Post(`/interview-rounds/${roundId}/criteria`, payload)
  },

  async update(id: number, payload: Partial<InterviewRoundCriterionPayload>): Promise<InterviewRoundCriterionResponse> {
    return await Patch(`/interview-round-criteria/${id}`, payload)
  },

  async removeFromRound(id: number): Promise<void> {
    await Delete(`/interview-round-criteria/${id}`)
  },
}

// ========== Evaluations ==========
export const EvaluationService = {
  async getAll(filters?: {
    interview_round_id?: number
    application_scholarship_id?: number
    status?: string
  }): Promise<EvaluationResponse[]> {
    let url = '/evaluations'
    const params = new URLSearchParams()

    if (filters?.interview_round_id) {
      params.append('interview_round_id', filters.interview_round_id.toString())
    }
    if (filters?.application_scholarship_id) {
      params.append('application_scholarship_id', filters.application_scholarship_id.toString())
    }
    if (filters?.status) {
      params.append('status', filters.status)
    }

    const queryString = params.toString()
    if (queryString) {
      url += `?${queryString}`
    }

    return await Get(url)
  },

  async getById(id: number): Promise<EvaluationResponse> {
    return await Get(`/evaluations/${id}`)
  },

  async create(payload: EvaluationPayload): Promise<EvaluationResponse> {
    return await Post('/evaluations', payload)
  },

  async update(id: number, payload: { status?: string; remark?: string }): Promise<EvaluationResponse> {
    return await Patch(`/evaluations/${id}`, payload)
  },

  async delete(id: number): Promise<void> {
    await Delete(`/evaluations/${id}`)
  },

  async complete(id: number, payload: CompleteEvaluationPayload): Promise<EvaluationResponse> {
    return await Post(`/evaluations/${id}/complete`, payload)
  },
}

// ========== Evaluation Scores ==========
export const EvaluationScoreService = {
  async addScore(evaluationId: number, payload: EvaluationScorePayload): Promise<EvaluationScoreResponse> {
    return await Post(`/evaluations/${evaluationId}/scores`, payload)
  },

  async update(id: number, payload: { score_value?: number; comment?: string }): Promise<EvaluationScoreResponse> {
    return await Patch(`/evaluation-scores/${id}`, payload)
  },

  async delete(id: number): Promise<void> {
    await Delete(`/evaluation-scores/${id}`)
  },
}
