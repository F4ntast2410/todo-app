import { apiRequest } from './client.js';

export async function login(email, password) {
    return apiRequest('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    });
}
export async function loginWithTg(user) {
    return apiRequest('/auth/telegram/login', {
        method: 'POST',
        body: JSON.stringify(user)
    })
}
export async function verify2FA(email, inputCode) {
    return apiRequest('/auth/verify2fa', {
        method: 'POST',
        body: JSON.stringify({ email: email, input_code: inputCode })
    });
}

export async function register(email) {
    return apiRequest('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email })
    });
}
export async function registerVerify(email, password, username, inputCode) {
    return apiRequest('/auth/register/verify', {
        method: 'POST',
        body: JSON.stringify({ email, password, username, input_code: inputCode })
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
