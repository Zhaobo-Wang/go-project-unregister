<template>
    <div class="app">
        <header>
        <h2>Todos</h2>
        </header>


        <section class="new">
            <h3>Todo Form</h3>
            <form @submit.prevent="onAdd">
                <input v-model="title" placeholder="Add todo title" required />
                <input v-model="description" placeholder="Add todo description" required />
                <button type="submit">{{ editingId === null ? 'Post Todo' : 'Update Todo' }}</button>
                <button type="button" @click="reload">Reload</button>
                <button v-if="editingId !== null" type="button" @click="cancelEdit">Cancel</button>
            </form>
        </section>

        <section class="list">
            <h3>Todo List</h3>
            <label>
                <select v-model="filter">
                    <option value="">Get All</option>
                    <option value="true">Filter Completed</option>
                    <option value="false">Filter Active</option>
                </select>
            </label>
            <div v-if="store.loading">Loading...</div>
            <div v-if="store.error" class="error">{{ store.error }}</div>

            <ul>
                <li v-for="t in store.todos" :key="t.id" class="todo-item">
                    <input type="checkbox" :checked="t.completed" @change="toggle(t)" />
                    <span :class="{ done: t.completed }">{{ t.title }}</span>
                    <span :class="{ done: t.completed }">{{ t.description }}</span>
                    <button @click="edit(t.id)">Update</button>
                    <button @click="del(t.id)">Delete</button>
                </li>
            </ul>
        </section>


        <footer class="pagination" v-if="store.meta">
            <button :disabled="!hasPrev" @click="goPrev">Prev</button>
            <span>Page {{ store.meta.page }} / {{ store.meta.total_pages }}</span>
            <button :disabled="!hasNext" @click="goNext">Next</button>
        </footer>
    </div>

</template>


<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { useTodoStore } from './stores/todo'


const store = useTodoStore()
const title = ref('')
const description = ref('')
const filter = ref('')
const editingId = ref<number | null>(null)

function reload() {
    store.fetchTodos()
}


async function onAdd() {
    
    if (editingId.value === null) {
        await store.addTodo(title.value.trim(), description.value.trim() || undefined)
        store.setPage(1)
        await store.fetchTodos()
    } else {
        try {
          const payload: any = { title: title.value.trim(), description: description.value.trim() || null }
          await store.patchTodo(editingId.value, payload)
          await store.fetchTodos()
        } catch (e) {
        } finally {
          editingId.value = null
        }
    }
    title.value = ''
    description.value = ''
}


async function toggle(t: any) {
    await store.patchTodo(t.id, { completed: !t.completed })
}

async function del(id: number | null) {
    if (id === null) return
    try {
        if (!confirm('Delete this todo?')) return
        await store.removeTodo(id)
    } catch (e) {
        console.error('failed to delete todo', e)
    }
}

async function edit(id: number | null) {
    if (id === null) return
    try {
        const todo = await store.getTodo(id)
        title.value = todo?.title || ''
        description.value = todo?.description || ''
        editingId.value = todo?.id ?? null
    } catch (e) {
        console.error('failed to load todo for edit', e)
    }
} 

function cancelEdit() {
    editingId.value = null
    title.value = ''
    description.value = ''
}

watch(filter, (v) => {
    if (v === '') store.setFilterCompleted(undefined)
    else store.setFilterCompleted(v)
    store.setPage(1)
    store.fetchTodos()
})


const hasPrev = computed(() => !!(store.meta && store.meta.page > 1))
const hasNext = computed(() => !!(store.meta && store.meta.page < (store.meta.total_pages || 0)))


function goPrev() {
    if (!store.meta) return
    store.setPage(store.meta.page - 1)
    store.fetchTodos()
}

function goNext() {
    if (!store.meta) return
    store.setPage(store.meta.page + 1)
    store.fetchTodos()
}


onMounted(() => {
    store.fetchTodos()
})
</script>


<style scoped>
    .app { max-width: 640px; margin: 32px auto; font-family: system-ui, -apple-system, 'Segoe UI', Roboto, 'Helvetica Neue', Arial }
    header h1 { margin: 0 0 16px }
    .new form { display:flex; gap:8px }
    .new input { flex:1; padding:8px }
    button { padding:6px 10px }
    .list ul { list-style:none; padding:0 }
    .todo-item { display:flex; align-items:center; gap:8px; padding:8px 0; border-bottom:1px solid #eee }
    .todo-item .done { text-decoration: line-through; opacity:0.7 }
    .error { color: #b00020 }
    .pagination { display:flex; gap:12px; align-items:center; justify-content:center; margin-top:16px }
</style>