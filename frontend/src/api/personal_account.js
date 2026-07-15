import { apiRequest } from './client.js';

export async function changeInfo(payload) {
    return apiRequest('/auth/me', {
        method: 'POST',
        body: JSON.stringify(payload)
    });
}

export async function changePassword(payload) {
    return apiRequest('/auth/me/password', {
        method: 'POST',
        body: JSON.stringify(payload)
    });
}