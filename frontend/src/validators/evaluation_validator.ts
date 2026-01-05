// src/validators/evaluation_validator.ts
import type { EvaluationCriterionPayload, EvaluationScorePayload } from '@/interfaces/evaluation'

// TYPES & INTERFACES

// ผลลัพธ์การ validate
export interface ValidationResult {
  valid: boolean
  errors: Record<string, string | undefined>
}

// Error fields สำหรับ EvaluationCriterion form
export interface CriterionFormErrors {
  name?: string
  description?: string
  max_score?: string
  weight?: string
}

// Error fields สำหรับ EvaluationScore form
export interface ScoreFormErrors {
  score_value?: string
  comment?: string
}

// VALIDATION RULES
export const CRITERION_VALIDATION_RULES = {
  name: {
    minLength: 2,
    maxLength: 100,
    required: true,
  },
  description: {
    maxLength: 500,
    required: false,
  },
  max_score: {
    min: 0,
    max: 1000,
  },
  weight: {
    min: 0,
    max: 10,
  },
} as const

// กฎการ validate สำหรับ Evaluation Score
export const SCORE_VALIDATION_RULES = {
  score: {
    min: 0,
  },
  comment: {
    maxLength: 1000,
    required: false,
  },
} as const

// VALIDATION FUNCTIONS
export function validateCriterionForm(
  form: Partial<EvaluationCriterionPayload>
): { valid: boolean; errors: CriterionFormErrors } {
  const errors: CriterionFormErrors = {}

  // Required, 2-100 characters
  const name = form.name?.trim() || ''
  if (!name) {
    errors.name = 'กรุณากรอกชื่อเกณฑ์'
  } else if (name.length < CRITERION_VALIDATION_RULES.name.minLength) {
    errors.name = `ชื่อเกณฑ์ต้องมีอย่างน้อย ${CRITERION_VALIDATION_RULES.name.minLength} ตัวอักษร`
  } else if (name.length > CRITERION_VALIDATION_RULES.name.maxLength) {
    errors.name = `ชื่อเกณฑ์ต้องไม่เกิน ${CRITERION_VALIDATION_RULES.name.maxLength} ตัวอักษร`
  }

  // max 500 characters
  const description = form.description || ''
  if (description.length > CRITERION_VALIDATION_RULES.description.maxLength) {
    errors.description = `รายละเอียดต้องไม่เกิน ${CRITERION_VALIDATION_RULES.description.maxLength} ตัวอักษร`
  }

  // 0-1000
  const maxScore = form.max_score ?? 0
  if (maxScore < CRITERION_VALIDATION_RULES.max_score.min || maxScore > CRITERION_VALIDATION_RULES.max_score.max) {
    errors.max_score = `คะแนนเต็มต้องอยู่ระหว่าง ${CRITERION_VALIDATION_RULES.max_score.min} - ${CRITERION_VALIDATION_RULES.max_score.max}`
  }

  // 0-10
  const weight = form.weight ?? 0
  if (weight < CRITERION_VALIDATION_RULES.weight.min || weight > CRITERION_VALIDATION_RULES.weight.max) {
    errors.weight = `น้ำหนักต้องอยู่ระหว่าง ${CRITERION_VALIDATION_RULES.weight.min} - ${CRITERION_VALIDATION_RULES.weight.max}`
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

export function validateCriterionName(name: string | undefined): string | undefined {
  const trimmed = name?.trim() || ''

  if (!trimmed) {
    return 'กรุณากรอกชื่อเกณฑ์'
  }
  if (trimmed.length < CRITERION_VALIDATION_RULES.name.minLength) {
    return `ชื่อเกณฑ์ต้องมีอย่างน้อย ${CRITERION_VALIDATION_RULES.name.minLength} ตัวอักษร`
  }
  if (trimmed.length > CRITERION_VALIDATION_RULES.name.maxLength) {
    return `ชื่อเกณฑ์ต้องไม่เกิน ${CRITERION_VALIDATION_RULES.name.maxLength} ตัวอักษร`
  }

  return undefined
}

export function validateCriterionDescription(description: string | undefined): string | undefined {
  const text = description || ''

  if (text.length > CRITERION_VALIDATION_RULES.description.maxLength) {
    return `รายละเอียดต้องไม่เกิน ${CRITERION_VALIDATION_RULES.description.maxLength} ตัวอักษร`
  }

  return undefined
}

export function validateCriterionMaxScore(maxScore: number | null | undefined): string | undefined {
  if (maxScore === null || maxScore === undefined || isNaN(maxScore)) {
    return 'กรุณากรอกคะแนนเต็ม'
  }
  if (maxScore < CRITERION_VALIDATION_RULES.max_score.min || maxScore > CRITERION_VALIDATION_RULES.max_score.max) {
    return `คะแนนเต็มต้องอยู่ระหว่าง ${CRITERION_VALIDATION_RULES.max_score.min} - ${CRITERION_VALIDATION_RULES.max_score.max}`
  }

  return undefined
}

export function validateCriterionWeight(weight: number | null | undefined): string | undefined {
  if (weight === null || weight === undefined || isNaN(weight)) {
    return 'กรุณากรอกน้ำหนัก'
  }
  if (weight < CRITERION_VALIDATION_RULES.weight.min || weight > CRITERION_VALIDATION_RULES.weight.max) {
    return `น้ำหนักต้องอยู่ระหว่าง ${CRITERION_VALIDATION_RULES.weight.min} - ${CRITERION_VALIDATION_RULES.weight.max}`
  }

  return undefined
}

export function isCriterionFormValid(form: Partial<EvaluationCriterionPayload>): boolean {
  const name = form.name?.trim() || ''
  const description = form.description || ''
  const maxScore = form.max_score ?? 0
  const weight = form.weight ?? 0

  return (
    name.length >= CRITERION_VALIDATION_RULES.name.minLength &&
    name.length <= CRITERION_VALIDATION_RULES.name.maxLength &&
    description.length <= CRITERION_VALIDATION_RULES.description.maxLength &&
    maxScore >= CRITERION_VALIDATION_RULES.max_score.min &&
    maxScore <= CRITERION_VALIDATION_RULES.max_score.max &&
    weight >= CRITERION_VALIDATION_RULES.weight.min &&
    weight <= CRITERION_VALIDATION_RULES.weight.max
  )
}

// VALIDATION FUNCTIONS - EVALUATION SCORE
export function validateScoreForm(
  form: Partial<EvaluationScorePayload>,
  maxScore: number
): { valid: boolean; errors: ScoreFormErrors } {
  const errors: ScoreFormErrors = {}

  // 0 - maxScore
  const scoreValue = form.score_value ?? 0
  if (scoreValue < SCORE_VALIDATION_RULES.score.min || scoreValue > maxScore) {
    errors.score_value = `คะแนนต้องอยู่ระหว่าง ${SCORE_VALIDATION_RULES.score.min} - ${maxScore}`
  }

  // Optional, max 1000 characters
  const comment = form.comment || ''
  if (comment.length > SCORE_VALIDATION_RULES.comment.maxLength) {
    errors.comment = `หมายเหตุต้องไม่เกิน ${SCORE_VALIDATION_RULES.comment.maxLength} ตัวอักษร`
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors,
  }
}

export function validateScore(score: number | null | undefined, maxScore: number): string | undefined {
  if (score === null || score === undefined || isNaN(score)) {
    return 'กรุณากรอกคะแนน'
  }
  if (score < SCORE_VALIDATION_RULES.score.min || score > maxScore) {
    return `คะแนนต้องอยู่ระหว่าง ${SCORE_VALIDATION_RULES.score.min} - ${maxScore}`
  }

  return undefined
}

export function validateScoreComment(comment: string | undefined): string | undefined {
  const text = comment || ''

  if (text.length > SCORE_VALIDATION_RULES.comment.maxLength) {
    return `หมายเหตุต้องไม่เกิน ${SCORE_VALIDATION_RULES.comment.maxLength} ตัวอักษร`
  }

  return undefined
}


export const EvaluationValidators = {
  criterion: {
    rules: CRITERION_VALIDATION_RULES,
    validateForm: validateCriterionForm,
    validateName: validateCriterionName,
    validateDescription: validateCriterionDescription,
    validateMaxScore: validateCriterionMaxScore,
    validateWeight: validateCriterionWeight,
    isFormValid: isCriterionFormValid,
  },
  score: {
    rules: SCORE_VALIDATION_RULES,
    validateForm: validateScoreForm,
    validateScore,
    validateComment: validateScoreComment,
  },
} as const
