import type { UserResponse } from "./user"

export interface ChatroomRespone {
    ID: number

    user_id: number
    user: UserResponse

    status_chatroom: string
}

export interface ChatroomCreate{
    massage?: string

    user_id: number

    status_chatroom: string

}

export interface ChatroomUpdate{
    ID?: number

    user_id?: number

    status_chatroom?: string

}