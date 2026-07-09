import { apiRequest } from './client.js';

export async function login(email, password) {
    return apiRequest('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    });
}

// Задел на будущее
export async function register(email, password, username) {
    return apiRequest('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password, username })
    });
}

// Функция для проверки текущей сессии (получение профиля)
export async function getProfile() {
    return apiRequest('/auth/me', {
        method: 'GET'
    });
}

export async function logout() {
    return apiRequest('/auth/logout', {
        method: 'POST'
    });
}