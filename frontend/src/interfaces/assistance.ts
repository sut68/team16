import type { UserResponse } from "./user"
import type { ChatroomRespone } from "./chatroom"

export interface AssistanceRespone {
    ID: number
    massage: string
    CreatedAt: string

    sender_id: number
    sender: UserResponse

    chatroom_id: number
    chatroom: ChatroomRespone
}

export interface AssistanceCreate {
    massage: string

    sender_id: number
    chatroom_id: number

}

export interface AssistanceUpdate {
    ID: number
    massage?: string

    sender_id?: number
    chatroom_id?: number

}