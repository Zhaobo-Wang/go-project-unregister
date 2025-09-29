import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../services/api'

export const useTodoStore = defineStore('todo', () => {
        const todos = ref<api.Todo[]>([])
        const meta = ref<api.Meta | null>(null)
        const loading = ref(false)
        const error = ref<string | null>(null)


        // query params
        const page = ref(1)
        const pageSize = ref(5)
        const completedFilter = ref<string | undefined>(undefined) // 'true' | 'false' | undefined
        const sort = ref('created_at desc')


        async function fetchTodos() {
                loading.value = true
                error.value = null
                try {
                        const params: any = { page: page.value, page_size: pageSize.value, sort: sort.value }
                        if (completedFilter.value !== undefined) params.completed = completedFilter.value
                        const res = await api.listTodos(params)
                        todos.value = res.items
                        meta.value = res.meta
                } catch (e: any) {
                        error.value = e?.message || 'fetch error'
                } finally {
                        loading.value = false
                }
        }

        async function getTodo(id: number) {
                loading.value = true
                error.value = null
                try {
                        const t = await api.getTodo(id)
                        return t
                } catch (e: any) {
                        error.value = e?.message || 'fetch error'
                } finally {
                        loading.value = false
                }
        }

        async function addTodo(title: string, description?: string) {
                loading.value = true
                try {
                        const t = await api.createTodo({ title, description })
                        // prepend to list
                        todos.value.unshift(t)
                } catch (e: any) {
                        error.value = e?.message || 'create error'
                } finally {
                        loading.value = false
                }
        }

        async function patchTodo(id: number, payload: Partial<{ title: string; description: string | null; completed: boolean }>) {
                loading.value = true
                try {
                        const t = await api.updateTodo(id, payload)
                        const idx = todos.value.findIndex(x => x.id === t.id)
                        if (idx !== -1) todos.value[idx] = t
                } catch (e: any) {
                        error.value = e?.message || 'update error'
                } finally {
                        loading.value = false
                }
        }


        async function removeTodo(id: number) {
                loading.value = true
                try {
                        await api.deleteTodo(id)
                        todos.value = todos.value.filter(x => x.id !== id)
                } catch (e: any) {
                        error.value = e?.message || 'delete error'
                } finally {
                        loading.value = false
                }
        }


        function setPage(p: number) {
                page.value = p
        }


        function setFilterCompleted(v: string | undefined) {
                completedFilter.value = v
        }


        return {
                todos,
                meta,
                loading,
                error,
                page,
                pageSize,
                completedFilter,
                sort,
                fetchTodos,
                getTodo,
                addTodo,
                patchTodo,
                removeTodo,
                setPage,
                setFilterCompleted,
        }
})