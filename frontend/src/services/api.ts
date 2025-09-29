import axios from 'axios'


const api = axios.create({
        baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081',
        timeout: 5000,
})


export interface Todo {
        id: number | null 
        title: string
        description?: string | null
        completed: boolean
        user_id?: number
        created_at?: string
        updated_at?: string
}


export interface Meta {
        page: number
        page_size: number
        total: number
        total_pages: number
}


export async function listTodos(params: { page?: number; page_size?: number; completed?: boolean | string; sort?: string } = {}) {
        const res = await api.get('/api/todos', { params })
        return res.data.data as { items: Todo[]; meta: Meta }
}


export async function getTodo(id: number) {
        const res = await api.get(`/api/todos/${id}`)
        return res.data.data as Todo
}


export async function createTodo(payload: { title: string; description?: string; completed?: boolean }) {
        const res = await api.post('/api/todos', payload)
        return res.data.data as Todo
}


export async function updateTodo(id: number, payload: Partial<{ title: string; description: string | null; completed: boolean }>) {
        const res = await api.patch(`/api/todos/${id}`, payload)
        return res.data.data as Todo
}


export async function deleteTodo(id: number) {
        const res = await api.delete(`/api/todos/${id}`)
        return res.data.data
}