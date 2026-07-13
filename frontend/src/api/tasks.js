import { apiRequest } from './client.js';

// Получить все активные задачи
export async function getTasks() {
    return apiRequest('/tasks/', {
        method: 'GET'
    });
}
export async function getTrashTasks() {
    return apiRequest('/tasks/trash', {
        method: 'GET',
    });
}

// Создать новую задачу
export async function createTask(taskData) {
    return apiRequest('/tasks', {
        method: 'POST',
        body: JSON.stringify(taskData)
    });
}

export async function deleteTask(user_task_id) {
    return apiRequest(`/tasks/${user_task_id}`, {
        method: 'DELETE',
    })
}
export async function updateTask(user_task_id, taskData) {
    return apiRequest(`/tasks/${user_task_id}`, {
        method: 'PUT',
        body: JSON.stringify(taskData)
    })
}
export async function restoreTask(user_task_id) {
    return apiRequest(`/tasks/restore/${user_task_id}`, {
        method: 'PUT',
    })
}
