import { apiRequest } from './client.js';

export async function changeInfo(payload) {
    return apiRequest('/auth/me', {
        method: 'PUT',
        body: JSON.stringify(payload)
    });
}

export async function changePassword(payload) {
    return apiRequest('/auth/me/password', {
        method: 'PUT',
        body: JSON.stringify(payload)
    });
}

export async function checkLinkTg() {
    return apiRequest('/auth/telegram/link', {
        method: 'GET',
    });
}

export async function linkTg(user) {
    return apiRequest('/auth/telegram/link', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    });
}